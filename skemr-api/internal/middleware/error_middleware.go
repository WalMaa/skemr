package middleware

import (
	"net/http"
)

type APIError struct {
	Path    string `json:"path"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// ErrorHandler captures errors and returns a consistent JSON error response
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// In Chi, error handling is typically done in handlers, but you can add global error handling here if needed.
		// For now, just call the next handler.
		next.ServeHTTP(w, r)
	})
}
