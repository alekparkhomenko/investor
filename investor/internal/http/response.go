package http

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse is the standard error response body.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// jsonResponse writes a JSON response with the given status code and data.
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

// errorResponse writes a JSON error response with the given status code.
func errorResponse(w http.ResponseWriter, status int, code, message string) {
	jsonResponse(w, status, ErrorResponse{
		Error:   code,
		Message: message,
		Code:    status,
	})
}
