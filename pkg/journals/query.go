package journals

import (
	"encoding/json"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	errs "github.com/trussle/snowy/pkg/http"
	"github.com/trussle/snowy/pkg/uuid"
)

const (
	defaultContentType = "multipart/form-data"
)

const (
	defaultKB = 1024
	defaultMB = 1024 * defaultKB

	defaultMaxContentLength = 10 * defaultMB
)

// InsertQueryParams defines all the dimensions of a query.
type InsertQueryParams struct {
}

// DecodeFrom populates a InsertQueryParams from a Request.
func (qp *InsertQueryParams) DecodeFrom(u *url.URL, h http.Header, rb queryBehavior) error {
	if contentType := h.Get("Content-Type"); rb == queryRequired && !strings.Contains(strings.ToLower(contentType), "multipart/form-data") {
		return errors.Errorf("expected 'multipart/form-data' content-type, got %q", contentType)
	}

	return nil
}

// InsertQueryResult contains statistics about the query.
type InsertQueryResult struct {
	Params     InsertQueryParams `json:"query"`
	Duration   string            `json:"duration"`
	ResourceID uuid.UUID         `json:"resource_id"`
}

// EncodeTo encodes the InsertQueryResult to the HTTP response writer.
func (qr *InsertQueryResult) EncodeTo(logger log.Logger, w http.ResponseWriter) {
	w.Header().Set(httpHeaderContentType, defaultContentType)
	w.Header().Set(httpHeaderDuration, qr.Duration)

	if err := json.NewEncoder(w).Encode(struct {
		ResourceID uuid.UUID `json:"resource_id"`
	}{
		ResourceID: qr.ResourceID,
	}); err != nil {
		errs.Error(logger, w, err.Error(), http.StatusInternalServerError)
	}
}

// InsertFileQueryParams defines all the dimensions of a query.
type InsertFileQueryParams struct {
	contentType   string
	contentLength int64
}

// DecodeFrom populates a InsertFileQueryParams from a URL.
func (qp *InsertFileQueryParams) DecodeFrom(h textproto.MIMEHeader, rb queryBehavior) error {
	// Required depending on the query behavior
	if rb == queryRequired {
		// Get the content-type
		if qp.contentType = h.Get("Content-Type"); qp.contentType == "" {
			return errors.New("error reading 'content-type' (required) query")
		}

		// Get the content-length
		contentLength := h.Get("Content-Length")
		if contentLength == "" {
			return errors.New("error reading 'content-length' (required) query")
		}

		size, err := strconv.ParseInt(contentLength, 10, 64)
		if err != nil {
			return errors.New("error parsing 'content-length' (required) query")
		} else if size > defaultMaxContentLength {
			return errors.Errorf("error request body too large")
		} else if size < 1 {
			return errors.Errorf("error request body is empty")
		}

		qp.contentLength = size
	}

	return nil
}

// ContentType returns the content-type from the header
func (qp InsertFileQueryParams) ContentType() string { return qp.contentType }

// ContentLength returns the content-length from the header
func (qp InsertFileQueryParams) ContentLength() int64 { return qp.contentLength }

// AppendQueryParams defines all the dimensions of a query.
type AppendQueryParams struct {
	ResourceID uuid.UUID `json:"resource_id"`
}

// DecodeFrom populates a AppendQueryParams from a URL.
func (qp *AppendQueryParams) DecodeFrom(u *url.URL, h http.Header, rb queryBehavior) error {
	// Required depending on the query behavior
	if contentType := h.Get("Content-Type"); rb == queryRequired && !strings.Contains(strings.ToLower(contentType), "multipart/form-data") {
		return errors.Errorf("expected 'multipart/form-data' content-type, got %q", contentType)
	}

	var (
		err        error
		resourceID = u.Query().Get("resource_id")
	)
	if rb == queryRequired && resourceID == "" {
		return errors.New("error reading 'resource_id' (required) query")
	}
	if resourceID != "" {
		if qp.ResourceID, err = uuid.Parse(resourceID); err != nil {
			return errors.Wrap(err, "error parsing 'resource_id' (required) query")
		}
	}

	return nil
}

// AppendQueryResult contains statistics about the query.
type AppendQueryResult struct {
	Params     AppendQueryParams `json:"query"`
	Duration   string            `json:"duration"`
	ResourceID uuid.UUID         `json:"resource_id"`
}

// EncodeTo encodes the AppendQueryResult to the HTTP response writer.
func (qr *AppendQueryResult) EncodeTo(logger log.Logger, w http.ResponseWriter) {
	w.Header().Set(httpHeaderContentType, defaultContentType)
	w.Header().Set(httpHeaderDuration, qr.Duration)
	w.Header().Set(httpHeaderResourceID, qr.Params.ResourceID.String())

	if err := json.NewEncoder(w).Encode(struct {
		ResourceID uuid.UUID `json:"resource_id"`
	}{
		ResourceID: qr.ResourceID,
	}); err != nil {
		errs.Error(logger, w, err.Error(), http.StatusInternalServerError)
	}
}

const (
	httpHeaderContentType   = "Content-Type"
	httpHeaderDuration      = "X-Duration"
	httpHeaderResourceID    = "X-Resource-ID"
	httpHeaderQueryTags     = "X-Query-Tags"
	httpHeaderQueryAuthorID = "X-Query-Author-ID"
)

type queryBehavior int

const (
	queryRequired queryBehavior = iota
	queryOptional
)
