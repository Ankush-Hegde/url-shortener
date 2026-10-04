package resolveshorturl

import (
	"context"
	"errors"
	"testing"

	"url-shortener/internal/pkg/database/mongodb"
	"url-shortener/internal/pkg/redis"
)

type testStore struct {
	longURL string
	err     error
	calls   int
}

func (s *testStore) GetRedirectUrl(context.Context, string) (string, error) {
	s.calls++
	return s.longURL, s.err
}

type testCache struct {
	longURL   string
	getErr    error
	storeErr  error
	storeCall int
}

func (c *testCache) GetURL(context.Context, string) (string, error) {
	return c.longURL, c.getErr
}

func (c *testCache) StoreURL(_ context.Context, _ string, longURL string) error {
	c.storeCall++
	c.longURL = longURL
	return c.storeErr
}

func TestResolveShortURLUsesRedisHit(t *testing.T) {
	store := &testStore{}
	cache := &testCache{longURL: "https://cached.example"}

	got, err := ResolveShortURL(context.Background(), "abc123", store, cache)
	if err != nil {
		t.Fatalf("ResolveShortURL() error = %v", err)
	}
	if got != cache.longURL {
		t.Errorf("ResolveShortURL() = %q, want %q", got, cache.longURL)
	}
	if store.calls != 0 {
		t.Errorf("MongoDB lookup calls = %d, want 0", store.calls)
	}
}

func TestResolveShortURLLoadsMongoAndCachesMiss(t *testing.T) {
	store := &testStore{longURL: "https://database.example"}
	cache := &testCache{getErr: redis.ErrURLNotFound}

	got, err := ResolveShortURL(context.Background(), "abc123", store, cache)
	if err != nil {
		t.Fatalf("ResolveShortURL() error = %v", err)
	}
	if got != store.longURL {
		t.Errorf("ResolveShortURL() = %q, want %q", got, store.longURL)
	}
	if cache.storeCall != 1 || cache.longURL != store.longURL {
		t.Errorf("Redis cache was not populated: calls=%d value=%q", cache.storeCall, cache.longURL)
	}
}

func TestResolveShortURLFallsBackWhenRedisUnavailable(t *testing.T) {
	store := &testStore{longURL: "https://database.example"}
	cache := &testCache{getErr: errors.New("Redis unavailable")}

	got, err := ResolveShortURL(context.Background(), "abc123", store, cache)
	if err != nil {
		t.Fatalf("ResolveShortURL() error = %v", err)
	}
	if got != store.longURL {
		t.Errorf("ResolveShortURL() = %q, want %q", got, store.longURL)
	}
}

func TestResolveShortURLReturnsNotFound(t *testing.T) {
	store := &testStore{err: mongodb.ErrMappingNotFound}
	cache := &testCache{getErr: redis.ErrURLNotFound}

	_, err := ResolveShortURL(context.Background(), "missing", store, cache)
	if !errors.Is(err, ErrShortURLNotFound) {
		t.Fatalf("ResolveShortURL() error = %v, want ErrShortURLNotFound", err)
	}
}
