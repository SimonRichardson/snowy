package documents

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/trussle/snowy/pkg/document"
	errs "github.com/trussle/snowy/pkg/http"
	"github.com/trussle/snowy/pkg/uuid"
	"github.com/pkg/errors"
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
	Document document.Document `json:"document"`
}

// EncodeTo encodes the SelectQueryResult to the HTTP response writer.
func (qr *SelectQueryResult) EncodeTo(w http.ResponseWriter) {
	w.Header().Set(httpHeaderDuration, qr.Duration)
	w.Header().Set(httpHeaderResourceID, qr.Params.ResourceID.String())

	if err := json.NewEncoder(w).Encode(qr.Document); err != nil {
		errs.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// InsertQueryParams defines all the dimensions of a query.
type InsertQueryParams struct {
}

// DecodeFrom populates a InsertQueryParams from a URL.
func (qp *InsertQueryParams) DecodeFrom(u *url.URL, rb queryBehavior) error {
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

	if err := json.NewEncoder(w).Encode(struct {
		ResourceID uuid.UUID `json:"resource_id"`
	}{
		ResourceID: qr.ResourceID,
	}); err != nil {
		errs.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
