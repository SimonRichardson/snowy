package contents

import (
	"bytes"
	"net/http"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	errs "github.com/trussle/snowy/pkg/http"
	"github.com/trussle/snowy/pkg/metrics"
	"github.com/trussle/snowy/pkg/models"
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
	case (method == "PUT" || method == "POST") && path == APIPathPostQuery:
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
		result        = make(chan models.Content)
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
		// Make sure we collect the content for the result.
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

	// This shouldn't mutate the state :(
	r.Body = http.MaxBytesReader(w, r.Body, defaultMaxContentLength)
	defer r.Body.Close()

	// Validate user input.
	var qp InsertQueryParams
	if err := qp.DecodeFrom(r.URL, r.Header, queryRequired); err != nil {
		errs.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var (
		internalError   = make(chan error)
		badRequestError = make(chan error)
		result          = make(chan models.Content)
	)
	a.action <- func() {
		buffer := bytes.NewBuffer(make([]byte, 0, defaultMaxContentLength))
		if _, err := buffer.ReadFrom(r.Body); err != nil {
			badRequestError <- err
			return
		}

		bytes := buffer.Bytes()
		content, err := models.BuildContent(
			models.WithContentBytes(bytes),
			models.WithSize(int64(len(bytes))),
			models.WithContentType(qp.ContentType),
		)
		if err != nil {
			badRequestError <- err
			return
		}

		res, err := a.repository.PutContent(content)
		if err != nil {
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
		// Make sure we collect the content for the result.
		qr := InsertQueryResult{Params: qp}
		qr.Content = content

		// Finish
		qr.Duration = time.Since(begin).String()
		qr.EncodeTo(w)
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
