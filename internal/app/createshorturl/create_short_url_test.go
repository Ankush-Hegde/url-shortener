package createshorturl

import (
	"context"
	"errors"
	"testing"

	"url-shortener/internal/pkg/database/mongodb"
)

type testStore struct {
	byLong  map[string]string
	byShort map[string]string
}

func newTestStore() *testStore {
	return &testStore{
		byLong:  make(map[string]string),
		byShort: make(map[string]string),
	}
}

func (s *testStore) FindLongURL(_ context.Context, longURL string) (string, error) {
	code, ok := s.byLong[longURL]
	if !ok {
		return "", mongodb.ErrMappingNotFound
	}
	return code, nil
}

func (s *testStore) CreateMapping(_ context.Context, shortCode, longURL string) error {
	if _, exists := s.byLong[longURL]; exists {
		return mongodb.ErrMappingConflict
	}
	s.byLong[longURL] = shortCode
	s.byShort[shortCode] = longURL
	return nil
}

type testCache struct {
	byShort map[string]string
}

func (c *testCache) StoreURL(_ context.Context, shortCode, longURL string) error {
	c.byShort[shortCode] = longURL
	return nil
}

func TestCreateShortURLPersistsAndReusesMapping(t *testing.T) {
	store := newTestStore()
	cache := &testCache{byShort: make(map[string]string)}
	ctx := context.Background()
	longURL := "https://example.com/path"

	first, err := CreateShortURL(ctx, longURL, "https://short.example", store, cache)
	if err != nil {
		t.Fatalf("CreateShortURL() error = %v", err)
	}
	if first.ShortCode == "" {
		t.Fatal("CreateShortURL() returned an empty short code")
	}
	if first.ShortURL != "https://short.example/v1/"+first.ShortCode {
		t.Errorf("ShortURL = %q, want base URL with generated code", first.ShortURL)
	}
	if store.byShort[first.ShortCode] != longURL {
		t.Errorf("Mongo mapping = %q, want %q", store.byShort[first.ShortCode], longURL)
	}
	if cache.byShort[first.ShortCode] != longURL {
		t.Errorf("Redis cache = %q, want %q", cache.byShort[first.ShortCode], longURL)
	}

	second, err := CreateShortURL(ctx, longURL, "https://short.example", store, cache)
	if err != nil {
		t.Fatalf("CreateShortURL() for existing URL error = %v", err)
	}
	if second != first {
		t.Errorf("repeated CreateShortURL() = %+v, want existing result %+v", second, first)
	}
	if len(store.byLong) != 1 {
		t.Errorf("stored %d mappings, want exactly one", len(store.byLong))
	}
}

func TestCreateShortURLRejectsNonHTTPURL(t *testing.T) {
	_, err := CreateShortURL(
		context.Background(),
		"javascript:alert(1)",
		"https://short.example",
		newTestStore(),
		&testCache{byShort: make(map[string]string)},
	)
	if err == nil {
		t.Fatal("CreateShortURL() error = nil, want invalid URL error")
	}
}

func TestCreateShortURLRejectsEmptyBaseURL(t *testing.T) {
	_, err := CreateShortURL(
		context.Background(),
		"https://example.com",
		"",
		newTestStore(),
		&testCache{byShort: make(map[string]string)},
	)
	if err == nil {
		t.Fatal("CreateShortURL() error = nil, want invalid base URL error")
	}
}

type failingStore struct {
	testStore
	err error
}

func (s failingStore) FindLongURL(context.Context, string) (string, error) {
	return "", s.err
}

func TestCreateShortURLSurfacesStoreErrors(t *testing.T) {
	wantErr := errors.New("database unavailable")
	_, err := CreateShortURL(
		context.Background(),
		"https://example.com",
		"https://short.example",
		&failingStore{err: wantErr},
		&testCache{byShort: make(map[string]string)},
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf("CreateShortURL() error = %v, want wrapped database error", err)
	}
}
