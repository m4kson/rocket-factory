package config

import (
	"time"

	"github.com/IBM/sarama"
)

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type GrpcClientsConfig interface {
	InventoryGrpcAddr() string
	PaymentGrpcAddr() string
}

type PostgresConfig interface {
	User() string
	Password() string
	DbName() string
	Host() string
	Port() string
	URL() string
	MigrationsPath() string
}

type HttpServerConfig interface {
	Port() string
	ReadHeaderTimeout() time.Duration
	ShutdownTimeout() time.Duration
}

type OrderAssembledConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

type OrderPaidProducerConfig interface {
	TopicName() string
	Config() *sarama.Config
}

type KafkaConfig interface {
	Brokers() []string
}
