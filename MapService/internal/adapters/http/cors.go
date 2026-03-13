package http

import (
	"net/http"

	"github.com/go-chi/cors"
)

// NewCORSMiddleware returns a chi-compatible CORS middleware that only allows
// the specified origins. allowedOrigins must not be empty; the caller is
// responsible for validating the list before calling this function.
func NewCORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
