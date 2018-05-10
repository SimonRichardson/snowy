package journals

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	errs "github.com/trussle/snowy/pkg/http"
	"github.com/trussle/snowy/pkg/metrics"
	"github.com/trussle/snowy/pkg/models"
	"github.com/trussle/snowy/pkg/repository"
)

// These are the query API URL paths.
const (
	APIPathInsertQuery = "/"
	APIPathAppendQuery = "/"
)

const (
	contentFormFile  = "content"
	documentFormFile = "document"
)

// API serves the query API
type API struct {
	handler        http.Handler
	repository     repository.Repository
	logger         log.Logger
	clients        metrics.Gauge
	bytes, records metrics.Counter
	duration       metrics.HistogramVec
	errors         errs.Error
}

// NewAPI creates a API with correct dependencies.
func NewAPI(repository repository.Repository, logger log.Logger,
	clients metrics.Gauge,
	bytes, records metrics.Counter,
	duration metrics.HistogramVec,
) *API {
	api := &API{
		repository: repository,
		logger:     logger,
		clients:    clients,
		bytes:      bytes,
		records:    records,
		duration:   duration,
		errors:     errs.NewError(logger),
	}
	{
		router := mux.NewRouter().StrictSlash(true)
		router.Methods("POST").Path(APIPathInsertQuery).HandlerFunc(api.handleInsert)
		router.Methods("PUT").Path(APIPathAppendQuery).HandlerFunc(api.handleAppend)
		router.NotFoundHandler = http.HandlerFunc(api.errors.NotFound)

		api.handler = router
	}
	return api
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	iw := &interceptingWriter{http.StatusOK, w}
	w = iw

	// Metrics
	a.clients.Inc()
	defer a.clients.Dec()

	defer func(begin time.Time) {
		a.duration.WithLabelValues(
			r.Method,
			r.URL.Path,
			strconv.Itoa(iw.code),
		).Observe(time.Since(begin).Seconds())
	}(time.Now())

	a.handler.ServeHTTP(w, r)
}

func (a *API) handleInsert(w http.ResponseWriter, r *http.Request) {
	// useful metrics
	begin := time.Now()

	defer r.Body.Close()

	// Validate user input.
	var qp InsertQueryParams
	if err := qp.DecodeFrom(r.URL, r.Header, queryRequired); err != nil {
		a.errors.BadRequest(w, r, err.Error())
		return
	}

	var (
		internalError   = make(chan error)
		badRequestError = make(chan error)
		result          = make(chan models.Ledger)
	)
	go func() {
		file, fileHeader, err := r.FormFile(contentFormFile)
		if err != nil {
			badRequestError <- err
			return
		}

		var fqp InsertFileQueryParams
		if err = fqp.DecodeFrom(fileHeader.Header, queryRequired); err != nil {
			badRequestError <- err
			return
		}

		document, documentHeader, err := r.FormFile(documentFormFile)
		if err != nil {
			badRequestError <- err
			return
		}

		var lqp InsertFileQueryParams
		if err = lqp.DecodeFrom(documentHeader.Header, queryRequired); err != nil {
			badRequestError <- err
			return
		}

		content, contentLength, err := ingestContent(file, fqp)
		if err != nil {
			badRequestError <- err
			return
		}

		ledger, err := ingestLedger(document, content, lqp, func() models.DocOption {
			return models.WithNewResourceID()
		})
		if err != nil {
			badRequestError <- err
			return
		}

		if _, err = a.repository.PutContent(content); err != nil {
			internalError <- err
			return
		}

		ledgerResult, err := a.repository.InsertLedger(ledger)
		if err != nil {
			internalError <- err
			return
		}

		a.bytes.Add(float64(contentLength))
		a.records.Inc()

		result <- ledgerResult
	}()

	select {
	case err := <-internalError:
		a.errors.Error(w, err.Error(), http.StatusInternalServerError)
	case err := <-badRequestError:
		a.errors.Error(w, err.Error(), http.StatusBadRequest)
	case resource := <-result:
		// Make sure we collect the content for the result.
		qr := InsertQueryResult{Params: qp}
		qr.ResourceID = resource.ResourceID()

		// Finish
		qr.Duration = time.Since(begin).String()
		qr.EncodeTo(w)
	}
}

