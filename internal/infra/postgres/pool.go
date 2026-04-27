package core_infra_postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Close()

	GetOperationTimeout() time.Duration
}

type ConnectionPool struct {
	*pgxpool.Pool
	operationTimeout time.Duration
}

func NewConnectionPool(
	ctx context.Context,
	cfg Config,
) (*ConnectionPool, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Db,
	)

	pgxconfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("pgx config parse: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxconfig)
	if err != nil {
		return nil, fmt.Errorf("pgx pool init: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pgx ping: %v", err)
	}

	return &ConnectionPool{
		Pool:             pool,
		operationTimeout: cfg.Timeout,
	}, nil
}

func (p *ConnectionPool) GetOperationTimeout() time.Duration {
	return p.operationTimeout
}
