package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/m4kson/rocket-factory/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger     LoggerConfig
	GrpcClient GrpcClientsConfig
	Postgres   PostgresConfig
	HttpServer HttpServerConfig

	Kafka                  KafkaConfig
	OrderPaidProducer      OrderPaidProducerConfig
	OrderAssembledConsumer OrderAssembledConsumerConfig
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

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	grpcClientCfg, err := env.NewGRPCClientConfig()
	if err != nil {
		return err
	}

	httpServerCfg, err := env.NewHttpServerConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderPaidProducerCfg, err := env.NewOrderPaidProducerConfig()
	if err != nil {
		return err
	}

	orderAssembledConsumerCfg, err := env.NewOrderAssembledConsumerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                 loggerCfg,
		GrpcClient:             grpcClientCfg,
		Postgres:               postgresCfg,
		HttpServer:             httpServerCfg,
		Kafka:                  kafkaCfg,
		OrderPaidProducer:      orderPaidProducerCfg,
		OrderAssembledConsumer: orderAssembledConsumerCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