func (a *API) handleAppend(w http.ResponseWriter, r *http.Request) {
	// useful metrics
	begin := time.Now()

	defer r.Body.Close()

	// Validate user input.
	var qp AppendQueryParams
	if err := qp.DecodeFrom(r.URL, r.Header, queryRequired); err != nil {
		a.errors.BadRequest(w, r, err.Error())
		return
	}

	var (
		internalError   = make(chan error)
		badRequestError = make(chan error)
		result          = make(chan models.Ledger)
	)
	go func() {
		file, fileHeader, err := r.FormFile(contentFormFile)
		if err != nil {
			badRequestError <- err
			return
		}

		var fqp InsertFileQueryParams
		if err = fqp.DecodeFrom(fileHeader.Header, queryRequired); err != nil {
			badRequestError <- err
			return
		}

		document, documentHeader, err := r.FormFile(documentFormFile)
		if err != nil {
			badRequestError <- err
			return
		}

		var lqp InsertFileQueryParams
		if err = lqp.DecodeFrom(documentHeader.Header, queryRequired); err != nil {
			badRequestError <- err
			return
		}

		content, contentLength, err := ingestContent(file, fqp)
		if err != nil {
			badRequestError <- err
			return
		}

		ledger, err := ingestLedger(document, content, lqp, func() models.DocOption {
			return models.WithNewResourceID()
		})
		if err != nil {
			badRequestError <- err
			return
		}

		if _, err = a.repository.PutContent(content); err != nil {
			internalError <- err
			return
		}

		ledgerResult, err := a.repository.AppendLedger(qp.ResourceID, ledger)
		if err != nil {
			internalError <- err
			return
		}

		a.bytes.Add(float64(contentLength))
		a.records.Inc()

		result <- ledgerResult
	}()

	select {
	case err := <-internalError:
		a.errors.Error(w, err.Error(), http.StatusInternalServerError)
	case err := <-badRequestError:
		a.errors.Error(w, err.Error(), http.StatusBadRequest)
	case resource := <-result:
		// Make sure we collect the content for the result.
		qr := AppendQueryResult{Params: qp}
		qr.ResourceID = resource.ResourceID()

		// Finish
		qr.Duration = time.Since(begin).String()
		qr.EncodeTo(w)
	}
}

type interceptingWriter struct {
	code int
	http.ResponseWriter
}

func (iw *interceptingWriter) WriteHeader(code int) {
	iw.code = code
	iw.ResponseWriter.WriteHeader(code)
}

// ContentHeader returns a header with both type and length
type ContentHeader interface {
	ContentType() string
	ContentLength() int64
}

func ingestContent(reader io.Reader, header ContentHeader) (models.Content, int, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, header.ContentLength()))
	if _, err := buffer.ReadFrom(reader); err != nil {
		return models.Content{}, -1, err
	}

	bytes := buffer.Bytes()
	if len(bytes) < 1 {
		return models.Content{}, -1, errors.New("no body content")
	}

	content, err := models.BuildContent(
		models.WithContentBytes(bytes),
		models.WithSize(int64(len(bytes))),
		models.WithContentType(header.ContentType()),
	)
	if err != nil {
		return content, -1, err
	}
	return content, len(bytes), nil
}

func ingestLedger(reader io.ReadCloser, content models.Content, header ContentHeader, fn func() models.DocOption) (models.Ledger, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, header.ContentLength()))
	if _, err := buffer.ReadFrom(reader); err != nil {
		return models.Ledger{}, err
	}

	bytes := buffer.Bytes()
	if len(bytes) < 1 {
		return models.Ledger{}, errors.New("no body content")
	}

	var input models.LedgerInput
	if err := json.Unmarshal(bytes, &input); err != nil {
		return models.Ledger{}, err
	}
	if err := models.ValidateLedgerInput(input); err != nil {
		return models.Ledger{}, err
	}

	return models.BuildLedger(
		fn(),
		models.WithName(input.Name),
		models.WithResourceAddress(content.Address()),
		models.WithResourceSize(content.Size()),
		models.WithResourceContentType(content.ContentType()),
		models.WithAuthorID(input.AuthorID),
		models.WithTags(input.Tags),
		models.WithCreatedOn(time.Now()),
	)
}
