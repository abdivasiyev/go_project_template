package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/abdivasiyev/project_template/config"
	"github.com/abdivasiyev/project_template/pkg/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"time"
)

var Module = fx.Provide(New)

type pgQuerier struct {
	db *sqlx.DB
}

type Params struct {
	fx.In
	Config config.Config
}

func New(params Params) (storage.Querier, error) {
	var (
		dbConn *sqlx.DB
		err    error
	)

	postgresURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		params.Config.GetString(config.PostgresUserKey),
		params.Config.GetString(config.PostgresPasswordKey),
		params.Config.GetString(config.PostgresHostKey),
		params.Config.GetInt(config.PostgresPortKey),
		params.Config.GetString(config.PostgresDatabaseKey),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbConn, err = sqlx.ConnectContext(ctx, "postgres", postgresURL)
	return &pgQuerier{db: dbConn}, err
}

func (q *pgQuerier) Query(ctx context.Context, query string, args ...any) (*sqlx.Rows, error) {
	return q.db.QueryxContext(ctx, query, args...)
}

func (q *pgQuerier) QueryRow(ctx context.Context, query string, args ...any) *sqlx.Row {
	return q.db.QueryRowxContext(ctx, query, args...)
}

func (q *pgQuerier) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return q.db.ExecContext(ctx, query, args...)
}

func (q *pgQuerier) Begin(ctx context.Context, options *sql.TxOptions) (*sqlx.Tx, error) {
	return q.db.BeginTxx(ctx, options)
}

func (q *pgQuerier) Prepare(ctx context.Context, query string) (*sqlx.Stmt, error) {
	return q.db.PreparexContext(ctx, query)
}

func (q *pgQuerier) PrepareNamed(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	return q.db.PrepareNamedContext(ctx, query)
}

func (q *pgQuerier) Close() error {
	return q.db.Close()
}
