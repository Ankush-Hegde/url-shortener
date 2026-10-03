package middleware

import (
	"net/http"
	"strings"
)

func ValidateRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validateRequest(w, r, next)
	})
}

func validateRequest(w http.ResponseWriter, r *http.Request, next http.Handler) {
	// 1. Validate Method (Only allow GET and POST)
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed. Only GET and POST are supported.", http.StatusMethodNotAllowed)
		return
	}

	// 2. Validate Content-Type
	if r.Method == http.MethodPost {
		contentType := r.Header.Get("Content-Type")

		// Use HasPrefix to handle cases like "application/json; charset=utf-8"
		if !strings.HasPrefix(strings.ToLower(contentType), "application/json") {
			http.Error(w, "Unsupported Media Type. Content-Type must be application/json", http.StatusUnsupportedMediaType)
			return
		}
	}

	// If everything passes, move to the next handler/controller
	next.ServeHTTP(w, r)
}
