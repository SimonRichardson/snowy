// +build documentation

package contents

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/SimonRichardson/betwixt"
	"github.com/SimonRichardson/betwixt/pkg/output"
	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
	"github.com/trussle/snowy/pkg/document"
	metricMocks "github.com/trussle/snowy/pkg/metrics/mocks"
	repoMocks "github.com/trussle/snowy/pkg/repository/mocks"
	"github.com/trussle/snowy/pkg/uuid"
)

func TestDocumentation_Flow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	file, err := os.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}

	var (
		clients  = metricMocks.NewMockGauge(ctrl)
		duration = metricMocks.NewMockHistogramVec(ctrl)
		observer = metricMocks.NewMockObserver(ctrl)
		repo     = repoMocks.NewMockRepository(ctrl)

		api = NewAPI(repo, log.NewNopLogger(), clients, duration)

		outputs = []betwixt.Output{
			output.NewMarkdown(file, output.Options{
				Header:    "# Snowy",
				Optionals: true,
			}),
		}
		capture = betwixt.New(api, outputs)
		server  = httptest.NewServer(capture)

		uid    = uuid.New()
		source = make([]byte, rand.Intn(100)+50)
	)
	if _, err = rand.Read(source); err != nil {
		t.Fatal(err)
	}

	var (
		base64Source = base64.URLEncoding.EncodeToString(source)
		bytes        = []byte(base64Source)
	)
	address, err := document.ContentAddress(bytes)
	if err != nil {
		t.Fatal(err)
	}

	var (
		outputContent, _ = document.BuildContent(
			document.WithAddress(address),
			document.WithContentType("application/octet-stream"),
			document.WithSize(int64(len(bytes))),
			document.WithBytes(bytes),
		)
	)

	defer func() {
		if err := capture.Output(); err != nil {
			t.Fatal(err)
		}

		file.Sync()
		file.Close()
	}()

	t.Run("get", func(t *testing.T) {
		clients.EXPECT().Inc().Times(1)
		clients.EXPECT().Dec().Times(1)

		duration.EXPECT().WithLabelValues("GET", "/", "200").Return(observer).Times(1)
		observer.EXPECT().Observe(Float64()).Times(1)

		repo.EXPECT().GetContent(uid).Times(1).Return(outputContent, nil)

		resp, err := http.Get(fmt.Sprintf("%s?resource_id=%s", server.URL, uid))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
	})
}
