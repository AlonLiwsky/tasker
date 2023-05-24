package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/tasker/service"
)

const (
	InsertTaskQr = "INSERT INTO task (Name) VALUES (?);"
)

func (r repository) CreateTask(ctx context.Context, task service.Task) (int, error) {
	result, err := r.db.ExecContext(ctx, InsertTaskQr, task.Name)
	if err != nil {
		return 0, fmt.Errorf("inserting event: %w", err)
	}
	rAffect, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rAffect != 1 {
		return 0, errors.New(fmt.Sprintf("rows affected while inserting task different than 1: %v", rAffect))
	}

	id, err := result.LastInsertId()
	return int(id), err
}
