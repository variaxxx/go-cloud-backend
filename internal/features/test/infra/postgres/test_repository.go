package test_postgres

import (
	infra_postgres "cloud/internal/infra/postgres"
	"context"
	"fmt"
)

type testRepository struct {
	pool infra_postgres.Pool
}

func NewTestRepository(
	pool infra_postgres.Pool,
) *testRepository {
	return &testRepository{
		pool: pool,
	}
}

func (r *testRepository) WriteStr(
	ctx context.Context,
	str string,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetOperationTimeout())
	defer cancel()

	const query = `
		INSERT INTO cloud.test_strings (value)
		VALUES ($1);
	`

	if _, err := r.pool.Exec(ctx, query, str); err != nil {
		return fmt.Errorf("write test string: %w", err)
	}

	return nil
}
