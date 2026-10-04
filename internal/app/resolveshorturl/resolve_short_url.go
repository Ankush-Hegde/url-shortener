package resolveshorturl

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"url-shortener/internal/pkg/database/mongodb"
	"url-shortener/internal/pkg/redis"
)

type MappingStore interface {
	FindShortCode(context.Context, string) (string, error)
}

type URLCache interface {
	GetURL(context.Context, string) (string, error)
	StoreURL(context.Context, string, string) error
}

var ErrShortURLNotFound = errors.New("short URL not found")

func ResolveShortURL(ctx context.Context, shortCode string, store MappingStore, cache URLCache) (string, error) {
	if strings.TrimSpace(shortCode) == "" {
		return "", errors.New("short code is required")
	}

	longURL, cacheErr := cache.GetURL(ctx, shortCode)
	if cacheErr == nil {
		return longURL, nil
	}
	if !errors.Is(cacheErr, redis.ErrURLNotFound) {
		slog.Warn("Redis lookup failed; falling back to MongoDB", "short_code", shortCode, "error", cacheErr)
	}

	longURL, err := store.FindShortCode(ctx, shortCode)
	if errors.Is(err, mongodb.ErrMappingNotFound) {
		return "", ErrShortURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("find URL mapping in MongoDB: %w", err)
	}

	if err := cache.StoreURL(ctx, shortCode, longURL); err != nil {
		slog.Warn("URL mapping loaded from MongoDB but not cached in Redis", "short_code", shortCode, "error", err)
	}
	return longURL, nil
}

var _ MappingStore = (*mongodb.Client)(nil)
var _ URLCache = (*redis.Client)(nil)
