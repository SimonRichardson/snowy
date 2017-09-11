// +build documentation

package journals

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/SimonRichardson/betwixt"
	"github.com/SimonRichardson/betwixt/pkg/output"
	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
	metricMocks "github.com/trussle/snowy/pkg/metrics/mocks"
	"github.com/trussle/snowy/pkg/models"
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
		clients      = metricMocks.NewMockGauge(ctrl)
		writtenBytes = metricMocks.NewMockCounter(ctrl)
		records      = metricMocks.NewMockCounter(ctrl)
		duration     = metricMocks.NewMockHistogramVec(ctrl)
		observer     = metricMocks.NewMockObserver(ctrl)
		repo         = repoMocks.NewMockRepository(ctrl)

		api = NewAPI(repo, log.NewNopLogger(), clients, writtenBytes, records, duration)

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
		conBytes     = []byte(base64Source)
	)
	address, err := models.ContentAddress(conBytes)
	if err != nil {
		t.Fatal(err)
	}

	var (
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
		inputContent, _ = models.BuildContent(
			models.WithSize(int64(len(conBytes))),
			models.WithContentBytes(conBytes),
			models.WithContentType("application/octet-stream"),
		)
		outputContent, _ = models.BuildContent(
			models.WithAddress(address),
			models.WithContentType("application/octet-stream"),
			models.WithSize(int64(len(conBytes))),
			models.WithBytes(conBytes),
		)
	)

	defer func() {
		if err := capture.Output(); err != nil {
			t.Fatal(err)
		}

		file.Sync()
		file.Close()
	}()

	t.Run("insert", func(t *testing.T) {
		clients.EXPECT().Inc().Times(1)
		clients.EXPECT().Dec().Times(1)

		duration.EXPECT().WithLabelValues("POST", "/", "200").Return(observer).Times(1)
		writtenBytes.EXPECT().Add(float64(len(conBytes))).Times(1)
		records.EXPECT().Inc().Times(1)
		observer.EXPECT().Observe(Float64()).Times(1)

		repo.EXPECT().PutContent(Content(inputContent)).Return(outputContent, nil).Times(1)
		repo.EXPECT().InsertLedger(Ledger(inputDoc)).Times(1).Return(outputDoc, nil)

		docBytes, err := json.Marshal(struct {
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
	})

	t.Run("append", func(t *testing.T) {
		clients.EXPECT().Inc().Times(1)
		clients.EXPECT().Dec().Times(1)

		duration.EXPECT().WithLabelValues("PUT", "/", "200").Return(observer).Times(1)
		writtenBytes.EXPECT().Add(float64(len(conBytes))).Times(1)
		records.EXPECT().Inc().Times(1)
		observer.EXPECT().Observe(Float64()).Times(1)

		repo.EXPECT().PutContent(Content(inputContent)).Return(outputContent, nil).Times(1)
		repo.EXPECT().AppendLedger(uid, Ledger(inputDoc)).Return(outputDoc, nil).Times(1)

		docBytes, err := json.Marshal(struct {
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

		var (
			buffer bytes.Buffer
			writer = multipart.NewWriter(&buffer)
		)

		MustWriteField(writer, contentFormFile, "application/octet-stream", conBytes)
		MustWriteField(writer, documentFormFile, "application/json", docBytes)

		if err := writer.Close(); err != nil {
			t.Fatal(err)
		}

		resp, err := Put(fmt.Sprintf("%s?resource_id=%s", server.URL, uid), writer.FormDataContentType(), &buffer)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
	})
}
