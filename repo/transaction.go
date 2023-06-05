package repo

import (
	"context"
	"database/sql"
	"errors"
)

var txKey = "tx_key"

// Basic db wrapper to handle transaction with the context
type dbTransactionAware struct {
	db DataBase
}

func (d dbTransactionAware) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if tx := getTx(ctx); tx != nil {
		//return tx.ExecContext(ctx, query, args)
		return tx.ExecContext(ctx, query, args...)
	}
	return d.db.ExecContext(ctx, query, args...)
}

func (d dbTransactionAware) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	if tx := getTx(ctx); tx != nil {
		return tx.PrepareContext(ctx, query)
	}
	return d.db.PrepareContext(ctx, query)
}

func (d dbTransactionAware) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if tx := getTx(ctx); tx != nil {
		return tx.QueryRowContext(ctx, query, args)
	}
	return d.db.QueryRowContext(ctx, query, args)
}

func (d dbTransactionAware) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if tx := getTx(ctx); tx != nil {
		return tx.QueryContext(ctx, query, args)
	}
	return d.db.QueryContext(ctx, query, args)
}

func (d dbTransactionAware) Begin(ctx context.Context) (context.Context, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, txKey, tx), nil
}

func (d dbTransactionAware) Commit(ctx context.Context) error {
	tx := getTx(ctx)
	if tx == nil {
		return errors.New("transaction not found, failed to commit")
	}

	return tx.Commit()
}

func (d dbTransactionAware) Rollback(ctx context.Context) error {
	tx := getTx(ctx)
	if tx == nil {
		return errors.New("transaction not found, failed to rollback")
	}

	return tx.Rollback()
}

func getTx(ctx context.Context) *sql.Tx {
	txValue := ctx.Value(txKey)
	if txValue == nil {
		return nil
	}
	tx, _ := txValue.(*sql.Tx) //if the assertion fails it will set tx to nil
	return tx
}
