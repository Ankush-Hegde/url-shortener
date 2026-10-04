package redis

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"strings"
	"time"

	redisclient "github.com/redis/go-redis/v9"
)

const URLTTL = 30 * 24 * time.Hour

var ErrURLNotFound = errors.New("short URL not found")

type Options struct {
	Addr     string
	Username string
	Password string
	DB       int
	TLS      bool
}

type Client struct {
	client *redisclient.Client
}

func NewClient(ctx context.Context, options Options) (*Client, error) {
	if strings.TrimSpace(options.Addr) == "" {
		return nil, errors.New("Redis address is required")
	}

	redisOptions := &redisclient.Options{
		Addr:         options.Addr,
		Username:     options.Username,
		Password:     options.Password,
		DB:           options.DB,
		PoolSize:     20, // Max number of socket connections
		MinIdleConns: 5,  // Minimum idle connections to keep open
	}
	if options.TLS {
		redisOptions.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	client := redisclient.NewClient(redisOptions)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("connect to Redis at %q: %w", options.Addr, err)
	}

	return &Client{client: client}, nil
}

func (c *Client) StoreURL(ctx context.Context, shortCode, longURL string) error {
	if strings.TrimSpace(shortCode) == "" {
		return errors.New("short code is required")
	}
	if strings.TrimSpace(longURL) == "" {
		return errors.New("long URL is required")
	}

	if err := c.client.Set(ctx, urlKey(shortCode), longURL, URLTTL).Err(); err != nil {
		return fmt.Errorf("store URL for short code %q: %w", shortCode, err)
	}
	return nil
}

func (c *Client) GetURL(ctx context.Context, shortCode string) (string, error) {
	if strings.TrimSpace(shortCode) == "" {
		return "", errors.New("short code is required")
	}

	longURL, err := c.client.Get(ctx, urlKey(shortCode)).Result()
	if errors.Is(err, redisclient.Nil) {
		return "", ErrURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("get URL for short code %q: %w", shortCode, err)
	}
	return longURL, nil
}

func (c *Client) DeleteURL(ctx context.Context, shortCode string) error {
	if strings.TrimSpace(shortCode) == "" {
		return errors.New("short code is required")
	}

	if err := c.client.Del(ctx, urlKey(shortCode)).Err(); err != nil {
		return fmt.Errorf("delete URL for short code %q: %w", shortCode, err)
	}
	return nil
}

func (c *Client) Close() error {
	if err := c.client.Close(); err != nil {
		return fmt.Errorf("close Redis client: %w", err)
	}
	return nil
}

func urlKey(shortCode string) string {
	return "url:" + shortCode
}
