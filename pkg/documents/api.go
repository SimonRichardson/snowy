package documents

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	errs "github.com/trussle/snowy/pkg/http"
	"github.com/trussle/snowy/pkg/metrics"
	"github.com/trussle/snowy/pkg/repository"
	"github.com/go-kit/kit/log"
)

// These are the query API URL paths.
const (
	APIPathGetQuery  = "/"
	APIPathPostQuery = "/"
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
		a.handlePut(w, r)
	default:
		// Make sure we send a permanent redirect if it ends with a `/`
		if strings.HasSuffix(path, "/") {
			http.Redirect(w, r, strings.TrimRight(path, "/"), http.StatusPermanentRedirect)
			return
		}
		// Nothing found
		errs.NotFound(w, r)
	}
}

func (a *API) handleGet(w http.ResponseWriter, r *http.Request) {
	// useful metrics
	begin := time.Now()

	defer r.Body.Close()

	// Validate user input.
	var qp SelectQueryParams
	if err := qp.DecodeFrom(r.URL, queryRequired); err != nil {
		errs.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	doc, err := a.repository.GetDocument(qp.ResourceID)
	if err != nil {
		if repository.ErrNotFound(err) {
			errs.NotFound(w, r)
			return
		}
		errs.InternalServerError(w, r, err.Error())
		return
	}

	// Make sure we collect the document for the result.
	qr := SelectQueryResult{Params: qp}
	qr.Document = doc

	// Finish
	qr.Duration = time.Since(begin).String()
	qr.EncodeTo(w)
}

func (a *API) handlePut(w http.ResponseWriter, r *http.Request) {
	// useful metrics
	begin := time.Now()

	defer r.Body.Close()

	// Validate user input.
	var qp InsertQueryParams
	if err := qp.DecodeFrom(r.URL, queryRequired); err != nil {
		errs.BadRequest(w, r, err.Error())
		return
	}

	resourceID, err := a.repository.PutDocument(nil)
	if err != nil {
		errs.InternalServerError(w, r, err.Error())
		return
	}
	// Make sure we collect the document for the result.
	qr := InsertQueryResult{Params: qp}
	qr.ResourceID = resourceID

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
