package repo

import (
	"context"
	"database/sql"

	"github.com/tasker/service"
)

type DataBase interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type Repository interface {
	SaveTask(ctx context.Context, task service.Task) (service.Task, error)
	SaveExecution(ctx context.Context, execution service.Execution) error
}

type repository struct {
	db DataBase
}

func (r repository) SaveExecution(ctx context.Context, execution service.Execution) error {
	//TODO implement me
	panic("implement me")
}

func NewRepository(db DataBase) Repository {
	return &repository{db: db}
}
