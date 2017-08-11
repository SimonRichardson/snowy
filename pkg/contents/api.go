package contents

import (
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/trussle/snowy/pkg/document"
	errs "github.com/trussle/snowy/pkg/http"
	"github.com/trussle/snowy/pkg/metrics"
	"github.com/trussle/snowy/pkg/repository"
)

// These are the query API URL paths.
const (
	APIPathGetQuery  = "/"
	APIPathPostQuery = "/"
)

// API serves the query API
type API struct {
	repository repository.Repository
	action     chan func()
	stop       chan chan struct{}
	logger     log.Logger
	clients    metrics.Gauge
	duration   metrics.HistogramVec
}

// NewAPI creates a API with correct dependencies.
func NewAPI(repository repository.Repository, logger log.Logger,
	clients metrics.Gauge,
	duration metrics.HistogramVec,
) *API {
	api := &API{
		repository: repository,
		action:     make(chan func()),
		stop:       make(chan chan struct{}),
		logger:     logger,
		clients:    clients,
		duration:   duration,
	}
	go api.run()
	return api
}

// Close out the API
func (a *API) Close() {
	c := make(chan struct{})
	a.stop <- c
	<-c
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
		a.handlePost(w, r)
	default:
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

	var (
		notFound      = make(chan struct{})
		internalError = make(chan error)
		result        = make(chan document.Content)
	)
	a.action <- func() {
		content, err := a.repository.GetContent(qp.ResourceID)
		if err != nil {
			if repository.ErrNotFound(err) {
				notFound <- struct{}{}
				return
			}
			internalError <- err
			return
		}
		result <- content
	}

	select {
	case <-notFound:
		errs.Error(w, "not found", http.StatusNotFound)
	case err := <-internalError:
		errs.Error(w, err.Error(), http.StatusInternalServerError)
	case content := <-result:
		// Make sure we collect the document for the result.
		qr := SelectQueryResult{Params: qp}
		qr.Content = content

		// Finish
		qr.Duration = time.Since(begin).String()
		qr.EncodeTo(w)
	}
}

func (a *API) handlePost(w http.ResponseWriter, r *http.Request) {
	// useful metrics
	begin := time.Now()

	defer r.Body.Close()

	// Validate user input.
	var qp InsertQueryParams
	if err := qp.DecodeFrom(r.URL, r.Header, queryRequired); err != nil {
		errs.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var (
		notFound        = make(chan struct{})
		internalError   = make(chan error)
		badRequestError = make(chan error)
		result          = make(chan document.Content)
	)
	a.action <- func() {
		body := io.LimitReader(r.Body, defaultMaxContentLength)
		bytes, err := ioutil.ReadAll(body)
		if err != nil {
			badRequestError <- err
			return
		}

		address, err := document.ContentAddress(bytes)
		if err != nil {
			internalError <- err
			return
		}

		content, err := document.BuildContent(
			document.WithAddress(address),
			document.WithSize(int64(len(bytes))),
			document.WithContentType(qp.ContentType),
			document.WithBytes(bytes),
		)
		if err != nil {
			badRequestError <- err
			return
		}

		res, err := a.repository.PutContent(content)
		if err != nil {
			if repository.ErrNotFound(err) {
				notFound <- struct{}{}
				return
			}
			internalError <- err
			return
		}
		result <- res
	}

	select {
	case err := <-internalError:
		errs.Error(w, err.Error(), http.StatusInternalServerError)
	case err := <-badRequestError:
		errs.Error(w, err.Error(), http.StatusBadRequest)
	case content := <-result:
		// Make sure we collect the document for the result.
		qr := InsertQueryResult{Params: qp}
		qr.Content = content

		// Finish
		qr.Duration = time.Since(begin).String()
		qr.EncodeTo(w)
	default:
		errs.Error(w, "unknown error", http.StatusInternalServerError)
	}
}

func (a *API) run() {
	for {
		select {
		case f := <-a.action:
			f()

		case c := <-a.stop:
			close(c)
			return
		}
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
