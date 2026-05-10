package mongo

import (
	"context"
	"log/slog"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	mongoPort           = "27017"
	mongoStartupTimeout = 1 * time.Minute

	mongoEnvUsernameKey = "MONGO_INITDB_ROOT_USERNAME"
	mongoEnvPasswordKey = "MONGO_INITDB_ROOT_PASSWORD" //nolint:gosec
)

type Container struct {
	container testcontainers.Container
	client    *mongo.Client
	cfg       *Config
}

func NewContainer(ctx context.Context, opts ...Option) (*Container, error) {
	cfg := buildConfig(opts...)

	container, err := startMongoContainer(ctx, cfg)
	if err != nil {
		return nil, err
	}

	success := false
	defer func() {
		if !success {
			if err = container.Terminate(ctx); err != nil {
				cfg.Logger.Error("failed to terminate mongo container", slog.String("error", err.Error()))
			}
		}
	}()

	cfg.Host, cfg.Port, err = getContainerHostPort(ctx, container)
	if err != nil {
		return nil, err
	}

	uri := buildMongoURI(cfg)

	client, err := connectMongoClient(ctx, uri)
	if err != nil {
		return nil, err
	}

	cfg.Logger.Info("Mongo container started", slog.String("host", cfg.Host), slog.String("port", cfg.Port))
	success = true

	return &Container{
		container: container,
		client:    client,
		cfg:       cfg,
	}, nil
}

func (c *Container) Client() *mongo.Client {
	return c.client
}

func (c *Container) Config() *Config {
	return c.cfg
}

func (c *Container) Terminate(ctx context.Context) error {
	if err := c.client.Disconnect(ctx); err != nil {
		c.cfg.Logger.Error("failed to disconnect mongo client", slog.String("error", err.Error()))
	}

	if err := c.container.Terminate(ctx); err != nil {
		c.cfg.Logger.Error("failed to terminate mongo container", slog.String("error", err.Error()))
	}

	c.cfg.Logger.Info("Mongo container terminated")

	return nil
}
