package contents

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/trussle/snowy/pkg/document"
	errs "github.com/trussle/snowy/pkg/http"
	"github.com/trussle/snowy/pkg/repository"
	"github.com/trussle/snowy/pkg/uuid"
	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

// These are the query API URL paths.
const (
	APIPathGetQuery = "/"
	APIPathPutQuery = "/"
)

// API serves the query API
type API struct {
	repository repository.Repository
	action     chan func()
	stop       chan chan struct{}
	logger     log.Logger
	clients    prometheus.Gauge
	duration   *prometheus.HistogramVec
}

// NewAPI creates a API with correct dependencies.
func NewAPI(repository repository.Repository, logger log.Logger,
	clients prometheus.Gauge,
	duration *prometheus.HistogramVec,
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
	case method == "POST" && path == APIPathPutQuery:
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
		errs.NotFound(w, r)
	case err := <-internalError:
		errs.Error(w, err.Error(), http.StatusInternalServerError)
	case <-result:
		fmt.Println(begin)
		// TODO: Implement stuff
	default:
		errs.Error(w, "unknown error", http.StatusInternalServerError)
	}
}

func (a *API) handlePut(w http.ResponseWriter, r *http.Request) {
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
		notFound      = make(chan struct{})
		internalError = make(chan error)
		result        = make(chan uuid.UUID)
	)
	a.action <- func() {
		body := io.LimitReader(r.Body, defaultMaxContentLength)
		bytes, err := ioutil.ReadAll(body)
		if err != nil {
			internalError <- err
		}

		fmt.Println(bytes)

		res, err := a.repository.PutContent(nil)
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
	case resourceID := <-result:
		// Make sure we collect the document for the result.
		qr := InsertQueryResult{Params: qp}
		qr.ResourceID = resourceID

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
