package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/m4kson/rocket-factory/payment/internal/config/env"
)

var appConfig *config

type config struct {
	Logger LoggerConfig
	Grpc   GrpcConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	grpcCfg, err := env.NewGrpcConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger: loggerCfg,
		Grpc:   grpcCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
