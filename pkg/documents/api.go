package documents

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/trussle/snowy/pkg/document"
	errs "github.com/trussle/snowy/pkg/http"
	"github.com/trussle/snowy/pkg/metrics"
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
		errs.BadRequest(w, r, err.Error())
		return
	}

	options, err := repository.BuildQuery(
		repository.WithQueryTags(qp.Tags),
		repository.WithQueryAuthorID(qp.AuthorID),
	)
	if err != nil {
		errs.BadRequest(w, r, err.Error())
		return
	}

	doc, err := a.repository.GetDocument(qp.ResourceID, options)
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

func (a *API) handlePost(w http.ResponseWriter, r *http.Request) {
	// useful metrics
	begin := time.Now()

	defer r.Body.Close()

	// Validate user input.
	var qp InsertQueryParams
	if err := qp.DecodeFrom(r.URL, r.Header, queryRequired); err != nil {
		errs.BadRequest(w, r, err.Error())
		return
	}

	doc, err := ingestRequest(r.Body, func() document.DocOption {
		return document.WithNewResourceID()
	})
	if err != nil {
		errs.BadRequest(w, r, err.Error())
		return
	}

	resource, err := a.repository.InsertDocument(doc)
	if err != nil {
		errs.InternalServerError(w, r, err.Error())
		return
	}

	// Make sure we collect the document for the result.
	qr := InsertQueryResult{Params: qp}
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
		errs.BadRequest(w, r, err.Error())
		return
	}

	doc, err := ingestRequest(r.Body, func() document.DocOption {
		return document.WithResourceID(qp.ResourceID)
	})
	if err != nil {
		errs.BadRequest(w, r, err.Error())
		return
	}

	resource, err := a.repository.AppendDocument(qp.ResourceID, doc)
	if err != nil {
		errs.InternalServerError(w, r, err.Error())
		return
	}

	// Make sure we collect the document for the result.
	qr := AppendQueryResult{Params: qp}
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
		errs.BadRequest(w, r, err.Error())
		return
	}

	options, err := repository.BuildQuery(
		repository.WithQueryTags(qp.Tags),
		repository.WithQueryAuthorID(qp.AuthorID),
	)
	if err != nil {
		errs.BadRequest(w, r, err.Error())
		return
	}

	docs, err := a.repository.GetDocuments(qp.ResourceID, options)
	if err != nil {
		errs.InternalServerError(w, r, err.Error())
		return
	}

	// Make sure we collect the documents for the result.
	qr := SelectMultipleQueryResult{Params: qp}
	qr.Documents = docs

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

func ingestRequest(reader io.ReadCloser, fn func() document.DocOption) (document.Document, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return document.Document{}, err
	}

	if len(bytes) < 1 {
		return document.Document{}, errors.New("no body content")
	}

	var input documentInput
	if err = json.Unmarshal(bytes, &input); err != nil {
		return document.Document{}, err
	}
	if err = validateInput(input); err != nil {
		return document.Document{}, err
	}

	return document.BuildDocument(
		fn(),
		document.WithName(input.Name),
		document.WithResourceAddress(input.ResourceAddress),
		document.WithResourceSize(input.ResourceSize),
		document.WithResourceContentType(input.ResourceContentType),
		document.WithAuthorID(input.AuthorID),
		document.WithTags(input.Tags),
		document.WithCreatedOn(time.Now()),
	)
}

type documentInput struct {
	Name                string   `json:"name"`
	ResourceAddress     string   `json:"resource_address"`
	ResourceSize        int64    `json:"resource_size"`
	ResourceContentType string   `json:"resource_content_type"`
	AuthorID            string   `json:"author_id"`
	Tags                []string `json:"tags"`
}

func validateInput(input documentInput) error {
	if len(strings.TrimSpace(input.Name)) == 0 {
		return errors.New("input.name is empty")
	}

	if len(strings.TrimSpace(input.AuthorID)) == 0 {
		return errors.New("input.author_id is empty")
	}

	return nil
}
