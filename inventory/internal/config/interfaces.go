package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type GrpcConfig interface {
	Host() string
	Port() string
}

type MongoConfig interface {
	User() string
	Password() string
	DbName() string
	Port() string
	AuthDbName() string
	URL() string
}
