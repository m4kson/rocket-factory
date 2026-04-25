package env

import "github.com/caarlos0/env/v11"

type postgresEnvConfig struct {
	User           string `env:"POSTGRES_USER,required"`
	Password       string `env:"POSTGRES_PASSWORD,required"`
	DbName         string `env:"POSTGRES_DB,required"`
	Host           string `env:"POSTGRES_HOST,required"`
	Port           string `env:"POSTGRES_PORT,required"`
	URL            string `env:"DATABASE_URL,required"`
	MigrationsPath string `env:"MIGRATIONS_PATH,required"`
}

type postgresConfig struct {
	raw postgresEnvConfig
}

func NewPostgresConfig() (*postgresConfig, error) {
	var raw postgresEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &postgresConfig{raw: raw}, nil
}

func (cfg *postgresConfig) User() string {
	return cfg.raw.User
}

func (cfg *postgresConfig) Password() string {
	return cfg.raw.Password
}

func (cfg *postgresConfig) DbName() string {
	return cfg.raw.DbName
}

func (cfg *postgresConfig) Host() string {
	return cfg.raw.Host
}

func (cfg *postgresConfig) Port() string {
	return cfg.raw.Port
}

func (cfg *postgresConfig) URL() string {
	return cfg.raw.URL
}

func (cfg *postgresConfig) MigrationsPath() string {
	return cfg.raw.MigrationsPath
}
