package postgres

import (
	"context"
	"fmt"
	"net"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mokuhyo-driven-test/api/internal/repository"
)

// PostgresDB はローカルPostgreSQLデータベース接続の実装です
type PostgresDB struct {
	pool *pgxpool.Pool
}

// NewPostgresDB は新しいPostgreSQLデータベース接続を作成します
func NewPostgresDB(connString string) (repository.DBInterface, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// IPv4接続を強制（IPv6接続の問題を回避）
	config.ConnConfig.DialFunc = func(ctx context.Context, network, addr string) (net.Conn, error) {
		d := &net.Dialer{}
		return d.DialContext(ctx, "tcp4", addr)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresDB{pool: pool}, nil
}

// QueryRow は単一の行を取得するクエリを実行します
func (db *PostgresDB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return db.pool.QueryRow(ctx, sql, args...)
}

// Query は複数の行を取得するクエリを実行します
func (db *PostgresDB) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return db.pool.Query(ctx, sql, args...)
}

// Exec はINSERT、UPDATE、DELETEなどのクエリを実行します
func (db *PostgresDB) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return db.pool.Exec(ctx, sql, args...)
}

// Begin はトランザクションを開始します
func (db *PostgresDB) Begin(ctx context.Context) (repository.TxInterface, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &PostgresTx{tx: tx}, nil
}

// Close はデータベース接続を閉じます
func (db *PostgresDB) Close() {
	db.pool.Close()
}

// PostgresTx はPostgreSQLトランザクションの実装です
type PostgresTx struct {
	tx pgx.Tx
}

// QueryRow はトランザクション内で単一の行を取得するクエリを実行します
func (tx *PostgresTx) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return tx.tx.QueryRow(ctx, sql, args...)
}

// Query はトランザクション内で複数の行を取得するクエリを実行します
func (tx *PostgresTx) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return tx.tx.Query(ctx, sql, args...)
}

// Exec はトランザクション内でINSERT、UPDATE、DELETEなどのクエリを実行します
func (tx *PostgresTx) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return tx.tx.Exec(ctx, sql, args...)
}

// Commit はトランザクションをコミットします
func (tx *PostgresTx) Commit(ctx context.Context) error {
	return tx.tx.Commit(ctx)
}

// Rollback はトランザクションをロールバックします
func (tx *PostgresTx) Rollback(ctx context.Context) error {
	return tx.tx.Rollback(ctx)
}
