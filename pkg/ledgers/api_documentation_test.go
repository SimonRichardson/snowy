// +build documentation

package ledgers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/SimonRichardson/betwixt"
	"github.com/SimonRichardson/betwixt/pkg/output"
	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
	metricMocks "github.com/trussle/snowy/pkg/metrics/mocks"
	"github.com/trussle/snowy/pkg/models"
	"github.com/trussle/snowy/pkg/repository"
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

		uid         = uuid.New()
		forkUID     = uuid.New()
		tags        = []string{"abc", "def", "g"}
		inputDoc, _ = models.BuildLedger(
			models.WithAuthorID(uid.String()),
			models.WithResourceAddress("abcdefghij"),
			models.WithResourceSize(10),
			models.WithResourceContentType("application/octet-stream"),
			models.WithName("document-name"),
			models.WithTags(tags),
		)
		outputDoc, _ = models.BuildLedger(
			models.WithResourceID(uid),
			models.WithResourceAddress("abcdefghij"),
			models.WithResourceSize(10),
			models.WithResourceContentType("application/octet-stream"),
			models.WithAuthorID(uuid.New().String()),
			models.WithName("document-name"),
			models.WithTags(tags),
			models.WithCreatedOn(time.Now()),
			models.WithDeletedOn(time.Time{}),
		)
		outputForkDoc, _ = models.BuildLedger(
			models.WithResourceID(forkUID),
			models.WithResourceAddress("abcdefghij"),
			models.WithResourceSize(10),
			models.WithResourceContentType("application/octet-stream"),
			models.WithAuthorID(uuid.New().String()),
			models.WithName("document-name"),
			models.WithTags(tags),
			models.WithCreatedOn(time.Now()),
			models.WithDeletedOn(time.Time{}),
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

		query, _ := repository.BuildQuery(
			repository.WithQueryTags(tags),
			repository.WithQueryAuthorID(""),
		)

		repo.EXPECT().SelectLedger(uid, query).Times(1).Return(outputDoc, nil)

		resp, err := http.Get(fmt.Sprintf("%s/?resource_id=%s&query.tags=%s", server.URL, uid, strings.Join(tags, ",")))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
	})

	t.Run("get revisions", func(t *testing.T) {
		clients.EXPECT().Inc().Times(1)
		clients.EXPECT().Dec().Times(1)

		duration.EXPECT().WithLabelValues("GET", "/revisions/", "200").Return(observer).Times(1)
		observer.EXPECT().Observe(Float64()).Times(1)

		query, _ := repository.BuildQuery(
			repository.WithQueryTags(tags),
			repository.WithQueryAuthorID(""),
		)

		repo.EXPECT().SelectLedgers(uid, query).Times(1).Return([]models.Ledger{
			outputDoc,
		}, nil)

		resp, err := http.Get(fmt.Sprintf("%s/revisions/?resource_id=%s&query.tags=%s", server.URL, uid, strings.Join(tags, ",")))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
	})

	t.Run("insert", func(t *testing.T) {
		clients.EXPECT().Inc().Times(1)
		clients.EXPECT().Dec().Times(1)

		duration.EXPECT().WithLabelValues("POST", "/", "200").Return(observer).Times(1)
		observer.EXPECT().Observe(Float64()).Times(1)

		repo.EXPECT().InsertLedger(Ledger(inputDoc)).Times(1).Return(outputDoc, nil)

		b, err := json.Marshal(struct {
			Name     string   `json:"name"`
			AuthorID string   `json:"author_id"`
			Tags     []string `json:"tags"`
		}{
			Name:     "document-name",
			AuthorID: uid.String(),
			Tags:     []string{"abc", "def", "g"},
		})
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.Post(server.URL, "application/json", bytes.NewReader(b))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
	})

	t.Run("append", func(t *testing.T) {
		clients.EXPECT().Inc().Times(1)
		clients.EXPECT().Dec().Times(1)

		duration.EXPECT().WithLabelValues("PUT", "/", "200").Return(observer).Times(1)
		observer.EXPECT().Observe(Float64()).Times(1)

		repo.EXPECT().AppendLedger(uid, Ledger(inputDoc)).Return(outputDoc, nil).Times(1)

		b, err := json.Marshal(struct {
			Name     string   `json:"name"`
			AuthorID string   `json:"author_id"`
			Tags     []string `json:"tags"`
		}{
			Name:     "document-name",
			AuthorID: uid.String(),
			Tags:     []string{"abc", "def", "g"},
		})
		if err != nil {
			t.Fatal(err)
		}

		resp, err := Put(fmt.Sprintf("%s?resource_id=%s", server.URL, uid), "application/json", bytes.NewReader(b))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
	})

	t.Run("fork", func(t *testing.T) {
		clients.EXPECT().Inc().Times(1)
		clients.EXPECT().Dec().Times(1)

		duration.EXPECT().WithLabelValues("PUT", "/fork/", "200").Return(observer).Times(1)
		observer.EXPECT().Observe(Float64()).Times(1)

		repo.EXPECT().ForkLedger(uid, Ledger(inputDoc)).Return(outputForkDoc, nil).Times(1)

		b, err := json.Marshal(struct {
			Name     string   `json:"name"`
			AuthorID string   `json:"author_id"`
			Tags     []string `json:"tags"`
		}{
			Name:     "document-name",
			AuthorID: uid.String(),
			Tags:     []string{"abc", "def", "g"},
		})
		if err != nil {
			t.Fatal(err)
		}

		resp, err := Put(fmt.Sprintf("%s/fork/?resource_id=%s", server.URL, uid), "application/json", bytes.NewReader(b))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
	})
}
