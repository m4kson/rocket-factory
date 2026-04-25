package config

import "time"

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
