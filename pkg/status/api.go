package status

import (
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	errs "github.com/trussle/snowy/pkg/http"
)

// These are the status API URL paths.
const (
	APIPathGetQuery = "/"
)

// API serves the status API
type API struct {
	logger log.Logger
}

// NewAPI creates a API with the correct dependencies.
func NewAPI(logger log.Logger) *API {
	return &API{
		logger: logger,
	}
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	iw := &interceptingWriter{http.StatusOK, w}
	w = iw

	// Routing table
	method, path := r.Method, r.URL.Path
	switch {
	case method == "GET" && path == APIPathGetQuery:
		w.WriteHeader(http.StatusOK)

		// Handle empty ledgers
		if err := json.NewEncoder(w).Encode(struct{}{}); err != nil {
			errs.Error(a.logger, w, err.Error(), http.StatusInternalServerError)
		}
	default:
		// Nothing found
		errs.NotFound(a.logger, w, r)
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
