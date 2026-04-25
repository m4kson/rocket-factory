package env

import "github.com/caarlos0/env/v11"

type MongoEnvConfig struct {
	User       string `env:"MONGO_INITDB_ROOT_USERNAME,required"`
	Password   string `env:"MONGO_INITDB_ROOT_PASSWORD,required"`
	DbName     string `env:"MONGO_INITDB_DATABASE,required"`
	Port       string `env:"MONGO_PORT,required"`
	AuthDbName string `env:"MONGO_AUTH_DB,required"`
	Url        string `env:"MONGO_URI,required"`
}

type MongoConfig struct {
	raw MongoEnvConfig
}

func NewMongoConfig() (*MongoConfig, error) {
	var raw MongoEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &MongoConfig{raw: raw}, nil
}

func (cfg *MongoConfig) User() string {
	return cfg.raw.User
}

func (cfg *MongoConfig) Password() string {
	return cfg.raw.Password
}

func (cfg *MongoConfig) DbName() string {
	return cfg.raw.DbName
}

func (cfg *MongoConfig) Port() string {
	return cfg.raw.Port
}

func (cfg *MongoConfig) AuthDbName() string {
	return cfg.raw.AuthDbName
}

func (cfg *MongoConfig) URL() string {
	return cfg.raw.Url
}
