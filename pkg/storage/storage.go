package storage

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"io"
	"time"
)

type Querier interface {
	io.Closer
	Query(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) *sqlx.Row
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
	Begin(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	Prepare(ctx context.Context, query string) (*sqlx.Stmt, error)
	PrepareNamed(ctx context.Context, query string) (*sqlx.NamedStmt, error)
}

type Cacher interface {
	io.Closer
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, duration time.Duration) error
	GetObj(ctx context.Context, key string, value any) error
	SetObj(ctx context.Context, key string, value any, duration time.Duration) error
	Delete(ctx context.Context, keys ...string) error
}
