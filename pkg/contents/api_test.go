package contents

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"testing/quick"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
	metricMocks "github.com/trussle/snowy/pkg/metrics/mocks"
	"github.com/trussle/snowy/pkg/models"
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
				clients      = metricMocks.NewMockGauge(ctrl)
				writtenBytes = metricMocks.NewMockCounter(ctrl)
				records      = metricMocks.NewMockCounter(ctrl)
				duration     = metricMocks.NewMockHistogramVec(ctrl)
				observer     = metricMocks.NewMockObserver(ctrl)
				repo         = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, writtenBytes, records, duration)
				server = httptest.NewServer(api)
			)
			defer api.Close()

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
				clients      = metricMocks.NewMockGauge(ctrl)
				writtenBytes = metricMocks.NewMockCounter(ctrl)
				records      = metricMocks.NewMockCounter(ctrl)
				duration     = metricMocks.NewMockHistogramVec(ctrl)
				observer     = metricMocks.NewMockObserver(ctrl)
				repo         = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, writtenBytes, records, duration)
				server = httptest.NewServer(api)
			)
			defer api.Close()

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

		fn := func(uid uuid.UUID, bytes []byte) bool {
			var (
				clients      = metricMocks.NewMockGauge(ctrl)
				writtenBytes = metricMocks.NewMockCounter(ctrl)
				records      = metricMocks.NewMockCounter(ctrl)
				duration     = metricMocks.NewMockHistogramVec(ctrl)
				observer     = metricMocks.NewMockObserver(ctrl)
				repo         = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, writtenBytes, records, duration)
				server = httptest.NewServer(api)

				content, err = models.BuildContent(
					models.WithSize(int64(len(bytes))),
					models.WithBytes(bytes),
				)
			)
			defer api.Close()

			if err != nil {
				t.Fatal(err)
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("GET", "/", "200").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			repo.EXPECT().GetContent(uid, Query()).Times(1).Return(content, nil)

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

	t.Run("get with resource_id but repo not found failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(uid uuid.UUID, bytes []byte) bool {
			var (
				clients      = metricMocks.NewMockGauge(ctrl)
				writtenBytes = metricMocks.NewMockCounter(ctrl)
				records      = metricMocks.NewMockCounter(ctrl)
				duration     = metricMocks.NewMockHistogramVec(ctrl)
				observer     = metricMocks.NewMockObserver(ctrl)
				repo         = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, writtenBytes, records, duration)
				server = httptest.NewServer(api)

				content, err = models.BuildContent(
					models.WithSize(int64(len(bytes))),
					models.WithBytes(bytes),
				)
			)
			defer api.Close()

			if err != nil {
				t.Fatal(err)
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("GET", "/", "404").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			repo.EXPECT().GetContent(uid, Query()).Times(1).Return(content, errNotFound{errors.New("failure")})

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

		fn := func(uid uuid.UUID, bytes []byte) bool {
			var (
				clients      = metricMocks.NewMockGauge(ctrl)
				writtenBytes = metricMocks.NewMockCounter(ctrl)
				records      = metricMocks.NewMockCounter(ctrl)
				duration     = metricMocks.NewMockHistogramVec(ctrl)
				observer     = metricMocks.NewMockObserver(ctrl)
				repo         = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, writtenBytes, records, duration)
				server = httptest.NewServer(api)

				content, err = models.BuildContent(
					models.WithSize(int64(len(bytes))),
					models.WithBytes(bytes),
				)
			)
			defer api.Close()

			if err != nil {
				t.Fatal(err)
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("GET", "/", "500").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			repo.EXPECT().GetContent(uid, Query()).Times(1).Return(content, errors.New("failure"))

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

func TestGetMultipleAPI(t *testing.T) {
	t.Parallel()

	t.Run("get multiple with no resource_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func() bool {
			var (
				clients      = metricMocks.NewMockGauge(ctrl)
				writtenBytes = metricMocks.NewMockCounter(ctrl)
				records      = metricMocks.NewMockCounter(ctrl)
				duration     = metricMocks.NewMockHistogramVec(ctrl)
				observer     = metricMocks.NewMockObserver(ctrl)
				repo         = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, writtenBytes, records, duration)
				server = httptest.NewServer(api)
			)
			defer api.Close()

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("GET", "/multiple/", "400").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			resp, err := http.Get(fmt.Sprintf("%s/multiple/", server.URL))
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

	t.Run("get multiple with invalid resource_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func() bool {
			var (
				clients      = metricMocks.NewMockGauge(ctrl)
				writtenBytes = metricMocks.NewMockCounter(ctrl)
				records      = metricMocks.NewMockCounter(ctrl)
				duration     = metricMocks.NewMockHistogramVec(ctrl)
				observer     = metricMocks.NewMockObserver(ctrl)
				repo         = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, writtenBytes, records, duration)
				server = httptest.NewServer(api)
			)
			defer api.Close()

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("GET", "/multiple/", "400").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			resp, err := http.Get(fmt.Sprintf("%s/multiple/?resource_id=%s", server.URL, "bad"))
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

	t.Run("get multiple with resource_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(uid uuid.UUID, bytes []byte) bool {
			var (
				clients      = metricMocks.NewMockGauge(ctrl)
				writtenBytes = metricMocks.NewMockCounter(ctrl)
				records      = metricMocks.NewMockCounter(ctrl)
				duration     = metricMocks.NewMockHistogramVec(ctrl)
				observer     = metricMocks.NewMockObserver(ctrl)
				repo         = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, writtenBytes, records, duration)
				server = httptest.NewServer(api)

				content, err = models.BuildContent(
					models.WithSize(int64(len(bytes))),
					models.WithBytes(bytes),
				)
			)
			defer api.Close()

			if err != nil {
				t.Fatal(err)
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("GET", "/multiple/", "200").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			repo.EXPECT().GetContents(UUID(uid), Query()).Times(1).Return([]models.Content{content}, nil)

			resp, err := http.Get(fmt.Sprintf("%s/multiple/?resource_id=%s", server.URL, uid))
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

	t.Run("get with resource_id but repo failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(uid uuid.UUID, bytes []byte) bool {
			var (
				clients      = metricMocks.NewMockGauge(ctrl)
				writtenBytes = metricMocks.NewMockCounter(ctrl)
				records      = metricMocks.NewMockCounter(ctrl)
				duration     = metricMocks.NewMockHistogramVec(ctrl)
				observer     = metricMocks.NewMockObserver(ctrl)
				repo         = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, writtenBytes, records, duration)
				server = httptest.NewServer(api)
			)
			defer api.Close()

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("GET", "/multiple/", "500").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			repo.EXPECT().GetContents(UUID(uid), Query()).Times(1).Return(nil, errNotFound{errors.New("failure")})

			resp, err := http.Get(fmt.Sprintf("%s/multiple/?resource_id=%s", server.URL, uid))
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

func TestPutAPI(t *testing.T) {
	t.Parallel()

	guard := func(fn func([]byte) bool) func([]byte) bool {
		return func(b []byte) bool {
			if len(b) < 1 {
				return true
			}
			return fn(b)
		}
	}

	t.Run("put with content too large", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			clients      = metricMocks.NewMockGauge(ctrl)
			writtenBytes = metricMocks.NewMockCounter(ctrl)
			records      = metricMocks.NewMockCounter(ctrl)
			duration     = metricMocks.NewMockHistogramVec(ctrl)
			observer     = metricMocks.NewMockObserver(ctrl)
			repo         = repoMocks.NewMockRepository(ctrl)

			api    = NewAPI(repo, log.NewNopLogger(), clients, writtenBytes, records, duration)
			server = httptest.NewServer(api)

			b = make([]byte, defaultMaxContentLength+1)
		)
		defer api.Close()

		clients.EXPECT().Inc().Times(1)
		clients.EXPECT().Dec().Times(1)

		duration.EXPECT().WithLabelValues("POST", "/", "400").Return(observer).Times(1)
		observer.EXPECT().Observe(Float64()).Times(1)

		resp, err := http.Post(server.URL, "plain/text", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if expected, actual := http.StatusBadRequest, resp.StatusCode; expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
	})

	t.Run("put with content too large, with content-length", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			clients      = metricMocks.NewMockGauge(ctrl)
			writtenBytes = metricMocks.NewMockCounter(ctrl)
			records      = metricMocks.NewMockCounter(ctrl)
			duration     = metricMocks.NewMockHistogramVec(ctrl)
			observer     = metricMocks.NewMockObserver(ctrl)
			repo         = repoMocks.NewMockRepository(ctrl)

			api    = NewAPI(repo, log.NewNopLogger(), clients, writtenBytes, records, duration)
			server = httptest.NewServer(api)

			b = make([]byte, defaultMaxContentLength+1)
		)
		defer api.Close()

		clients.EXPECT().Inc().Times(1)
		clients.EXPECT().Dec().Times(1)

		duration.EXPECT().WithLabelValues("POST", "/", "400").Return(observer).Times(1)
		observer.EXPECT().Observe(Float64()).Times(1)

		req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Length", strconv.FormatInt(defaultMaxContentLength-10, 10))
		req.Header.Set("Content-Type", "plain/text")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if expected, actual := http.StatusBadRequest, resp.StatusCode; expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
	})

	t.Run("put with repo failure", func(t *testing.T) {
		fn := guard(func(b []byte) bool {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var (
				clients      = metricMocks.NewMockGauge(ctrl)
				writtenBytes = metricMocks.NewMockCounter(ctrl)
				records      = metricMocks.NewMockCounter(ctrl)
				duration     = metricMocks.NewMockHistogramVec(ctrl)
				observer     = metricMocks.NewMockObserver(ctrl)
				repo         = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, writtenBytes, records, duration)
				server = httptest.NewServer(api)

				content, err = models.BuildContent(
					models.WithSize(int64(len(b))),
					models.WithContentBytes(b),
					models.WithContentType("plain/text"),
				)
			)
			defer api.Close()

			if err != nil {
				t.Fatal(err)
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("POST", "/", "500").Return(observer).Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			repo.EXPECT().PutContent(Content(content)).Return(models.Content{}, errors.New("failure")).Times(1)

			resp, err := http.Post(server.URL, "plain/text", bytes.NewBuffer(b))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusInternalServerError, resp.StatusCode; expected != actual {
				t.Errorf("expected: %d, actual: %d", expected, actual)
			}

			return true
		})

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("put", func(t *testing.T) {
		fn := guard(func(b []byte) bool {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var (
				clients      = metricMocks.NewMockGauge(ctrl)
				writtenBytes = metricMocks.NewMockCounter(ctrl)
				records      = metricMocks.NewMockCounter(ctrl)
				duration     = metricMocks.NewMockHistogramVec(ctrl)
				observer     = metricMocks.NewMockObserver(ctrl)
				repo         = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, writtenBytes, records, duration)
				server = httptest.NewServer(api)

				content, err = models.BuildContent(
					models.WithSize(int64(len(b))),
					models.WithContentBytes(b),
					models.WithContentType("plain/text"),
				)
			)
			defer api.Close()

			if err != nil {
				t.Fatal(err)
			}

			clients.EXPECT().Inc().Times(1)
			clients.EXPECT().Dec().Times(1)

			duration.EXPECT().WithLabelValues("POST", "/", "200").Return(observer).Times(1)
			writtenBytes.EXPECT().Add(float64(len(b))).Times(1)
			records.EXPECT().Inc().Times(1)
			observer.EXPECT().Observe(Float64()).Times(1)

			repo.EXPECT().PutContent(Content(content)).Return(content, nil).Times(1)

			resp, err := http.Post(server.URL, "plain/text", bytes.NewBuffer(b))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if expected, actual := http.StatusOK, resp.StatusCode; expected != actual {
				t.Errorf("expected: %d, actual: %d", expected, actual)
			}

			return true
		})

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
				clients      = metricMocks.NewMockGauge(ctrl)
				writtenBytes = metricMocks.NewMockCounter(ctrl)
				records      = metricMocks.NewMockCounter(ctrl)
				duration     = metricMocks.NewMockHistogramVec(ctrl)
				observer     = metricMocks.NewMockObserver(ctrl)
				repo         = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, writtenBytes, records, duration)
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

type float64Matcher struct{}

func (float64Matcher) Matches(x interface{}) bool {
	_, ok := x.(float64)
	return ok
}

func (float64Matcher) String() string {
	return "is float64"
}

func Float64() gomock.Matcher { return float64Matcher{} }

type errNotFound struct {
	err error
}

func (e errNotFound) Error() string {
	return e.err.Error()
}

func (e errNotFound) NotFound() bool {
	return true
}

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

type contentsMatcher struct {
	contents []models.Content
}

func (m contentsMatcher) Matches(x interface{}) bool {
	d, ok := x.([]models.Content)
	if !ok {
		return false
	}

	res := true
	for k, content := range m.contents {
		c := d[k]
		res = content.Address() == c.Address() &&
			content.ContentType() == c.ContentType() &&
			content.Size() == c.Size()
	}

	return res
}

func (contentsMatcher) String() string {
	return "is contents"
}

func Contents(contents []models.Content) gomock.Matcher { return contentsMatcher{contents} }

type uuidMatcher struct {
	uuid uuid.UUID
}

func (m uuidMatcher) Matches(x interface{}) bool {
	d, ok := x.(uuid.UUID)
	if !ok {
		return false
	}

	return m.uuid.Equals(d)
}

func (uuidMatcher) String() string {
	return "is uuid"
}

func UUID(uuid uuid.UUID) gomock.Matcher { return uuidMatcher{uuid} }

type queryMatcher struct{}

func (m queryMatcher) Matches(x interface{}) bool {
	_, ok := x.(repository.Query)
	return ok
}

func (queryMatcher) String() string {
	return "is query"
}

func Query() gomock.Matcher { return queryMatcher{} }
