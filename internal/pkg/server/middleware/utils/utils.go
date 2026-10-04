package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
)

var trustedOrigins = []string{"http://localhost:3000", "http://localhost:8089"}

func ApplyCorsPolicy(w http.ResponseWriter, r *http.Request, next http.Handler) {
	origin := r.Header.Get("Origin")

	if slices.Contains(trustedOrigins, origin) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Vary", "Origin")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	next.ServeHTTP(w, r)
}

func ValidateRequest(w http.ResponseWriter, r *http.Request, next http.Handler) {
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

func AddLocationHeader(w http.ResponseWriter, r *http.Request, next http.Handler) {
	response := httptest.NewRecorder()
	next.ServeHTTP(response, r)

	for name, values := range response.Header() {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	if response.Code == http.StatusFound {
		var body struct {
			URL string `json:"url"`
		}
		if err := json.Unmarshal(response.Body.Bytes(), &body); err == nil || body.URL != "" {
			w.Header().Set("Location", body.URL)
		}

		w.WriteHeader(response.Code)
		return
	}

	w.WriteHeader(response.Code)
	_, _ = w.Write(response.Body.Bytes())
}
