package contents

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/quick"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
	"github.com/trussle/snowy/pkg/document"
	metricMocks "github.com/trussle/snowy/pkg/metrics/mocks"
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
				clients  = metricMocks.NewMockGauge(ctrl)
				duration = metricMocks.NewMockHistogramVec(ctrl)
				observer = metricMocks.NewMockObserver(ctrl)
				repo     = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
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
				clients  = metricMocks.NewMockGauge(ctrl)
				duration = metricMocks.NewMockHistogramVec(ctrl)
				observer = metricMocks.NewMockObserver(ctrl)
				repo     = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
				server = httptest.NewServer(api)

				content, err = document.BuildContent(
					document.WithSize(int64(len(bytes))),
					document.WithBytes(bytes),
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

			repo.EXPECT().GetContent(uid).Times(1).Return(content, nil)

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
				clients  = metricMocks.NewMockGauge(ctrl)
				duration = metricMocks.NewMockHistogramVec(ctrl)
				observer = metricMocks.NewMockObserver(ctrl)
				repo     = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
				server = httptest.NewServer(api)

				content, err = document.BuildContent(
					document.WithSize(int64(len(bytes))),
					document.WithBytes(bytes),
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

			repo.EXPECT().GetContent(uid).Times(1).Return(content, errNotFound{errors.New("failure")})

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
				clients  = metricMocks.NewMockGauge(ctrl)
				duration = metricMocks.NewMockHistogramVec(ctrl)
				observer = metricMocks.NewMockObserver(ctrl)
				repo     = repoMocks.NewMockRepository(ctrl)

				api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
				server = httptest.NewServer(api)

				content, err = document.BuildContent(
					document.WithSize(int64(len(bytes))),
					document.WithBytes(bytes),
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

			repo.EXPECT().GetContent(uid).Times(1).Return(content, errors.New("failure"))

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
