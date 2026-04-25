package env

import "github.com/caarlos0/env/v11"

type GrpcEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type GrpcConfig struct {
	raw GrpcEnvConfig
}

func NewGrpcConfig() (*GrpcConfig, error) {
	var raw GrpcEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &GrpcConfig{raw: raw}, nil
}

func (cfg *GrpcConfig) Host() string {
	return cfg.raw.Host
}

func (cfg *GrpcConfig) Port() string {
	return cfg.raw.Port
}
