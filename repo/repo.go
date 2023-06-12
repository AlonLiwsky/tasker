package repo

import (
	"context"
	"database/sql"

	"github.com/tasker/entities"
)

type DataBase interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Repository interface {
	SaveTask(ctx context.Context, task entities.Task) (entities.Task, error)
	GetTask(ctx context.Context, taskID int) (entities.Task, error)
	SaveExecution(ctx context.Context, exec entities.Execution) (entities.Execution, error)
}

type repository struct {
	db dbTransactionAware
}

func NewRepository(db DataBase) Repository {
	return &repository{
		db: dbTransactionAware{
			db: db,
		},
	}
}
