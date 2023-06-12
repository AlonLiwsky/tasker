package repo

import (
	"context"
	"fmt"

	"github.com/tasker/entities"
)

const (
	InsertExecQr = "INSERT INTO execution (scheduled_task_id, status, executed_time) VALUES (?, ?, ?);"
)

func (r repository) SaveExecution(ctx context.Context, exec entities.Execution) (entities.Execution, error) {
	result, err := r.db.ExecContext(ctx, InsertExecQr, exec.ScheduledTask, exec.Status, exec.ExecutedTime)
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
