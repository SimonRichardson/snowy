package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// Error replies to the request with the specified error message and HTTP code.
// It does not otherwise end the request; the caller should ensure no further
// writes are done to w.
// The error message should be application/json.
func Error(logger log.Logger, w http.ResponseWriter, err string, code int) {
	level.Error(logger).Log("err", err, "code", code)

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
func NotFound(logger log.Logger, w http.ResponseWriter, r *http.Request) {
	Error(logger, w, "not found", http.StatusNotFound)
}

// BadRequest to the request with an HTTP 400 bad request error.
func BadRequest(logger log.Logger, w http.ResponseWriter, r *http.Request, err string) {
	Error(logger, w, err, http.StatusBadRequest)
}

// InternalServerError to the request with an HTTP 500 bad request error.
func InternalServerError(logger log.Logger, w http.ResponseWriter, r *http.Request, err string) {
	Error(logger, w, err, http.StatusInternalServerError)
}
