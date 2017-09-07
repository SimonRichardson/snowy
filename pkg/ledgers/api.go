package ledgers

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	errs "github.com/trussle/snowy/pkg/http"
	"github.com/trussle/snowy/pkg/metrics"
	"github.com/trussle/snowy/pkg/models"
	"github.com/trussle/snowy/pkg/repository"
)

// These are the query API URL paths.
const (
	APIPathGetQuery         = "/"
	APIPathPostQuery        = "/"
	APIPathPutQuery         = "/"
	APIPathGetMultipleQuery = "/multiple/"
)

// API serves the query API
type API struct {
	repository repository.Repository
	logger     log.Logger
	clients    metrics.Gauge
	duration   metrics.HistogramVec
}

// NewAPI creates a API with correct dependencies.
func NewAPI(repository repository.Repository, logger log.Logger,
	clients metrics.Gauge,
	duration metrics.HistogramVec,
) *API {
	return &API{
		repository: repository,
		logger:     logger,
		clients:    clients,
		duration:   duration,
	}
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	level.Info(a.logger).Log("url", r.URL.String())

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

	// Routing table
	method, path := r.Method, r.URL.Path
	switch {
	case method == "GET" && path == APIPathGetQuery:
		a.handleGet(w, r)
	case method == "POST" && path == APIPathPostQuery:
		a.handlePost(w, r)
	case method == "PUT" && path == APIPathPutQuery:
		a.handlePut(w, r)
	case method == "GET" && path == APIPathGetMultipleQuery:
		a.handleGetMultiple(w, r)
	default:
		// Nothing found
		errs.NotFound(a.logger, w, r)
	}
}

func (a *API) handleGet(w http.ResponseWriter, r *http.Request) {
	// useful metrics
	begin := time.Now()

	defer r.Body.Close()

	// Validate user input.
	var qp SelectQueryParams
	if err := qp.DecodeFrom(r.URL, queryRequired); err != nil {
		errs.BadRequest(a.logger, w, r, err.Error())
		return
	}

	options, err := repository.BuildQuery(
		repository.WithQueryTags(qp.Tags),
		repository.WithQueryAuthorID(qp.AuthorID),
	)
	if err != nil {
		errs.BadRequest(a.logger, w, r, err.Error())
		return
	}

	doc, err := a.repository.GetLedger(qp.ResourceID, options)
	if err != nil {
		if repository.ErrNotFound(err) {
			errs.NotFound(a.logger, w, r)
			return
		}
		errs.InternalServerError(a.logger, w, r, err.Error())
		return
	}

	// Make sure we collect the document for the result.
	qr := SelectQueryResult{Logger: a.logger, Params: qp}
	qr.Ledger = doc

	// Finish
	qr.Duration = time.Since(begin).String()
	qr.EncodeTo(w)
}

func (a *API) handlePost(w http.ResponseWriter, r *http.Request) {
	// useful metrics
	begin := time.Now()

	defer r.Body.Close()

	// Validate user input.
	var qp InsertQueryParams
	if err := qp.DecodeFrom(r.URL, r.Header, queryRequired); err != nil {
		errs.BadRequest(a.logger, w, r, err.Error())
		return
	}

	doc, err := ingestLedger(r.Body, func() models.DocOption {
		return models.WithNewResourceID()
	})
	if err != nil {
		errs.BadRequest(a.logger, w, r, err.Error())
		return
	}

	resource, err := a.repository.InsertLedger(doc)
	if err != nil {
		errs.InternalServerError(a.logger, w, r, err.Error())
		return
	}

	// Make sure we collect the document for the result.
	qr := InsertQueryResult{Logger: a.logger, Params: qp}
	qr.ResourceID = resource.ResourceID()

	// Finish
	qr.Duration = time.Since(begin).String()
	qr.EncodeTo(w)
}

func (a *API) handlePut(w http.ResponseWriter, r *http.Request) {
	// useful metrics
	begin := time.Now()

	defer r.Body.Close()

	// Validate user input.
	var qp AppendQueryParams
	if err := qp.DecodeFrom(r.URL, r.Header, queryRequired); err != nil {
		errs.BadRequest(a.logger, w, r, err.Error())
		return
	}

	doc, err := ingestLedger(r.Body, func() models.DocOption {
		return models.WithResourceID(qp.ResourceID)
	})
	if err != nil {
		errs.BadRequest(a.logger, w, r, err.Error())
		return
	}

	resource, err := a.repository.AppendLedger(qp.ResourceID, doc)
	if err != nil {
		errs.InternalServerError(a.logger, w, r, err.Error())
		return
	}

	// Make sure we collect the document for the result.
	qr := AppendQueryResult{Logger: a.logger, Params: qp}
	qr.ResourceID = resource.ID()

	// Finish
	qr.Duration = time.Since(begin).String()
	qr.EncodeTo(w)
}

func (a *API) handleGetMultiple(w http.ResponseWriter, r *http.Request) {
	// useful metrics
	begin := time.Now()

	defer r.Body.Close()

	// Validate user input.
	var qp SelectQueryParams
	if err := qp.DecodeFrom(r.URL, queryRequired); err != nil {
		errs.BadRequest(a.logger, w, r, err.Error())
		return
	}

	options, err := repository.BuildQuery(
		repository.WithQueryTags(qp.Tags),
		repository.WithQueryAuthorID(qp.AuthorID),
	)
	if err != nil {
		errs.BadRequest(a.logger, w, r, err.Error())
		return
	}

	ledgers, err := a.repository.GetLedgers(qp.ResourceID, options)
	if err != nil {
		errs.InternalServerError(a.logger, w, r, err.Error())
		return
	}

	// Make sure we collect the documents for the result.
	qr := SelectMultipleQueryResult{Logger: a.logger, Params: qp}
	qr.Ledgers = ledgers

	// Finish
	qr.Duration = time.Since(begin).String()
	qr.EncodeTo(w)
}

type interceptingWriter struct {
	code int
	http.ResponseWriter
}

func (iw *interceptingWriter) WriteHeader(code int) {
	iw.code = code
	iw.ResponseWriter.WriteHeader(code)
}

func ingestLedger(reader io.ReadCloser, fn func() models.DocOption) (models.Ledger, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return models.Ledger{}, err
	}

	if len(bytes) < 1 {
		return models.Ledger{}, errors.New("no body content")
	}

	var input models.LedgerInput
	if err = json.Unmarshal(bytes, &input); err != nil {
		return models.Ledger{}, err
	}
	if err = models.ValidateLedgerInput(input); err != nil {
		return models.Ledger{}, err
	}

	return models.BuildLedger(
		fn(),
		models.WithName(input.Name),
		models.WithResourceAddress(input.ResourceAddress),
		models.WithResourceSize(input.ResourceSize),
		models.WithResourceContentType(input.ResourceContentType),
		models.WithAuthorID(input.AuthorID),
		models.WithTags(input.Tags),
		models.WithCreatedOn(time.Now()),
	)
}
