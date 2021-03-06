package contents

import (
	"bytes"
	"io"
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
	APIPathMultipleQuery        = "/multiple/"
	APIPathSelectRevisionsQuery = "/revisions/"
)

// API serves the query API
type API struct {
	handler        http.Handler
	repository     repository.Repository
	action         chan func()
	stop           chan chan struct{}
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
		action:     make(chan func()),
		stop:       make(chan chan struct{}),
		logger:     logger,
		clients:    clients,
		bytes:      bytes,
		records:    records,
		duration:   duration,
		errors:     errs.NewError(logger),
	}
	{
		router := mux.NewRouter().StrictSlash(true)
		router.Methods("GET").Path(APIPathSelectQuery).HandlerFunc(api.handleSelect)
		router.Methods("PUT").Path(APIPathInsertQuery).HandlerFunc(api.handleInsert)
		router.Methods("POST").Path(APIPathInsertQuery).HandlerFunc(api.handleInsert)
		router.Methods("GET").Path(APIPathMultipleQuery).HandlerFunc(api.handleMultiple)
		router.Methods("GET").Path(APIPathSelectRevisionsQuery).HandlerFunc(api.handleSelectRevisions)
		router.NotFoundHandler = http.HandlerFunc(api.errors.NotFound)

		api.handler = router
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
		a.errors.Error(w, err.Error(), http.StatusBadRequest)
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

	var (
		notFound      = make(chan struct{})
		internalError = make(chan error)
		result        = make(chan models.Content)
	)
	go func() {
		content, err := a.repository.SelectContent(qp.ResourceID, options)
		if err != nil {
			if repository.ErrNotFound(err) {
				notFound <- struct{}{}
				return
			}
			internalError <- err
			return
		}
		result <- content
	}()

	select {
	case <-notFound:
		a.errors.Error(w, "not found", http.StatusNotFound)
	case err := <-internalError:
		a.errors.Error(w, err.Error(), http.StatusInternalServerError)
	case content := <-result:
		// Make sure we collect the content for the result.
		qr := SelectQueryResult{Errors: a.errors, Params: qp}
		qr.Content = content

		// Finish
		qr.Duration = time.Since(begin).String()
		qr.EncodeTo(w)
	}
}

func (a *API) handleInsert(w http.ResponseWriter, r *http.Request) {
	// useful metrics
	begin := time.Now()

	// This shouldn't mutate the state :(
	r.Body = http.MaxBytesReader(w, r.Body, defaultMaxContentLength)
	defer r.Body.Close()

	// Validate user input.
	var qp InsertQueryParams
	if err := qp.DecodeFrom(r.URL, r.Header, queryRequired); err != nil {
		a.errors.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var (
		internalError   = make(chan error)
		badRequestError = make(chan error)
		result          = make(chan models.Content)
	)
	a.action <- func() {
		content, contentLength, err := ingestContent(r.Body, qp)
		if err != nil {
			badRequestError <- err
			return
		}

		res, err := a.repository.PutContent(content)
		if err != nil {
			internalError <- err
			return
		}

		a.bytes.Add(float64(contentLength))
		a.records.Inc()

		result <- res
	}

	select {
	case err := <-internalError:
		a.errors.Error(w, err.Error(), http.StatusInternalServerError)
	case err := <-badRequestError:
		a.errors.Error(w, err.Error(), http.StatusBadRequest)
	case content := <-result:
		// Make sure we collect the content for the result.
		qr := InsertQueryResult{Errors: a.errors, Params: qp}
		qr.Content = content

		// Finish
		qr.Duration = time.Since(begin).String()
		qr.EncodeTo(w)
	}
}

func (a *API) handleMultiple(w http.ResponseWriter, r *http.Request) {
	// useful metrics
	begin := time.Now()

	defer r.Body.Close()

	// Validate user input.
	var qp MultipleQueryParams
	if err := qp.DecodeFrom(r.URL, queryRequired); err != nil {
		a.errors.BadRequest(w, r, err.Error())
		return
	}

	var (
		internalError = make(chan error)
		result        = make(chan []models.Content)
	)
	go func() {
		contents := make([]models.Content, len(qp.ResourceIDs))
		for k, v := range qp.ResourceIDs {
			c, err := a.repository.SelectContent(v, repository.BuildEmptyQuery())
			if err != nil {
				internalError <- err
				return
			}
			contents[k] = c
		}
		result <- contents
	}()

	select {
	case err := <-internalError:
		a.errors.Error(w, err.Error(), http.StatusInternalServerError)
	case contents := <-result:
		// Make sure we collect the content for the result.
		qr := MultipleQueryResult{Errors: a.errors, Params: qp}
		qr.Contents = contents

		// Finish
		qr.Duration = time.Since(begin).String()
		qr.EncodeTo(w)
	}
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

	var (
		internalError = make(chan error)
		result        = make(chan []models.Content)
	)
	go func() {
		contents, err := a.repository.SelectContents(qp.ResourceID, options)
		if err != nil {
			internalError <- err
			return
		}
		result <- contents
	}()

	select {
	case err := <-internalError:
		a.errors.Error(w, err.Error(), http.StatusInternalServerError)
	case contents := <-result:
		// Make sure we collect the content for the result.
		qr := SelectRevisionsQueryResult{Errors: a.errors, Params: qp}
		qr.Contents = contents

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

// ContentHeader returns a header with both type and length
type ContentHeader interface {
	ContentType() string
	ContentLength() int64
}

func ingestContent(file io.Reader, header ContentHeader) (models.Content, int, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, header.ContentLength()))
	if _, err := buffer.ReadFrom(file); err != nil {
		return models.Content{}, -1, err
	}

	bytes := buffer.Bytes()
	content, err := models.BuildContent(
		models.WithContentBytes(bytes),
		models.WithSize(int64(len(bytes))),
		models.WithContentType(header.ContentType()),
	)
	if err != nil {
		return content, -1, err
	}
	return content, len(bytes), err
}
