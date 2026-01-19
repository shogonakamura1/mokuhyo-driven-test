package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DBInterface はデータベース接続のインターフェースです
// クリーンアーキテクチャに基づき、実装詳細を抽象化します
type DBInterface interface {
	// QueryRow は単一の行を取得するクエリを実行します
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row

	// Query は複数の行を取得するクエリを実行します
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)

	// Exec はINSERT、UPDATE、DELETEなどのクエリを実行します
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)

	// Begin はトランザクションを開始します
	Begin(ctx context.Context) (TxInterface, error)

	// Close はデータベース接続を閉じます
	Close()
}

// TxInterface はトランザクションのインターフェースです
type TxInterface interface {
	// QueryRow はトランザクション内で単一の行を取得するクエリを実行します
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row

	// Query はトランザクション内で複数の行を取得するクエリを実行します
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)

	// Exec はトランザクション内でINSERT、UPDATE、DELETEなどのクエリを実行します
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)

	// Commit はトランザクションをコミットします
	Commit(ctx context.Context) error

	// Rollback はトランザクションをロールバックします
	Rollback(ctx context.Context) error
}
