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
	"github.com/gorilla/mux"
	errs "github.com/trussle/snowy/pkg/http"
	"github.com/trussle/snowy/pkg/metrics"
	"github.com/trussle/snowy/pkg/models"
	"github.com/trussle/snowy/pkg/repository"
)

// These are the query API URL paths.
const (
	APIPathSelectQuery          = "/"
	APIPathInsertQuery          = "/"
	APIPathAppendQuery          = "/"
	APIPathSelectRevisionsQuery = "/revisions/"
	APIPathForkQuery            = "/fork/"
	APIPathForkRevisionsQuery   = "/fork/revisions/"
)

// API serves the query API
type API struct {
	handler    http.Handler
	repository repository.Repository
	logger     log.Logger
	clients    metrics.Gauge
	duration   metrics.HistogramVec
	errors     errs.Error
}

// NewAPI creates a API with correct dependencies.
func NewAPI(repository repository.Repository, logger log.Logger,
	clients metrics.Gauge,
	duration metrics.HistogramVec,
) *API {
	api := &API{
		repository: repository,
		logger:     logger,
		clients:    clients,
		duration:   duration,
		errors:     errs.NewError(logger),
	}
	{
		router := mux.NewRouter().StrictSlash(true)
		router.Methods("GET").Path(APIPathSelectQuery).HandlerFunc(api.handleSelect)
		router.Methods("POST").Path(APIPathInsertQuery).HandlerFunc(api.handleInsert)
		router.Methods("PUT").Path(APIPathAppendQuery).HandlerFunc(api.handleAppend)
		router.Methods("GET").Path(APIPathSelectRevisionsQuery).HandlerFunc(api.handleSelectRevisions)
		router.Methods("PUT").Path(APIPathForkQuery).HandlerFunc(api.handleFork)
		router.Methods("GET").Path(APIPathForkRevisionsQuery).HandlerFunc(api.handleForkRevisions)
		router.NotFoundHandler = http.HandlerFunc(api.errors.NotFound)

		api.handler = router
	}
	return api
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

	a.handler.ServeHTTP(w, r)
}

func (a *API) handleSelect(w http.ResponseWriter, r *http.Request) {
	// useful metrics
	begin := time.Now()

	defer r.Body.Close()

	// Validate user input.
	var qp SelectQueryParams
	if err := qp.DecodeFrom(r.URL, queryRequired); err != nil {
		a.errors.BadRequest(w, r, err.Error())
		return
	}

	options, err := repository.BuildQuery(
		repository.WithQueryTags(qp.Tags),
		repository.WithQueryAuthorID(qp.AuthorID),
	)
	if err != nil {
		a.errors.BadRequest(w, r, err.Error())
		return
	}

	doc, err := a.repository.SelectLedger(qp.ResourceID, options)
	if err != nil {
		if repository.ErrNotFound(err) {
			a.errors.NotFound(w, r)
			return
		}
		a.errors.InternalServerError(w, r, err.Error())
		return
	}

	// Make sure we collect the document for the result.
	qr := SelectQueryResult{Errors: a.errors, Params: qp}
	qr.Ledger = doc

	// Finish
	qr.Duration = time.Since(begin).String()
	qr.EncodeTo(w)
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

	doc, err := ingestLedger(r.Body, func() models.DocOption {
		return models.WithNewResourceID()
	})
	if err != nil {
		a.errors.BadRequest(w, r, err.Error())
		return
	}

	resource, err := a.repository.InsertLedger(doc)
	if err != nil {
		a.errors.InternalServerError(w, r, err.Error())
		return
	}

	// Make sure we collect the document for the result.
	qr := InsertQueryResult{Errors: a.errors, Params: qp}
	qr.ResourceID = resource.ResourceID()

	// Finish
	qr.Duration = time.Since(begin).String()
	qr.EncodeTo(w)
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

	doc, err := ingestLedger(r.Body, func() models.DocOption {
		return models.WithResourceID(qp.ResourceID)
	})
	if err != nil {
		a.errors.BadRequest(w, r, err.Error())
		return
	}

	resource, err := a.repository.AppendLedger(qp.ResourceID, doc)
	if err != nil {
		a.errors.InternalServerError(w, r, err.Error())
		return
	}

	// Make sure we collect the document for the result.
	qr := AppendQueryResult{Errors: a.errors, Params: qp}
	qr.ResourceID = resource.ResourceID()

	// Finish
	qr.Duration = time.Since(begin).String()
	qr.EncodeTo(w)
}

func (a *API) handleFork(w http.ResponseWriter, r *http.Request) {
	// useful metrics
	begin := time.Now()

	defer r.Body.Close()

	// Validate user input.
	var qp ForkQueryParams
	if err := qp.DecodeFrom(r.URL, r.Header, queryRequired); err != nil {
		a.errors.BadRequest(w, r, err.Error())
		return
	}

	doc, err := ingestLedger(r.Body, func() models.DocOption {
		return models.WithNewResourceID()
	})
	if err != nil {
		a.errors.BadRequest(w, r, err.Error())
		return
	}

	resource, err := a.repository.ForkLedger(qp.ResourceID, doc)
	if err != nil {
		a.errors.InternalServerError(w, r, err.Error())
		return
	}

	// Make sure we collect the document for the result.
	qr := ForkQueryResult{Errors: a.errors, Params: qp}
	qr.ResourceID = resource.ResourceID()

	// Finish
	qr.Duration = time.Since(begin).String()
	qr.EncodeTo(w)
}

func (a *API) handleSelectRevisions(w http.ResponseWriter, r *http.Request) {
	// useful metrics
	begin := time.Now()

	defer r.Body.Close()

	// Validate user input.
	var qp SelectQueryParams
	if err := qp.DecodeFrom(r.URL, queryRequired); err != nil {
		a.errors.BadRequest(w, r, err.Error())
		return
	}

	options, err := repository.BuildQuery(
		repository.WithQueryTags(qp.Tags),
		repository.WithQueryAuthorID(qp.AuthorID),
	)
	if err != nil {
		a.errors.BadRequest(w, r, err.Error())
		return
	}

	ledgers, err := a.repository.SelectLedgers(qp.ResourceID, options)
	if err != nil {
		a.errors.InternalServerError(w, r, err.Error())
		return
	}

	// Make sure we collect the documents for the result.
	qr := SelectRevisionsQueryResult{Errors: a.errors, Params: qp}
	qr.Ledgers = ledgers

	// Finish
	qr.Duration = time.Since(begin).String()
	qr.EncodeTo(w)
}

func (a *API) handleForkRevisions(w http.ResponseWriter, r *http.Request) {
	// useful metrics
	begin := time.Now()

	defer r.Body.Close()

	// Validate user input.
	var qp SelectQueryParams
	if err := qp.DecodeFrom(r.URL, queryRequired); err != nil {
		a.errors.BadRequest(w, r, err.Error())
		return
	}

	ledgers, err := a.repository.SelectForkLedgers(qp.ResourceID)
	if err != nil {
		a.errors.InternalServerError(w, r, err.Error())
		return
	}

	// Make sure we collect the documents for the result.
	qr := SelectRevisionsQueryResult{Errors: a.errors, Params: qp}
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
