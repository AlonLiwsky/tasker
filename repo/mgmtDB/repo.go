package mgmtDB

import (
	"context"
	"database/sql"
	"time"

	"github.com/tasker/entities"
)

type DataBase interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type DataBaseTransactionAware interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Repository interface {
	SaveTask(ctx context.Context, task entities.Task) (entities.Task, error)
	GetTask(ctx context.Context, taskID int) (entities.Task, error)
	SaveExecution(ctx context.Context, exec entities.Execution) (entities.Execution, error)
	GetExecutionIdempotency(ctx context.Context, idempToken string) (entities.Execution, error)
	SaveSchedule(ctx context.Context, sch entities.ScheduledTask) (entities.ScheduledTask, error)
	GetEnabledSchedules(ctx context.Context) ([]entities.ScheduledTask, error)
	SetScheduleLastRun(ctx context.Context, schID int, time time.Time) error
}

type repository struct {
	db DataBaseTransactionAware
}

func NewRepository(db DataBase) Repository {
	return &repository{
		db: dbTransactionAware{
			db: db,
		},
	}
}
