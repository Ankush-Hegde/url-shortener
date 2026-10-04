package middleware

import (
	"net/http"

	middlewareUtils "url-shortener/internal/pkg/server/middleware/utils"
)

func ValidateRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middlewareUtils.ValidateRequest(w, r, next)
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middlewareUtils.ApplyCorsPolicy(w, r, next)
	})
}

func RedirectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middlewareUtils.AddLocationHeader(w, r, next)
	})
}
