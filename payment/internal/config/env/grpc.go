package env

import "github.com/caarlos0/env/v11"

type grpcEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type grpcConfig struct {
	raw grpcEnvConfig
}

func NewGrpcConfig() (*grpcConfig, error) {
	var raw grpcEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &grpcConfig{raw: raw}, nil
}

func (cfg *grpcConfig) Host() string {
	return cfg.raw.Host
}

func (cfg *grpcConfig) Port() string {
	return cfg.raw.Port
}
