package ledgers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	errs "github.com/trussle/snowy/pkg/http"
	"github.com/trussle/snowy/pkg/models"
	"github.com/trussle/snowy/pkg/uuid"
)

const (
	defaultContentType = "application/json"
)

// SelectQueryParams defines all the dimensions of a query.
type SelectQueryParams struct {
	ResourceID uuid.UUID `json:"resource_id"`
	Tags       []string  `json:"query.tags"`
	AuthorID   string    `json:"query.author_id"`
}

// DecodeFrom populates a SelectQueryParams from a URL.
func (qp *SelectQueryParams) DecodeFrom(u *url.URL, rb queryBehavior) error {
	// Required depending on the query behavior
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

	// Tags are optional here.
	tags := u.Query().Get("query.tags")
	if tags != "" {
		qp.Tags = strings.Split(tags, ",")
	}

	// Author ID is optional here.
	if authorID := u.Query().Get("query.author_id"); authorID != "" {
		qp.AuthorID = authorID
	}

	return nil
}

// SelectQueryResult contains statistics about the query.
type SelectQueryResult struct {
	Errors   errs.Error
	Params   SelectQueryParams `json:"query"`
	Duration string            `json:"duration"`
	Ledger   models.Ledger     `json:"ledger"`
}

// EncodeTo encodes the SelectQueryResult to the HTTP response writer.
func (qr *SelectQueryResult) EncodeTo(w http.ResponseWriter) {
	w.Header().Set(httpHeaderContentType, defaultContentType)
	w.Header().Set(httpHeaderDuration, qr.Duration)
	w.Header().Set(httpHeaderResourceID, qr.Params.ResourceID.String())
	w.Header().Set(httpHeaderQueryTags, strings.Join(qr.Params.Tags, ","))
	w.Header().Set(httpHeaderQueryAuthorID, qr.Params.AuthorID)

	// Handle empty ledgers
	if err := json.NewEncoder(w).Encode(qr.Ledger); err != nil {
		qr.Errors.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// InsertQueryParams defines all the dimensions of a query.
type InsertQueryParams struct {
}

// DecodeFrom populates a InsertQueryParams from a URL.
func (qp *InsertQueryParams) DecodeFrom(u *url.URL, h http.Header, rb queryBehavior) error {
	if contentType := h.Get("Content-Type"); rb == queryRequired && strings.ToLower(contentType) != "application/json" {
		return errors.Errorf("expected 'application/json' content-type, got %q", contentType)
	}

	return nil
}

// InsertQueryResult contains statistics about the query.
type InsertQueryResult struct {
	Errors     errs.Error
	Params     InsertQueryParams `json:"query"`
	Duration   string            `json:"duration"`
	ResourceID uuid.UUID         `json:"resource_id"`
}

// EncodeTo encodes the InsertQueryResult to the HTTP response writer.
func (qr *InsertQueryResult) EncodeTo(w http.ResponseWriter) {
	w.Header().Set(httpHeaderContentType, defaultContentType)
	w.Header().Set(httpHeaderDuration, qr.Duration)

	if err := json.NewEncoder(w).Encode(struct {
		ResourceID uuid.UUID `json:"resource_id"`
	}{
		ResourceID: qr.ResourceID,
	}); err != nil {
		qr.Errors.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// SelectMultipleQueryResult contains statistics about the query.
type SelectMultipleQueryResult struct {
	Errors   errs.Error
	Params   SelectQueryParams `json:"query"`
	Duration string            `json:"duration"`
	Ledgers  []models.Ledger   `json:"ledger"`
}

// EncodeTo encodes the SelectMultipleQueryResult to the HTTP response writer.
func (qr *SelectMultipleQueryResult) EncodeTo(w http.ResponseWriter) {
	w.Header().Set(httpHeaderContentType, defaultContentType)
	w.Header().Set(httpHeaderDuration, qr.Duration)
	w.Header().Set(httpHeaderResourceID, qr.Params.ResourceID.String())
	w.Header().Set(httpHeaderQueryTags, strings.Join(qr.Params.Tags, ","))
	w.Header().Set(httpHeaderQueryAuthorID, qr.Params.AuthorID)

	// Make sure that we encode empty ledgers correctly (i.e. they're not
	// null in the json output)
	docs := qr.Ledgers
	if qr.Ledgers == nil {
		docs = make([]models.Ledger, 0)
	}

	if err := json.NewEncoder(w).Encode(docs); err != nil {
		qr.Errors.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AppendQueryParams defines all the dimensions of a query.
type AppendQueryParams struct {
	ResourceID uuid.UUID `json:"resource_id"`
}

// DecodeFrom populates a AppendQueryParams from a URL.
func (qp *AppendQueryParams) DecodeFrom(u *url.URL, h http.Header, rb queryBehavior) error {
	// Required depending on the query behavior
	if contentType := h.Get("Content-Type"); rb == queryRequired && strings.ToLower(contentType) != "application/json" {
		return errors.Errorf("expected 'application/json' content-type, got %q", contentType)
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
	Errors     errs.Error
	Params     AppendQueryParams `json:"query"`
	Duration   string            `json:"duration"`
	ResourceID uuid.UUID         `json:"resource_id"`
}

// EncodeTo encodes the AppendQueryResult to the HTTP response writer.
func (qr *AppendQueryResult) EncodeTo(w http.ResponseWriter) {
	w.Header().Set(httpHeaderContentType, defaultContentType)
	w.Header().Set(httpHeaderDuration, qr.Duration)
	w.Header().Set(httpHeaderResourceID, qr.Params.ResourceID.String())

	if err := json.NewEncoder(w).Encode(struct {
		ResourceID uuid.UUID `json:"resource_id"`
	}{
		ResourceID: qr.ResourceID,
	}); err != nil {
		qr.Errors.Error(w, err.Error(), http.StatusInternalServerError)
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
