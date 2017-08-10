package documents

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"testing/quick"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
	"github.com/trussle/snowy/pkg/document"
	metricMocks "github.com/trussle/snowy/pkg/metrics/mocks"
	"github.com/trussle/snowy/pkg/repository"
	repoMocks "github.com/trussle/snowy/pkg/repository/mocks"
	"github.com/trussle/snowy/pkg/uuid"
)

func TestGetAPI(t *testing.T) {
	t.Parallel()

	t.Run("get with no resource_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func() bool {
			var (
				clients  = metricMocks.NewMockGauge(ctrl)
				duration = metricMocks.NewMockHistogramVec(ctrl)
				observer = metricMocks.NewMockObserver(ctrl)
				repo     = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
				server = httptest.NewServer(api)
			)

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("GET", "/", "400").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			resp, err := http.Get(server.URL)
			if err != nil {
				t.Error(err)
			}
			defer resp.Body.Close()

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get with invalid resource_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func() bool {
			var (
				clients  = metricMocks.NewMockGauge(ctrl)
				duration = metricMocks.NewMockHistogramVec(ctrl)
				observer = metricMocks.NewMockObserver(ctrl)
				repo     = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
				server = httptest.NewServer(api)
			)

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("GET", "/", "400").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			resp, err := http.Get(fmt.Sprintf("%s?resource_id=%s", server.URL, "bad"))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusBadRequest, resp.StatusCode; expected != actual {
				t.Errorf("expected: %d, actual: %d", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get with resource_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(uid uuid.UUID) bool {
			var (
				clients  = metricMocks.NewMockGauge(ctrl)
				duration = metricMocks.NewMockHistogramVec(ctrl)
				observer = metricMocks.NewMockObserver(ctrl)
				repo     = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
				server = httptest.NewServer(api)

				doc, err = document.Build(
					document.WithResourceID(uid),
				)
			)
			if err != nil {
				t.Fatal(err)
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("GET", "/", "200").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			repo.EXPECT().GetDocument(uid, repository.Query{}).Times(1).Return(doc, nil)

			resp, err := http.Get(fmt.Sprintf("%s?resource_id=%s", server.URL, uid))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusOK, resp.StatusCode; expected != actual {
				t.Errorf("expected: %d, actual: %d", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get with resource_id and tags", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(uid uuid.UUID, tags Tags) bool {
			var (
				clients  = metricMocks.NewMockGauge(ctrl)
				duration = metricMocks.NewMockHistogramVec(ctrl)
				observer = metricMocks.NewMockObserver(ctrl)
				repo     = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
				server = httptest.NewServer(api)

				doc, err = document.Build(
					document.WithResourceID(uid),
				)
			)
			if err != nil {
				t.Fatal(err)
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("GET", "/", "200").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			repo.EXPECT().GetDocument(uid, repository.Query{
				Tags: tags.Slice(),
			}).Times(1).Return(doc, nil)

			resp, err := http.Get(fmt.Sprintf("%s?resource_id=%s&query.tags=%s", server.URL, uid, tags.String()))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusOK, resp.StatusCode; expected != actual {
				t.Errorf("expected: %d, actual: %d", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get with resource_id but repo not found failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(uid uuid.UUID) bool {
			var (
				clients  = metricMocks.NewMockGauge(ctrl)
				duration = metricMocks.NewMockHistogramVec(ctrl)
				observer = metricMocks.NewMockObserver(ctrl)
				repo     = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
				server = httptest.NewServer(api)

				doc, err = document.Build(
					document.WithResourceID(uid),
				)
			)
			if err != nil {
				t.Fatal(err)
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("GET", "/", "404").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			repo.EXPECT().GetDocument(uid, repository.Query{}).Times(1).Return(doc, errNotFound{errors.New("failure")})

			resp, err := http.Get(fmt.Sprintf("%s?resource_id=%s", server.URL, uid))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusNotFound, resp.StatusCode; expected != actual {
				t.Errorf("expected: %d, actual: %d", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get with resource_id but repo failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(uid uuid.UUID) bool {
			var (
				clients  = metricMocks.NewMockGauge(ctrl)
				duration = metricMocks.NewMockHistogramVec(ctrl)
				observer = metricMocks.NewMockObserver(ctrl)
				repo     = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
				server = httptest.NewServer(api)

				doc, err = document.Build(
					document.WithResourceID(uid),
				)
			)
			if err != nil {
				t.Fatal(err)
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("GET", "/", "500").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			repo.EXPECT().GetDocument(uid, repository.Query{}).Times(1).Return(doc, errors.New("failure"))

			resp, err := http.Get(fmt.Sprintf("%s?resource_id=%s", server.URL, uid))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusInternalServerError, resp.StatusCode; expected != actual {
				t.Errorf("expected: %d, actual: %d", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestPostAPI(t *testing.T) {
	t.Parallel()

	t.Run("post with no body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			clients  = metricMocks.NewMockGauge(ctrl)
			duration = metricMocks.NewMockHistogramVec(ctrl)
			observer = metricMocks.NewMockObserver(ctrl)
			repo     = repoMocks.NewMockRepository(ctrl)

			api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
			server = httptest.NewServer(api)
		)

		fn := func() bool {
			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("POST", "/", "400").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			resp, err := http.Post(server.URL, "application/json", nil)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusBadRequest, resp.StatusCode; expected != actual {
				t.Errorf("expected: %d, actual: %d", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("post with empty body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			clients  = metricMocks.NewMockGauge(ctrl)
			duration = metricMocks.NewMockHistogramVec(ctrl)
			observer = metricMocks.NewMockObserver(ctrl)
			repo     = repoMocks.NewMockRepository(ctrl)

			api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
			server = httptest.NewServer(api)
		)

		fn := func() bool {

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("POST", "/", "400").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			body := strings.NewReader("{}")
			resp, err := http.Post(server.URL, "application/json", body)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusBadRequest, resp.StatusCode; expected != actual {
				t.Errorf("expected: %d, actual: %d", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("post with body with no authorID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			clients  = metricMocks.NewMockGauge(ctrl)
			duration = metricMocks.NewMockHistogramVec(ctrl)
			observer = metricMocks.NewMockObserver(ctrl)
			repo     = repoMocks.NewMockRepository(ctrl)

			api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
			server = httptest.NewServer(api)
		)

		fn := func(name string, tags Tags) bool {
			if len(name) == 0 {
				return true
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("POST", "/", "400").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			b, err := json.Marshal(struct {
				Name     string   `json:"name"`
				AuthorID string   `json:"author_id"`
				Tags     []string `json:"tags"`
			}{
				Name:     name,
				AuthorID: "",
				Tags:     tags,
			})
			if err != nil {
				t.Fatal(err)
			}

			resp, err := http.Post(server.URL, "application/json", bytes.NewReader(b))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusBadRequest, resp.StatusCode; expected != actual {
				t.Fatalf("expected: %d, actual: %d", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("post with body invalid content-type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			clients  = metricMocks.NewMockGauge(ctrl)
			duration = metricMocks.NewMockHistogramVec(ctrl)
			observer = metricMocks.NewMockObserver(ctrl)
			repo     = repoMocks.NewMockRepository(ctrl)

			api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
			server = httptest.NewServer(api)
		)

		fn := func(name, authorID, contentType string, tags Tags) bool {
			if len(name) == 0 || len(authorID) == 0 {
				return true
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("POST", "/", "400").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			b, err := json.Marshal(struct {
				Name     string   `json:"name"`
				AuthorID string   `json:"author_id"`
				Tags     []string `json:"tags"`
			}{
				Name:     name,
				AuthorID: authorID,
				Tags:     tags,
			})
			if err != nil {
				t.Fatal(err)
			}

			resp, err := http.Post(server.URL, contentType, bytes.NewReader(b))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusBadRequest, resp.StatusCode; expected != actual {
				t.Fatalf("expected: %d, actual: %d", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("post with body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			clients  = metricMocks.NewMockGauge(ctrl)
			duration = metricMocks.NewMockHistogramVec(ctrl)
			observer = metricMocks.NewMockObserver(ctrl)
			repo     = repoMocks.NewMockRepository(ctrl)

			api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
			server = httptest.NewServer(api)
		)

		fn := func(name, authorID string, tags Tags) bool {
			if len(name) == 0 || len(authorID) == 0 {
				return true
			}

			doc, err := document.Build(
				document.WithName(name),
				document.WithAuthorID(authorID),
				document.WithTags(tags),
			)
			if err != nil {
				t.Fatal(err)
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("POST", "/", "200").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)
			repo.EXPECT().InsertDocument(Document(doc)).Return(doc, nil).Times(1)

			b, err := json.Marshal(struct {
				Name     string   `json:"name"`
				AuthorID string   `json:"author_id"`
				Tags     []string `json:"tags"`
			}{
				Name:     name,
				AuthorID: authorID,
				Tags:     tags,
			})
			if err != nil {
				t.Fatal(err)
			}

			resp, err := http.Post(server.URL, "application/json", bytes.NewReader(b))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusOK, resp.StatusCode; expected != actual {
				t.Fatalf("expected: %d, actual: %d", expected, actual)
			}

			b, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			var resDoc struct {
				ResourceID uuid.UUID `json:"resource_id"`
			}
			if err := json.Unmarshal(b, &resDoc); err != nil {
				t.Fatal(err)
			}

			return !resDoc.ResourceID.Zero()
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("post with body but with repo failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			clients  = metricMocks.NewMockGauge(ctrl)
			duration = metricMocks.NewMockHistogramVec(ctrl)
			observer = metricMocks.NewMockObserver(ctrl)
			repo     = repoMocks.NewMockRepository(ctrl)

			api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
			server = httptest.NewServer(api)
		)

		fn := func(name, authorID string, tags Tags) bool {
			if len(name) == 0 || len(authorID) == 0 {
				return true
			}

			doc, err := document.Build(
				document.WithName(name),
				document.WithAuthorID(authorID),
				document.WithTags(tags),
			)
			if err != nil {
				t.Fatal(err)
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("POST", "/", "500").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)
			repo.EXPECT().InsertDocument(Document(doc)).Return(doc, errors.New("bad")).Times(1)

			b, err := json.Marshal(struct {
				Name     string   `json:"name"`
				AuthorID string   `json:"author_id"`
				Tags     []string `json:"tags"`
			}{
				Name:     name,
				AuthorID: authorID,
				Tags:     tags,
			})
			if err != nil {
				t.Fatal(err)
			}

			resp, err := http.Post(server.URL, "application/json", bytes.NewReader(b))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusInternalServerError, resp.StatusCode; expected != actual {
				t.Fatalf("expected: %d, actual: %d", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestPutAPI(t *testing.T) {
	t.Parallel()

	t.Run("put with no resource_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			clients  = metricMocks.NewMockGauge(ctrl)
			duration = metricMocks.NewMockHistogramVec(ctrl)
			observer = metricMocks.NewMockObserver(ctrl)
			repo     = repoMocks.NewMockRepository(ctrl)

			api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
			server = httptest.NewServer(api)
		)

		fn := func() bool {
			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("PUT", "/", "400").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			resp, err := Put(server.URL, "application/json", nil)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusBadRequest, resp.StatusCode; expected != actual {
				t.Errorf("expected: %d, actual: %d", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("put with invalid resource_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			clients  = metricMocks.NewMockGauge(ctrl)
			duration = metricMocks.NewMockHistogramVec(ctrl)
			observer = metricMocks.NewMockObserver(ctrl)
			repo     = repoMocks.NewMockRepository(ctrl)

			api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
			server = httptest.NewServer(api)
		)

		fn := func(resourceID ASCII) bool {
			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("PUT", "/", "400").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			resp, err := Put(fmt.Sprintf("%s?resource_id=%s", server.URL, resourceID), "application/json", nil)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusBadRequest, resp.StatusCode; expected != actual {
				t.Errorf("expected: %d, actual: %d", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("put with no body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			clients  = metricMocks.NewMockGauge(ctrl)
			duration = metricMocks.NewMockHistogramVec(ctrl)
			observer = metricMocks.NewMockObserver(ctrl)
			repo     = repoMocks.NewMockRepository(ctrl)

			api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
			server = httptest.NewServer(api)
		)

		fn := func(resourceID uuid.UUID) bool {
			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("PUT", "/", "400").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			resp, err := Put(fmt.Sprintf("%s?resource_id=%s", server.URL, resourceID), "application/json", nil)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusBadRequest, resp.StatusCode; expected != actual {
				t.Errorf("expected: %d, actual: %d", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("put with empty body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			clients  = metricMocks.NewMockGauge(ctrl)
			duration = metricMocks.NewMockHistogramVec(ctrl)
			observer = metricMocks.NewMockObserver(ctrl)
			repo     = repoMocks.NewMockRepository(ctrl)

			api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
			server = httptest.NewServer(api)
		)

		fn := func(resourceID uuid.UUID) bool {

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("PUT", "/", "400").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			body := strings.NewReader("{}")
			resp, err := Put(fmt.Sprintf("%s?resource_id=%s", server.URL, resourceID), "application/json", body)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusBadRequest, resp.StatusCode; expected != actual {
				t.Errorf("expected: %d, actual: %d", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("put with body with no authorID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			clients  = metricMocks.NewMockGauge(ctrl)
			duration = metricMocks.NewMockHistogramVec(ctrl)
			observer = metricMocks.NewMockObserver(ctrl)
			repo     = repoMocks.NewMockRepository(ctrl)

			api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
			server = httptest.NewServer(api)
		)

		fn := func(resourceID uuid.UUID, name string, tags Tags) bool {
			if len(name) == 0 {
				return true
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("PUT", "/", "400").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			b, err := json.Marshal(struct {
				Name     string   `json:"name"`
				AuthorID string   `json:"author_id"`
				Tags     []string `json:"tags"`
			}{
				Name:     name,
				AuthorID: "",
				Tags:     tags,
			})
			if err != nil {
				t.Fatal(err)
			}

			resp, err := Put(fmt.Sprintf("%s?resource_id=%s", server.URL, resourceID), "application/json", bytes.NewReader(b))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusBadRequest, resp.StatusCode; expected != actual {
				t.Fatalf("expected: %d, actual: %d", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("put with body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			clients  = metricMocks.NewMockGauge(ctrl)
			duration = metricMocks.NewMockHistogramVec(ctrl)
			observer = metricMocks.NewMockObserver(ctrl)
			repo     = repoMocks.NewMockRepository(ctrl)

			api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
			server = httptest.NewServer(api)
		)

		fn := func(resourceID uuid.UUID, name, authorID string, tags Tags) bool {
			if len(name) == 0 || len(authorID) == 0 {
				return true
			}

			doc, err := document.Build(
				document.WithName(name),
				document.WithAuthorID(authorID),
				document.WithTags(tags),
			)
			if err != nil {
				t.Fatal(err)
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("PUT", "/", "200").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)
			repo.EXPECT().AppendDocument(resourceID, Document(doc)).Return(doc, nil).Times(1)

			b, err := json.Marshal(struct {
				Name     string   `json:"name"`
				AuthorID string   `json:"author_id"`
				Tags     []string `json:"tags"`
			}{
				Name:     name,
				AuthorID: authorID,
				Tags:     tags,
			})
			if err != nil {
				t.Fatal(err)
			}

			resp, err := Put(fmt.Sprintf("%s?resource_id=%s", server.URL, resourceID), "application/json", bytes.NewReader(b))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusOK, resp.StatusCode; expected != actual {
				t.Fatalf("expected: %d, actual: %d", expected, actual)
			}

			b, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			var resDoc struct {
				ResourceID uuid.UUID `json:"resource_id"`
			}
			if err := json.Unmarshal(b, &resDoc); err != nil {
				t.Fatal(err)
			}

			return !resDoc.ResourceID.Zero()
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("put with body but with repo failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			clients  = metricMocks.NewMockGauge(ctrl)
			duration = metricMocks.NewMockHistogramVec(ctrl)
			observer = metricMocks.NewMockObserver(ctrl)
			repo     = repoMocks.NewMockRepository(ctrl)

			api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
			server = httptest.NewServer(api)
		)

		fn := func(resourceID uuid.UUID, name, authorID string, tags Tags) bool {
			if len(name) == 0 || len(authorID) == 0 {
				return true
			}

			doc, err := document.Build(
				document.WithName(name),
				document.WithAuthorID(authorID),
				document.WithTags(tags),
			)
			if err != nil {
				t.Fatal(err)
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("PUT", "/", "500").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)
			repo.EXPECT().AppendDocument(resourceID, Document(doc)).Return(doc, errors.New("bad")).Times(1)

			b, err := json.Marshal(struct {
				Name     string   `json:"name"`
				AuthorID string   `json:"author_id"`
				Tags     []string `json:"tags"`
			}{
				Name:     name,
				AuthorID: authorID,
				Tags:     tags,
			})
			if err != nil {
				t.Fatal(err)
			}

			resp, err := Put(fmt.Sprintf("%s?resource_id=%s", server.URL, resourceID), "application/json", bytes.NewReader(b))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusInternalServerError, resp.StatusCode; expected != actual {
				t.Fatalf("expected: %d, actual: %d", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestGetMultipleAPI(t *testing.T) {
	t.Parallel()

	t.Run("get with no resource_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func() bool {
			var (
				clients  = metricMocks.NewMockGauge(ctrl)
				duration = metricMocks.NewMockHistogramVec(ctrl)
				observer = metricMocks.NewMockObserver(ctrl)
				repo     = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
				server = httptest.NewServer(api)
			)

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("GET", "/multiple", "400").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			resp, err := http.Get(fmt.Sprintf("%s/multiple", server.URL))
			if err != nil {
				t.Error(err)
			}
			defer resp.Body.Close()

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get with invalid resource_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func() bool {
			var (
				clients  = metricMocks.NewMockGauge(ctrl)
				duration = metricMocks.NewMockHistogramVec(ctrl)
				observer = metricMocks.NewMockObserver(ctrl)
				repo     = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
				server = httptest.NewServer(api)
			)

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("GET", "/multiple", "400").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			resp, err := http.Get(fmt.Sprintf("%s/multiple?resource_id=%s", server.URL, "bad"))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusBadRequest, resp.StatusCode; expected != actual {
				t.Errorf("expected: %d, actual: %d", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get with resource_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(uid uuid.UUID) bool {
			var (
				clients  = metricMocks.NewMockGauge(ctrl)
				duration = metricMocks.NewMockHistogramVec(ctrl)
				observer = metricMocks.NewMockObserver(ctrl)
				repo     = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
				server = httptest.NewServer(api)

				doc, err = document.Build(
					document.WithResourceID(uid),
				)
				docs = []document.Document{doc}
			)
			if err != nil {
				t.Fatal(err)
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("GET", "/multiple", "200").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			repo.EXPECT().GetDocuments(uid, repository.Query{}).Times(1).Return(docs, nil)

			resp, err := http.Get(fmt.Sprintf("%s/multiple?resource_id=%s", server.URL, uid))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusOK, resp.StatusCode; expected != actual {
				t.Errorf("expected: %d, actual: %d", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get with resource_id and tags", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(uid uuid.UUID, tags Tags) bool {
			var (
				clients  = metricMocks.NewMockGauge(ctrl)
				duration = metricMocks.NewMockHistogramVec(ctrl)
				observer = metricMocks.NewMockObserver(ctrl)
				repo     = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
				server = httptest.NewServer(api)

				doc, err = document.Build(
					document.WithResourceID(uid),
				)
				docs = []document.Document{doc}
			)
			if err != nil {
				t.Fatal(err)
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("GET", "/multiple", "200").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			repo.EXPECT().GetDocuments(uid, repository.Query{
				Tags: tags.Slice(),
			}).Times(1).Return(docs, nil)

			resp, err := http.Get(fmt.Sprintf("%s/multiple?resource_id=%s&query.tags=%s", server.URL, uid, tags.String()))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusOK, resp.StatusCode; expected != actual {
				t.Errorf("expected: %d, actual: %d", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get with resource_id but with repo failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(uid uuid.UUID) bool {
			var (
				clients  = metricMocks.NewMockGauge(ctrl)
				duration = metricMocks.NewMockHistogramVec(ctrl)
				observer = metricMocks.NewMockObserver(ctrl)
				repo     = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
				server = httptest.NewServer(api)

				doc, err = document.Build(
					document.WithResourceID(uid),
				)
				docs = []document.Document{doc}
			)
			if err != nil {
				t.Fatal(err)
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("GET", "/multiple", "500").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			repo.EXPECT().GetDocuments(uid, repository.Query{}).Times(1).Return(docs, errors.New("bad"))

			resp, err := http.Get(fmt.Sprintf("%s/multiple?resource_id=%s", server.URL, uid))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusInternalServerError, resp.StatusCode; expected != actual {
				t.Errorf("expected: %d, actual: %d", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestNotFoundAPI(t *testing.T) {
	t.Parallel()

	t.Run("get", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(resource ASCII) bool {
			var (
				clients  = metricMocks.NewMockGauge(ctrl)
				duration = metricMocks.NewMockHistogramVec(ctrl)
				observer = metricMocks.NewMockObserver(ctrl)
				repo     = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
				server = httptest.NewServer(api)
			)

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("GET", fmt.Sprintf("/%s", resource), "404").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			resp, err := http.Get(fmt.Sprintf("%s/%s", server.URL, resource))
			if err != nil {
				t.Error(err)
			}
			defer resp.Body.Close()

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func Put(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return http.DefaultClient.Do(req)
}

type float64Matcher struct{}

func (float64Matcher) Matches(x interface{}) bool {
	_, ok := x.(float64)
	return ok
}

func (float64Matcher) String() string {
	return "is float64"
}

func Float64() gomock.Matcher { return float64Matcher{} }

type documentMatcher struct {
	doc document.Document
}

func (m documentMatcher) Matches(x interface{}) bool {
	d, ok := x.(document.Document)
	if !ok {
		return false
	}

	return m.doc.Name() == d.Name() &&
		m.doc.AuthorID() == d.AuthorID() &&
		reflect.DeepEqual(m.doc.Tags(), d.Tags())
}

func (documentMatcher) String() string {
	return "is document"
}

func Document(doc document.Document) gomock.Matcher { return documentMatcher{doc} }

type errNotFound struct {
	err error
}

func (e errNotFound) Error() string {
	return e.err.Error()
}

func (e errNotFound) NotFound() bool {
	return true
}
