package status

import (
	"net/http"
	"strings"

	errs "github.com/trussle/snowy/pkg/http"
)

// These are the status API URL paths.
const (
	APIPathGetQuery = "/"
)

// API serves the status API
type API struct{}

// NewAPI creates a API with the correct dependencies.
func NewAPI() *API {
	return &API{}
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	iw := &interceptingWriter{http.StatusOK, w}
	w = iw

	// Routing table
	method, path := r.Method, r.URL.Path
	switch {
	case method == "GET" && path == APIPathGetQuery:
		w.WriteHeader(http.StatusOK)
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

type interceptingWriter struct {
	code int
	http.ResponseWriter
}

func (iw *interceptingWriter) WriteHeader(code int) {
	iw.code = code
	iw.ResponseWriter.WriteHeader(code)
}
