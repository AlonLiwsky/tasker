package mgmtDB

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/tasker/entities"
	"github.com/tasker/http"
)

const (
	InsertExecQr         = "INSERT INTO execution (scheduled_task_id, task_id, status, idempotency_token, executed_time) VALUES (?, ?, ?, ?, ?);"
	GetExecIdempotencyQr = "SELECT * FROM execution WHERE idempotency_token = ?"
)

func (r repository) SaveExecution(ctx context.Context, exec entities.Execution) (entities.Execution, error) {
	result, err := r.db.ExecContext(ctx, InsertExecQr, exec.ScheduledTask, exec.TaskID, exec.Status, exec.IdempotencyToken, exec.ExecutedTime)
	if err != nil {
		return entities.Execution{}, fmt.Errorf("inserting execution: %w", err)
	}

	rAffect, err := result.RowsAffected()
	switch {
	case err != nil:
		return entities.Execution{}, err
	case rAffect != 1:
		return entities.Execution{}, fmt.Errorf("inserting execution: should affect 1 and affected #%d rows", rAffect)
	}

	execID, err := result.LastInsertId()
	if err != nil {
		return entities.Execution{}, err
	}

	exec.ID = int(execID)
	return exec, nil
}

func (r repository) GetExecutionIdempotency(ctx context.Context, idempToken string) (entities.Execution, error) {
	var aux *string
	exec := entities.Execution{}
	execTimeString := ""
	row := r.db.QueryRowContext(ctx, GetExecIdempotencyQr, idempToken)
	err := row.Scan(&exec.ID, &exec.ScheduledTask, &exec.TaskID, &aux, &exec.Status, &exec.IdempotencyToken, &aux, &execTimeString, &aux)
	switch {
	case err == sql.ErrNoRows:
		return entities.Execution{}, http.WrapError(err, http.ErrNotFound.WithMessage("execution not found"))
	case err != nil:
		return entities.Execution{}, fmt.Errorf("getting task: %w", err)
	}

	parsedTime, err := time.Parse(time.DateTime, execTimeString)
	if err != nil {
		log.Printf("Error unmarshalling JSON: %s. unmarshalling executed_time", err)
	}
	exec.ExecutedTime = parsedTime

	return exec, nil
}
