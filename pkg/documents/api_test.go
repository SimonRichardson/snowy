package documents

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/trussle/snowy/pkg/document"
	"github.com/trussle/snowy/pkg/metrics"
	"github.com/trussle/snowy/pkg/repository"
	"github.com/trussle/snowy/pkg/uuid"
	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
)

func TestAPI(t *testing.T) {
	t.Parallel()

	t.Run("get with no id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			clients  = metrics.NewMockGauge(ctrl)
			duration = metrics.NewMockHistogramVec(ctrl)
			observer = metrics.NewMockObserver(ctrl)
			repo     = repository.NewMockRepository(ctrl)

			api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
			server = httptest.NewServer(api)
		)

		clients.EXPECT().Inc().Times(1)
		clients.EXPECT().Dec().Times(1)

		duration.EXPECT().WithLabelValues("GET", "/", "400").Return(observer).Times(1)
		observer.EXPECT().Observe(Float64()).Times(1)

		_, err := http.Get(server.URL)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("get with resource_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			clients  = metrics.NewMockGauge(ctrl)
			duration = metrics.NewMockHistogramVec(ctrl)
			observer = metrics.NewMockObserver(ctrl)
			repo     = repository.NewMockRepository(ctrl)

			api    = NewAPI(repo, log.NewNopLogger(), clients, duration)
			server = httptest.NewServer(api)

			uid = uuid.New()

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

		repo.EXPECT().GetDocument(uid).Times(1).Return(doc, nil)

		resp, err := http.Get(fmt.Sprintf("%s?resource_id=%s", server.URL, uid))
		if err != nil {
			t.Error(err)
		}

		if expected, actual := http.StatusOK, resp.StatusCode; expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
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
