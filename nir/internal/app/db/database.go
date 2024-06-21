package db

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	cluster *pgxpool.Pool
}

type DBops interface {
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetPool(_ context.Context) *pgxpool.Pool
}

func newDatabase(cluster *pgxpool.Pool) *Database {
	return &Database{cluster: cluster}
}

func (db Database) GetPool(_ context.Context) *pgxpool.Pool {
	return db.cluster
}

func (db Database) GetExecutor() *executor {
	return newExecutor(db.cluster)
}

func (db Database) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	exec := db.GetExecutor()

	return exec.WithTx(ctx, pgx.Serializable, func(context context.Context) error {
		return pgxscan.Get(context, exec.pool, dest, query, args...)
	})
}

func (db Database) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	exec := db.GetExecutor()

	return exec.WithTx(ctx, pgx.Serializable, func(context context.Context) error {
		return pgxscan.Select(context, exec.pool, dest, query, args...)
	})
}

func (db Database) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	exec := db.GetExecutor()
	var res pgconn.CommandTag

	err := exec.WithTx(ctx, pgx.Serializable, func(context context.Context) error {
		var err error
		res, err = db.cluster.Exec(ctx, query, args...)
		return err
	})

	return res, err
}

func (db Database) ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	exec := db.GetExecutor()
	var res pgx.Row

	exec.WithTx(ctx, pgx.Serializable, func(context context.Context) error {
		res = db.cluster.QueryRow(ctx, query, args...)
		return nil

	})

	return res
}
