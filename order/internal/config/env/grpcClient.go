package env

import "github.com/caarlos0/env/v11"

type grpcClientEnvConfig struct {
	InventoryGrpcAddr string `env:"INVENTORY_GRPC_ADDR,required"`
	PaymentGrpcAddr   string `env:"PAYMENT_GRPC_ADDR,required"`
}

type grpcClientConfig struct {
	raw grpcClientEnvConfig
}

func NewGRPCClientConfig() (*grpcClientConfig, error) {
	var raw grpcClientEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &grpcClientConfig{raw: raw}, nil
}

func (cfg *grpcClientConfig) InventoryGrpcAddr() string {
	return cfg.raw.InventoryGrpcAddr
}

func (cfg *grpcClientConfig) PaymentGrpcAddr() string {
	return cfg.raw.PaymentGrpcAddr
}
