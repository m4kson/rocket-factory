package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type httpServerEnvConfig struct {
	Port              string        `env:"HTTP_PORT,required"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT,required"`
	ShutdownTimeout   time.Duration `env:"SHUTDOWN_TIMEOUT,required"`
}

type httpServerConfig struct {
	raw httpServerEnvConfig
}

func NewHttpServerConfig() (*httpServerConfig, error) {
	var raw httpServerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &httpServerConfig{raw: raw}, nil
}

func (cfg *httpServerConfig) Port() string {
	return cfg.raw.Port
}

func (cfg *httpServerConfig) ReadHeaderTimeout() time.Duration {
	return cfg.raw.ReadHeaderTimeout
}

func (cfg *httpServerConfig) ShutdownTimeout() time.Duration {
	return cfg.raw.ShutdownTimeout
}
