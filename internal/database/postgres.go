package database

import (
	"context"
	"errors"

	"github.com/itskyana/belajar-grpc-go/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(cfg *config.Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.GetDSN())
	if err != nil {
		return nil, errors.New("unable to parse configuration: " + err.Error())
	}

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, errors.New("unable to create connection pool: " + err.Error())
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, errors.New("unable to ping database: " + err.Error())
	}

	return pool, nil
}
