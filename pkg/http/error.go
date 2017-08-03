package http

import (
	"encoding/json"
	"net/http"
)

// Error replies to the request with the specified error message and HTTP code.
// It does not otherwise end the request; the caller should ensure no further
// writes are done to w.
// The error message should be application/json.
func Error(w http.ResponseWriter, err string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(struct {
		Description string `json:"description"`
		Code        int    `json:"code"`
	}{
		Description: err,
		Code:        code,
	}); err != nil {
		panic(err)
	}
}

// NotFound replies to the request with an HTTP 404 not found error.
func NotFound(w http.ResponseWriter, r *http.Request) {
	Error(w, "not found", http.StatusNotFound)
}

// BadRequest to the request with an HTTP 400 bad request error.
func BadRequest(w http.ResponseWriter, r *http.Request, err string) {
	Error(w, err, http.StatusBadRequest)
}

// InternalServerError to the request with an HTTP 500 bad request error.
func InternalServerError(w http.ResponseWriter, r *http.Request, err string) {
	Error(w, err, http.StatusInternalServerError)
}
