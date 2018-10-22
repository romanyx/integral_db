package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	internalServerErrorMessage = "internal server error"
)

// OK response.
func OK(w http.ResponseWriter, resp interface{}) {
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		InternalServerError(w)
	}
}

// BadRequest response.
func BadRequest(w http.ResponseWriter, resp interface{}) {
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		InternalServerError(w)
	}
}

// NotFound response.
func NotFound(w http.ResponseWriter, resp interface{}) {
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		InternalServerError(w)
	}
}

// InternalServerError response.
func InternalServerError(w http.ResponseWriter) {
	http.Error(w, fmt.Sprintf("{\"message\": \"%s\"}", internalServerErrorMessage), http.StatusInternalServerError)
}
