package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/tasker/service"
)

const (
	InsertTaskQr = "INSERT INTO task (name) VALUES (?);"
	InsertStepQr = "INSERT INTO step (task_id, step_type, params, failure_step, position) VALUES (?, ?, ?, ?, ?)"
)

func (r repository) SaveTask(ctx context.Context, task service.Task) (service.Task, error) {
	result, err := r.db.ExecContext(ctx, InsertTaskQr, task.Name)
	if err != nil {
		return service.Task{}, fmt.Errorf("inserting event: %w", err)
	}
	rAffect, err := result.RowsAffected()
	if err != nil {
		return service.Task{}, err
	}
	if rAffect != 1 {
		return service.Task{}, errors.New(fmt.Sprintf("rows affected while inserting task different than 1: %v", rAffect))
	}

	auxID, err := result.LastInsertId()
	task.ID = int(auxID)

	steps, err := r.saveSteps(ctx, task.Steps)
	if err != nil {
		return service.Task{}, err
	}
	task.Steps = steps

	return task, err
}

func (r repository) saveSteps(ctx context.Context, steps []service.Step) ([]service.Step, error) {
	//Save each step with the Order field

	// Prepare the SQL statement
	stmt, err := r.db.PrepareContext(ctx, InsertStepQr)
	if err != nil {
		return []service.Step{}, fmt.Errorf("preparing insert steps stmt: %w", err)
	}
	defer stmt.Close()

	// Iterate over the rows and execute the prepared statement in a batch
	for position, step := range steps {
		var failureStepID *int = nil
		if step.FailureStep != nil {
			*failureStepID, err = insertFailureStep(ctx, *step.FailureStep)
			if err != nil {
				return []service.Step{}, err //TODO: proper error handling
			}
		}

		result, err := stmt.Exec(step.Task.ID, step.Type, step.Params, failureStepID, position)
		if err != nil {
			// Handle the error
		}
		id, err := result.LastInsertId()
		if err != nil {
			return []service.Step{}, err //TODO: proper error handling
		}

		steps[position].ID = int(id)
		steps[position].FailureStep.ID = *failureStepID
	}

	return steps, err
}

func insertFailureStep(ctx context.Context, step service.Step) (int, error) {
	panic("implement me")
}
