package journals

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"testing/quick"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	metricMocks "github.com/trussle/snowy/pkg/metrics/mocks"
	"github.com/trussle/snowy/pkg/models"
	repoMocks "github.com/trussle/snowy/pkg/repository/mocks"
	"github.com/trussle/snowy/pkg/uuid"
)

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

			body := strings.NewReader("")
			resp, err := http.Post(server.URL, "multipart/form-data", body)
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

		fn := func(name string, tags Tags, conBytes []byte) bool {
			if len(name) == 0 {
				return true
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("POST", "/", "400").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			docBytes, err := json.Marshal(struct {
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

			var (
				buffer bytes.Buffer
				writer = multipart.NewWriter(&buffer)
			)

			MustWriteField(writer, contentFormFile, "application/octet-stream", conBytes)
			MustWriteField(writer, documentFormFile, "application/json", docBytes)

			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}

			resp, err := http.Post(server.URL, writer.FormDataContentType(), &buffer)
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

		fn := func(name, authorID, contentType string, tags Tags, conBytes []byte) bool {
			if len(name) == 0 || len(authorID) == 0 {
				return true
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("POST", "/", "400").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			docBytes, err := json.Marshal(struct {
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

			var (
				buffer bytes.Buffer
				writer = multipart.NewWriter(&buffer)
			)

			MustWriteField(writer, contentFormFile, "application/octet-stream", conBytes)
			MustWriteField(writer, documentFormFile, "application/json", docBytes)

			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}

			resp, err := http.Post(server.URL, contentType, &buffer)
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

	t.Run("post with no content body", func(t *testing.T) {
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

			docBytes, err := json.Marshal(struct {
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

			var (
				buffer bytes.Buffer
				writer = multipart.NewWriter(&buffer)
			)

			MustWriteField(writer, documentFormFile, "application/json", docBytes)

			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}

			resp, err := http.Post(server.URL, writer.FormDataContentType(), &buffer)
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

		fn := func(name, authorID string, tags Tags, conBytes []byte) bool {
			if len(name) == 0 || len(authorID) == 0 || len(conBytes) == 0 {
				return true
			}

			doc, err := models.BuildLedger(
				models.WithName(name),
				models.WithAuthorID(authorID),
				models.WithTags(tags),
			)
			if err != nil {
				t.Fatal(err)
			}

			content, err := models.BuildContent(
				models.WithSize(int64(len(conBytes))),
				models.WithContentBytes(conBytes),
				models.WithContentType("application/octet-stream"),
			)
			if err != nil {
				t.Fatal(err)
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("POST", "/", "200").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)
			repo.EXPECT().PutContent(Content(content)).Return(content, nil).Times(1)
			repo.EXPECT().InsertLedger(Ledger(doc)).Return(doc, nil).Times(1)

			docBytes, err := json.Marshal(struct {
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

			var (
				buffer bytes.Buffer
				writer = multipart.NewWriter(&buffer)
			)

			MustWriteField(writer, contentFormFile, "application/octet-stream", conBytes)
			MustWriteField(writer, documentFormFile, "application/json", docBytes)

			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}

			resp, err := http.Post(server.URL, writer.FormDataContentType(), &buffer)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusOK, resp.StatusCode; expected != actual {
				t.Fatalf("expected: %d, actual: %d", expected, actual)
			}

			b, err := ioutil.ReadAll(resp.Body)
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

	t.Run("post with body but with repo ledger failure", func(t *testing.T) {
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

		fn := func(name, authorID string, tags Tags, conBytes []byte) bool {
			if len(name) == 0 || len(authorID) == 0 || len(conBytes) == 0 {
				return true
			}

			doc, err := models.BuildLedger(
				models.WithName(name),
				models.WithAuthorID(authorID),
				models.WithTags(tags),
			)
			if err != nil {
				t.Fatal(err)
			}

			content, err := models.BuildContent(
				models.WithSize(int64(len(conBytes))),
				models.WithContentBytes(conBytes),
				models.WithContentType("application/octet-stream"),
			)
			if err != nil {
				t.Fatal(err)
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("POST", "/", "500").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)
			repo.EXPECT().PutContent(Content(content)).Return(content, nil).Times(1)
			repo.EXPECT().InsertLedger(Ledger(doc)).Return(doc, errors.New("bad")).Times(1)

			docBytes, err := json.Marshal(struct {
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

			var (
				buffer bytes.Buffer
				writer = multipart.NewWriter(&buffer)
			)

			MustWriteField(writer, contentFormFile, "application/octet-stream", conBytes)
			MustWriteField(writer, documentFormFile, "application/json", docBytes)

			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}

			resp, err := http.Post(server.URL, writer.FormDataContentType(), &buffer)
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

	t.Run("post with body but with repo contents failure", func(t *testing.T) {
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

		fn := func(name, authorID string, tags Tags, conBytes []byte) bool {
			if len(name) == 0 || len(authorID) == 0 || len(conBytes) == 0 {
				return true
			}

			content, err := models.BuildContent(
				models.WithSize(int64(len(conBytes))),
				models.WithContentBytes(conBytes),
				models.WithContentType("application/octet-stream"),
			)
			if err != nil {
				t.Fatal(err)
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("POST", "/", "500").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)
			repo.EXPECT().PutContent(Content(content)).Return(content, errors.New("bad")).Times(1)

			docBytes, err := json.Marshal(struct {
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

			var (
				buffer bytes.Buffer
				writer = multipart.NewWriter(&buffer)
			)

			MustWriteField(writer, contentFormFile, "application/octet-stream", conBytes)
			MustWriteField(writer, documentFormFile, "application/json", docBytes)

			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}

			resp, err := http.Post(server.URL, writer.FormDataContentType(), &buffer)
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

			body := strings.NewReader("")
			resp, err := Put(fmt.Sprintf("%s?resource_id=%s", server.URL, resourceID), "multipart/form-data", body)
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

		fn := func(resourceID uuid.UUID, name string, tags Tags, conBytes []byte) bool {
			if len(name) == 0 {
				return true
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("PUT", "/", "400").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			docBytes, err := json.Marshal(struct {
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

			var (
				buffer bytes.Buffer
				writer = multipart.NewWriter(&buffer)
			)

			MustWriteField(writer, contentFormFile, "application/octet-stream", conBytes)
			MustWriteField(writer, documentFormFile, "application/json", docBytes)

			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}

			resp, err := Put(fmt.Sprintf("%s?resource_id=%s", server.URL, resourceID), writer.FormDataContentType(), &buffer)
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

	t.Run("put with body invalid content-type", func(t *testing.T) {
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

		fn := func(resourceID uuid.UUID, name, authorID, contentType string, tags Tags, conBytes []byte) bool {
			if len(name) == 0 || len(authorID) == 0 {
				return true
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("PUT", "/", "400").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			docBytes, err := json.Marshal(struct {
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

			var (
				buffer bytes.Buffer
				writer = multipart.NewWriter(&buffer)
			)

			MustWriteField(writer, contentFormFile, "application/octet-stream", conBytes)
			MustWriteField(writer, documentFormFile, "application/json", docBytes)

			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}

			resp, err := Put(fmt.Sprintf("%s?resource_id=%s", server.URL, resourceID), contentType, &buffer)
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

	t.Run("put with no content body", func(t *testing.T) {
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

			docBytes, err := json.Marshal(struct {
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

			var (
				buffer bytes.Buffer
				writer = multipart.NewWriter(&buffer)
			)

			MustWriteField(writer, documentFormFile, "application/json", docBytes)

			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}

			resp, err := Put(fmt.Sprintf("%s?resource_id=%s", server.URL, resourceID), writer.FormDataContentType(), &buffer)
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

		fn := func(resourceID uuid.UUID, name, authorID string, tags Tags, conBytes []byte) bool {
			if len(name) == 0 || len(authorID) == 0 || len(conBytes) == 0 {
				return true
			}

			doc, err := models.BuildLedger(
				models.WithName(name),
				models.WithAuthorID(authorID),
				models.WithTags(tags),
			)
			if err != nil {
				t.Fatal(err)
			}

			content, err := models.BuildContent(
				models.WithSize(int64(len(conBytes))),
				models.WithContentBytes(conBytes),
				models.WithContentType("application/octet-stream"),
			)
			if err != nil {
				t.Fatal(err)
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("PUT", "/", "200").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)
			repo.EXPECT().PutContent(Content(content)).Return(content, nil).Times(1)
			repo.EXPECT().AppendLedger(resourceID, Ledger(doc)).Return(doc, nil).Times(1)

			docBytes, err := json.Marshal(struct {
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

			var (
				buffer bytes.Buffer
				writer = multipart.NewWriter(&buffer)
			)

			MustWriteField(writer, contentFormFile, "application/octet-stream", conBytes)
			MustWriteField(writer, documentFormFile, "application/json", docBytes)

			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}

			resp, err := Put(fmt.Sprintf("%s?resource_id=%s", server.URL, resourceID), writer.FormDataContentType(), &buffer)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusOK, resp.StatusCode; expected != actual {
				t.Fatalf("expected: %d, actual: %d", expected, actual)
			}

			b, err := ioutil.ReadAll(resp.Body)
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

	t.Run("put with body but with repo ledger failure", func(t *testing.T) {
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

		fn := func(resourceID uuid.UUID, name, authorID string, tags Tags, conBytes []byte) bool {
			if len(name) == 0 || len(authorID) == 0 || len(conBytes) == 0 {
				return true
			}

			doc, err := models.BuildLedger(
				models.WithName(name),
				models.WithAuthorID(authorID),
				models.WithTags(tags),
			)
			if err != nil {
				t.Fatal(err)
			}

			content, err := models.BuildContent(
				models.WithSize(int64(len(conBytes))),
				models.WithContentBytes(conBytes),
				models.WithContentType("application/octet-stream"),
			)
			if err != nil {
				t.Fatal(err)
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("PUT", "/", "500").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)
			repo.EXPECT().PutContent(Content(content)).Return(content, nil).Times(1)
			repo.EXPECT().AppendLedger(resourceID, Ledger(doc)).Return(doc, errors.New("bad")).Times(1)

			docBytes, err := json.Marshal(struct {
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

			var (
				buffer bytes.Buffer
				writer = multipart.NewWriter(&buffer)
			)

			MustWriteField(writer, contentFormFile, "application/octet-stream", conBytes)
			MustWriteField(writer, documentFormFile, "application/json", docBytes)

			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}

			resp, err := Put(fmt.Sprintf("%s?resource_id=%s", server.URL, resourceID), writer.FormDataContentType(), &buffer)
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

	t.Run("put with body but with repo contents failure", func(t *testing.T) {
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

		fn := func(resourceID uuid.UUID, name, authorID string, tags Tags, conBytes []byte) bool {
			if len(name) == 0 || len(authorID) == 0 || len(conBytes) == 0 {
				return true
			}

			content, err := models.BuildContent(
				models.WithSize(int64(len(conBytes))),
				models.WithContentBytes(conBytes),
				models.WithContentType("application/octet-stream"),
			)
			if err != nil {
				t.Fatal(err)
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("PUT", "/", "500").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)
			repo.EXPECT().PutContent(Content(content)).Return(content, errors.New("bad")).Times(1)

			docBytes, err := json.Marshal(struct {
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

			var (
				buffer bytes.Buffer
				writer = multipart.NewWriter(&buffer)
			)

			MustWriteField(writer, contentFormFile, "application/octet-stream", conBytes)
			MustWriteField(writer, documentFormFile, "application/json", docBytes)

			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}

			resp, err := Put(fmt.Sprintf("%s?resource_id=%s", server.URL, resourceID), writer.FormDataContentType(), &buffer)
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

func MustWriteField(w *multipart.Writer, name, contentType string, b []byte) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(name), escapeQuotes(name)))
	h.Set("Content-Type", contentType)
	h.Set("Content-Length", strconv.Itoa(len(b)))
	p, err := w.CreatePart(h)
	if err != nil {
		panic(err)
	}
	_, err = p.Write(b)
	if err != nil {
		panic(err)
	}
}

func Put(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return http.DefaultClient.Do(req)
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
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

type ledgerMatcher struct {
	doc models.Ledger
}

func (m ledgerMatcher) Matches(x interface{}) bool {
	d, ok := x.(models.Ledger)
	if !ok {
		return false
	}

	return m.doc.Name() == d.Name() &&
		m.doc.AuthorID() == d.AuthorID() &&
		reflect.DeepEqual(m.doc.Tags(), d.Tags())
}

func (ledgerMatcher) String() string {
	return "is ledger"
}

func Ledger(doc models.Ledger) gomock.Matcher { return ledgerMatcher{doc} }

type contentMatcher struct {
	content models.Content
}

func (m contentMatcher) Matches(x interface{}) bool {
	d, ok := x.(models.Content)
	if !ok {
		return false
	}

	res := m.content.Address() == d.Address() &&
		m.content.ContentType() == d.ContentType() &&
		m.content.Size() == d.Size()

	return res
}

func (contentMatcher) String() string {
	return "is content"
}

func Content(content models.Content) gomock.Matcher { return contentMatcher{content} }

type errNotFound struct {
	err error
}

func (e errNotFound) Error() string {
	return e.err.Error()
}

func (e errNotFound) NotFound() bool {
	return true
}
