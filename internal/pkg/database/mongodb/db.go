package mongodb

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var ErrMappingNotFound = errors.New("URL mapping not found")
var ErrMappingConflict = errors.New("URL mapping conflicts with an existing record")

type URLMapping struct {
	ShortCode string `bson:"short_code"`
	LongURL   string `bson:"long_url"`
}

type Client struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func Connect(ctx context.Context, uri string) (*Client, error) {
	if strings.TrimSpace(uri) == "" {
		return nil, errors.New("MongoDB connection string is required")
	}

	databaseName := os.Getenv("MONGODB_DATABASE_NAME")
	if strings.TrimSpace(databaseName) == "" {
		return nil, errors.New("MongoDB database name is required")
	}

	collectionName := os.Getenv("MONGODB_COLLECTION_NAME")
	if strings.TrimSpace(collectionName) == "" {
		return nil, errors.New("MongoDB collection name is required")
	}

	client, err := mongo.Connect(
		options.Client().
			ApplyURI(uri).
			SetMaxPoolSize(10).
			SetMinPoolSize(1), // setting a minimum pool size to ensure some connections are always available
	)
	if err != nil {
		return nil, fmt.Errorf("create MongoDB client: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("ping MongoDB: %w", err)
	}

	collection := client.Database(databaseName).Collection(collectionName)
	_, err = collection.Indexes().CreateMany(pingCtx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "short_code", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "long_url", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("create URL mapping indexes: %w", err)
	}

	return &Client{client: client, collection: collection}, nil
}

func (c *Client) QueryLongURL(ctx context.Context, longURL string) (string, error) {
	var mapping URLMapping
	err := c.collection.FindOne(ctx, bson.D{{Key: "long_url", Value: longURL}}).Decode(&mapping)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return "", ErrMappingNotFound
	}
	if err != nil {
		return "", fmt.Errorf("find mapping by long URL: %w", err)
	}
	return mapping.ShortCode, nil
}

func (c *Client) GetRedirectUrl(ctx context.Context, shortCode string) (string, error) {
	var mapping URLMapping
	err := c.collection.FindOne(ctx, bson.D{{Key: "short_code", Value: shortCode}}).Decode(&mapping)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return "", ErrMappingNotFound
	}
	if err != nil {
		return "", fmt.Errorf("find mapping by short code: %w", err)
	}
	return mapping.LongURL, nil
}

func (c *Client) CreateEntry(ctx context.Context, shortCode, longURL string) error {
	_, err := c.collection.InsertOne(ctx, URLMapping{
		ShortCode: shortCode,
		LongURL:   longURL,
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("%w: %v", ErrMappingConflict, err)
		}
		return fmt.Errorf("insert URL mapping: %w", err)
	}
	return nil
}

func (c *Client) Disconnect(ctx context.Context) error {
	if err := c.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("disconnect MongoDB: %w", err)
	}
	return nil
}
