package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type begin interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type TxManager struct {
	db begin
}

func NewTxManager(db begin) *TxManager {
	return &TxManager{db: db}
}

type txContextKey struct{}

// WithinTx executes the provided function (fn) within a database transaction.
// If a transaction is already present in the context, it will use that transaction instead of starting a new one.
//   - fn MUST pass the ctx it receives to every store call, because that ctx carries the transaction.
//
// If the function (fn) returns an error, the transaction will be rolled back; otherwise, it will be committed.
func (m *TxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := ctx.Value(txContextKey{}).(pgx.Tx); ok {
		return fn(ctx)
	}

	tx, err := m.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		rbCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		_ = tx.Rollback(rbCtx)
	}()

	if err := fn(context.WithValue(ctx, txContextKey{}, tx)); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

// Executor returns the transaction (pgx.Tx) carried by ctx if it exists, or fallback to the pool for executing standalone queries.
func Executor(ctx context.Context, pool DBTX) DBTX {
	if tx, ok := ctx.Value(txContextKey{}).(pgx.Tx); ok {
		return tx
	}
	return pool
}
