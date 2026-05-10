package mongo

import (
	"log/slog"

	log "github.com/m4kson/rocket-factory/platform/pkg/logger/slogLog"
	"github.com/moby/moby/api/types/container"
)

type Config struct {
	NetworkName   string
	ContainerName string
	ImageName     string
	Database      string
	Username      string
	Password      string
	AuthDB        string

	Logger *slog.Logger

	Host string
	Port string
}

func buildConfig(opts ...Option) *Config {
	cfg := &Config{
		NetworkName:   "test-network",
		ContainerName: "mongo-inventory-test",
		ImageName:     "mongo:8.0",
		Database:      "test",
		Username:      "root",
		Password:      "root",
		AuthDB:        "admin",
		Logger:        log.NewNop(),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

func defaultHostConfig() func(hc *container.HostConfig) {
	return func(hc *container.HostConfig) {
		hc.AutoRemove = true
	}
}
