package createshorturl

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"url-shortener/internal/pkg/database/mongodb"
	"url-shortener/internal/pkg/redis"
)

type MappingStore interface {
	FindLongURL(context.Context, string) (string, error)
	CreateMapping(context.Context, string, string) error
}

type URLCache interface {
	StoreURL(context.Context, string, string) error
}

type Result struct {
	ShortCode string
	ShortURL  string
}

func CreateShortURL(ctx context.Context, longURL, baseURL string, store MappingStore, cache URLCache) (Result, error) {
	if err := validateLongURL(longURL); err != nil {
		return Result{}, err
	}
	if err := validateBaseURL(baseURL); err != nil {
		return Result{}, err
	}

	shortCode, err := store.FindLongURL(ctx, longURL)
	if err != nil && !errors.Is(err, mongodb.ErrMappingNotFound) {
		return Result{}, fmt.Errorf("check for existing URL mapping: %w", err)
	}

	if errors.Is(err, mongodb.ErrMappingNotFound) {
		shortCode, err = createMapping(ctx, store, longURL)
		if err != nil {
			return Result{}, err
		}
	}

	if err := cache.StoreURL(ctx, shortCode, longURL); err != nil {
		slog.Warn("URL mapping saved in MongoDB but not cached in Redis", "short_code", shortCode, "error", err)
	}

	return Result{
		ShortCode: shortCode,
		ShortURL:  strings.TrimRight(baseURL, "/") + "/v1/" + shortCode,
	}, nil
}

func createMapping(ctx context.Context, store MappingStore, longURL string) (string, error) {
	const attempts = 5
	for range attempts {
		shortCode, err := generateShortCode()
		if err != nil {
			return "", fmt.Errorf("generate short code: %w", err)
		}

		if err := store.CreateMapping(ctx, shortCode, longURL); err == nil {
			return shortCode, nil
		} else if !errors.Is(err, mongodb.ErrMappingConflict) {
			return "", fmt.Errorf("save URL mapping: %w", err)
		}

		existingCode, lookupErr := store.FindLongURL(ctx, longURL)
		if lookupErr == nil {
			return existingCode, nil
		}
		if !errors.Is(lookupErr, mongodb.ErrMappingNotFound) {
			return "", fmt.Errorf("save URL mapping: %w", lookupErr)
		}
	}

	return "", fmt.Errorf("could not allocate a unique short code after %d attempts: %w", attempts, mongodb.ErrMappingConflict)
}

func generateShortCode() (string, error) {
	randomBytes := make([]byte, 9)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(randomBytes), nil
}

func validateLongURL(longURL string) error {
	parsed, err := url.ParseRequestURI(longURL)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return errors.New("long_url must be an absolute HTTP or HTTPS URL")
	}
	return nil
}

func validateBaseURL(baseURL string) error {
	parsed, err := url.ParseRequestURI(baseURL)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return errors.New("public base URL must be an absolute HTTP or HTTPS URL")
	}
	return nil
}

var _ MappingStore = (*mongodb.Client)(nil)
var _ URLCache = (*redis.Client)(nil)
