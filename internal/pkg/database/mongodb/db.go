package mongodb

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Client struct {
	client *mongo.Client
}

func Connect(ctx context.Context, uri string) (*Client, error) {
	if strings.TrimSpace(uri) == "" {
		return nil, errors.New("MongoDB connection string is required")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("create MongoDB client: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("ping MongoDB: %w", err)
	}

	return &Client{client: client}, nil
}

func (c *Client) Database(name string) *mongo.Database {
	return c.client.Database(name)
}

func (c *Client) Disconnect(ctx context.Context) error {
	if err := c.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("disconnect MongoDB: %w", err)
	}
	return nil
}
