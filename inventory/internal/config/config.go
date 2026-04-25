package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/m4kson/rocket-factory/inventory/internal/config/env"
)

var appConfig *config

type config struct {
	Logger LoggerConfig
	Grpc   GrpcConfig
	Mongo  MongoConfig
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

	mongoCfg, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger: loggerCfg,
		Grpc:   grpcCfg,
		Mongo:  mongoCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
