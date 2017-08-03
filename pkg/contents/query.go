package contents

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/trussle/snowy/pkg/document"
	"github.com/trussle/snowy/pkg/uuid"
	"github.com/pkg/errors"
)

const (
	defaultKB = 1024
	defaultMB = 1024 * defaultKB

	defaultMaxContentLength = 5 * defaultMB
)

// SelectQueryParams defines all the dimensions of a query.
type SelectQueryParams struct {
	ResourceID uuid.UUID `json:"resource_id"`
}

// DecodeFrom populates a SelectQueryParams from a URL.
func (qp *SelectQueryParams) DecodeFrom(u *url.URL, rb queryBehavior) error {
	// Required depending on the query behavior
	if rb == queryRequired {
		var (
			err        error
			resourceID = u.Query().Get("resource_id")
		)
		if resourceID == "" {
			return errors.New("error reading 'resource_id' (required) query")
		}
		if qp.ResourceID, err = uuid.Parse(resourceID); err != nil {
			return errors.Wrap(err, "error parsing 'resource_id' (required) query")
		}
	}

	return nil
}

// SelectQueryResult contains statistics about the query.
type SelectQueryResult struct {
	Params   SelectQueryParams `json:"query"`
	Duration string            `json:"duration"`
	Content  document.Content  `json:"content"`
}

// EncodeTo encodes the SelectQueryResult to the HTTP response writer.
func (qr *SelectQueryResult) EncodeTo(w http.ResponseWriter) {
	w.Header().Set(httpHeaderDuration, qr.Duration)
	w.Header().Set(httpHeaderResourceID, qr.Params.ResourceID.String())
}

// InsertQueryParams defines all the dimensions of a query.
type InsertQueryParams struct {
}

// DecodeFrom populates a InsertQueryParams from a URL.
func (qp *InsertQueryParams) DecodeFrom(u *url.URL, h http.Header, rb queryBehavior) error {
	// Required depending on the query behavior
	if rb == queryRequired {
		// Get the content-length
		contentLength := h.Get("Content-Length")
		if contentLength != "" {
			return errors.New("error reading 'content-length' (required) query")
		}

		if size, err := strconv.ParseInt(contentLength, 10, 64); err != nil {
			return errors.New("error parsing 'content-length' (required) query")
		} else if size > defaultMaxContentLength {
			return errors.Errorf("error request body too large")
		}
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
func (qr *InsertQueryResult) EncodeTo(w http.ResponseWriter) {
	w.Header().Set(httpHeaderDuration, qr.Duration)
}

const (
	httpHeaderDuration   = "X-Duration"
	httpHeaderResourceID = "X-ResourceID"
)

type queryBehavior int

const (
	queryRequired queryBehavior = iota
	queryOptional
)
