package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

func newTestClient(t *testing.T) (*Client, *miniredis.Miniredis) {
	t.Helper()

	server := miniredis.RunT(t)
	client, err := NewClient(context.Background(), Options{Addr: server.Addr()})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Errorf("Close() error = %v", err)
		}
	})

	return client, server
}

func TestStoreGetAndDeleteURL(t *testing.T) {
	client, server := newTestClient(t)
	ctx := context.Background()

	if err := client.StoreURL(ctx, "abc123", "https://example.com"); err != nil {
		t.Fatalf("StoreURL() error = %v", err)
	}
	if got := server.TTL("url:abc123"); got != URLTTL {
		t.Errorf("stored TTL = %s, want %s", got, URLTTL)
	}

	got, err := client.GetURL(ctx, "abc123")
	if err != nil {
		t.Fatalf("GetURL() error = %v", err)
	}
	if got != "https://example.com" {
		t.Errorf("GetURL() = %q, want %q", got, "https://example.com")
	}

	if err := client.DeleteURL(ctx, "abc123"); err != nil {
		t.Fatalf("DeleteURL() error = %v", err)
	}
	if _, err := client.GetURL(ctx, "abc123"); !errors.Is(err, ErrURLNotFound) {
		t.Errorf("GetURL() error = %v, want ErrURLNotFound", err)
	}
}

func TestGetURLReturnsNotFound(t *testing.T) {
	client, _ := newTestClient(t)

	_, err := client.GetURL(context.Background(), "missing")
	if !errors.Is(err, ErrURLNotFound) {
		t.Fatalf("GetURL() error = %v, want ErrURLNotFound", err)
	}
}

func TestNewClientRejectsEmptyAddress(t *testing.T) {
	_, err := NewClient(context.Background(), Options{})
	if err == nil {
		t.Fatal("NewClient() error = nil, want an error")
	}
}

func TestNewClientUsesCredentialsAndDatabase(t *testing.T) {
	server := miniredis.RunT(t)
	server.RequireAuth("test-password")

	client, err := NewClient(context.Background(), Options{
		Addr:     server.Addr(),
		Username: "default",
		Password: "test-password",
		DB:       2,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Errorf("Close() error = %v", err)
		}
	})

	if err := client.StoreURL(context.Background(), "db-check", "https://example.com"); err != nil {
		t.Fatalf("StoreURL() error = %v", err)
	}
	got, err := server.DB(2).Get("url:db-check")
	if err != nil {
		t.Fatalf("read stored value from Redis DB 2: %v", err)
	}
	if got != "https://example.com" {
		t.Errorf("stored value in Redis DB 2 = %q, want %q", got, "https://example.com")
	}
}

func TestURLTTLIsThirtyDays(t *testing.T) {
	if URLTTL != 30*24*time.Hour {
		t.Fatalf("URLTTL = %s, want 30 days", URLTTL)
	}
}
