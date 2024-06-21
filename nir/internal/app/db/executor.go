package db

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type executor struct {
	pool *pgxpool.Pool
}

func newExecutor(pool *pgxpool.Pool) *executor {
	return &executor{pool: pool}
}

func (e *executor) WithTx(ctx context.Context, level pgx.TxIsoLevel, action func(context.Context) error) error {
	tx, err := e.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   level,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return err
	}

	if err := action(ctx); err != nil {
		tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return err
	}

	return nil
}
