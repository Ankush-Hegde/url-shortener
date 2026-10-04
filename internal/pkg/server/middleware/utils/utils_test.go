package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddLocationHeaderSetsLocationForFoundResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusFound)
		_, _ = w.Write([]byte(`{"url":"https://example.com/path"}`))
	})
	response := httptest.NewRecorder()

	AddLocationHeader(response, httptest.NewRequest(http.MethodGet, "/short", nil), handler)

	if response.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusFound)
	}
	if got := response.Header().Get("Location"); got != "https://example.com/path" {
		t.Errorf("Location = %q, want %q", got, "https://example.com/path")
	}
	if got := response.Body.String(); got != `{"url":"https://example.com/path"}` {
		t.Errorf("body = %q, want redirect response body unchanged", got)
	}
	if got := response.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
}

func TestAddLocationHeaderPassesThroughNonFoundResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Test", "value")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("response"))
	})
	response := httptest.NewRecorder()

	AddLocationHeader(response, httptest.NewRequest(http.MethodGet, "/", nil), handler)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if got := response.Header().Get("Location"); got != "" {
		t.Errorf("Location = %q, want empty", got)
	}
	if got := response.Header().Get("X-Test"); got != "value" {
		t.Errorf("X-Test = %q, want value", got)
	}
	if got := response.Body.String(); got != "response" {
		t.Errorf("body = %q, want response", got)
	}
}

func TestAddLocationHeaderReturnsErrorForInvalidFoundResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusFound)
		_, _ = w.Write([]byte(`{"url":""}`))
	})
	response := httptest.NewRecorder()

	AddLocationHeader(response, httptest.NewRequest(http.MethodGet, "/short", nil), handler)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusInternalServerError)
	}
}
