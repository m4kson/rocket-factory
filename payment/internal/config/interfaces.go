package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type GrpcConfig interface {
	Host() string
	Port() string
}
