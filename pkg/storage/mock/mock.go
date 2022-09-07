package mock

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/abdivasiyev/project_template/config"
	"github.com/abdivasiyev/project_template/pkg/storage"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
)

var Module = fx.Provide(New)

type mockQuerier struct {
	db   *sqlx.DB
	mock sqlmock.Sqlmock
}

type Params struct {
	fx.In
	Config config.Config
}

func New(context.Context, Params) (storage.Querier, error) {
	db, mock, err := sqlmock.New()

	dbConn := sqlx.NewDb(db, "mock")

	return &mockQuerier{db: dbConn, mock: mock}, err
}

func (q *mockQuerier) GetMock() sqlmock.Sqlmock {
	return q.mock
}

func (q *mockQuerier) Query(ctx context.Context, query string, args ...any) (*sqlx.Rows, error) {
	return q.db.QueryxContext(ctx, query, args...)
}

func (q *mockQuerier) QueryRow(ctx context.Context, query string, args ...any) *sqlx.Row {
	return q.db.QueryRowxContext(ctx, query, args...)
}

func (q *mockQuerier) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return q.db.ExecContext(ctx, query, args...)
}

func (q *mockQuerier) Begin(ctx context.Context, options *sql.TxOptions) (*sqlx.Tx, error) {
	return q.db.BeginTxx(ctx, options)
}

func (q *mockQuerier) Prepare(ctx context.Context, query string) (*sqlx.Stmt, error) {
	return q.db.PreparexContext(ctx, query)
}

func (q *mockQuerier) PrepareNamed(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	return q.db.PrepareNamedContext(ctx, query)
}

func (q *mockQuerier) Close() error {
	return q.db.Close()
}
