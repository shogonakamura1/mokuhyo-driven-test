package supabase

import (
	"context"
	"fmt"
	"net"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mokuhyo-driven-test/api/internal/repository"
)

// SupabaseDB はSupabaseデータベース接続の実装です
// デプロイ時に使用されます
type SupabaseDB struct {
	pool *pgxpool.Pool
}

// NewSupabaseDB は新しいSupabaseデータベース接続を作成します
func NewSupabaseDB(connString string) (repository.DBInterface, error) {
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

	return &SupabaseDB{pool: pool}, nil
}

// QueryRow は単一の行を取得するクエリを実行します
func (db *SupabaseDB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return db.pool.QueryRow(ctx, sql, args...)
}

// Query は複数の行を取得するクエリを実行します
func (db *SupabaseDB) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return db.pool.Query(ctx, sql, args...)
}

// Exec はINSERT、UPDATE、DELETEなどのクエリを実行します
func (db *SupabaseDB) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return db.pool.Exec(ctx, sql, args...)
}

// Begin はトランザクションを開始します
func (db *SupabaseDB) Begin(ctx context.Context) (repository.TxInterface, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &SupabaseTx{tx: tx}, nil
}

// Close はデータベース接続を閉じます
func (db *SupabaseDB) Close() {
	db.pool.Close()
}

// SupabaseTx はSupabaseトランザクションの実装です
type SupabaseTx struct {
	tx pgx.Tx
}

// QueryRow はトランザクション内で単一の行を取得するクエリを実行します
func (tx *SupabaseTx) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return tx.tx.QueryRow(ctx, sql, args...)
}

// Query はトランザクション内で複数の行を取得するクエリを実行します
func (tx *SupabaseTx) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return tx.tx.Query(ctx, sql, args...)
}

// Exec はトランザクション内でINSERT、UPDATE、DELETEなどのクエリを実行します
func (tx *SupabaseTx) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return tx.tx.Exec(ctx, sql, args...)
}

// Commit はトランザクションをコミットします
func (tx *SupabaseTx) Commit(ctx context.Context) error {
	return tx.tx.Commit(ctx)
}

// Rollback はトランザクションをロールバックします
func (tx *SupabaseTx) Rollback(ctx context.Context) error {
	return tx.tx.Rollback(ctx)
}
