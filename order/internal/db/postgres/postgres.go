package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string

	MaxConns          int32
	MinConns          int32
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCehckPeriod time.Duration
}

func NewPool(ctx context.Context, config Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres: failed to parse config: %w", err)
	}

	poolCfg.MaxConns = config.MaxConns
	poolCfg.MinConns = config.MinConns
	poolCfg.MaxConnLifetime = config.MaxConnLifetime
	poolCfg.MaxConnIdleTime = config.MaxConnIdleTime
	poolCfg.HealthCheckPeriod = config.HealthCehckPeriod

	poolCfg.PrepareConn = func(ctx context.Context, conn *pgx.Conn) (bool, error) {
		if err := conn.Ping(ctx); err != nil {
			return false, nil
		}
		return true, nil
	}

	poolCfg.ConnConfig.RuntimeParams["statement_timeout"] = "5000"

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("postgres: failed to connect to pool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres: failed to ping pool: %w", err)
	}

	return pool, nil
}
