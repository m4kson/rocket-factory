package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Config struct {
	URI             string
	Database        string
	ConnectTimeout  time.Duration
	MaxPoolSize     uint64
	MinPoolSize     uint64
	MaxConnIdleTime time.Duration
}

type Client struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	opts := options.Client().
		ApplyURI(cfg.URI).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMinPoolSize(cfg.MinPoolSize).
		SetMaxConnIdleTime(cfg.MaxConnIdleTime).
		SetConnectTimeout(cfg.ConnectTimeout).
		SetServerSelectionTimeout(5 * time.Second).
		SetCompressors([]string{"zstd", "snappy"})

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("mongo: connect failed: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	if err = client.Ping(pingCtx, readpref.Primary()); err != nil {
		_ = client.Disconnect(ctx)
		return nil, fmt.Errorf("mongo: ping failed: %w", err)
	}

	return &Client{
		client: client,
		db:     client.Database(cfg.Database),
	}, nil
}

func (c *Client) Collection(name string) *mongo.Collection {
	return c.db.Collection(name)
}

func (c *Client) Disconnect(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

func (c *Client) Client() *mongo.Client {
	return c.client
}
