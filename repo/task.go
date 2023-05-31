package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/tasker/service"
)

const (
	InsertTaskQr = "INSERT INTO task (name) VALUES (?);"
	InsertStepQr = "INSERT INTO step (task_id, step_type, params, failure_step, position) VALUES (?, ?, ?, ?, ?)"
)

func (r repository) SaveTask(ctx context.Context, task service.Task) (service.Task, error) {
	result, err := r.db.ExecContext(ctx, InsertTaskQr, task.Name)
	if err != nil {
		return service.Task{}, fmt.Errorf("inserting task: %w", err)
	}
	rAffect, err := result.RowsAffected()
	if err != nil {
		return service.Task{}, err
	}
	if rAffect != 1 {
		return service.Task{}, errors.New(fmt.Sprintf("rows affected while inserting task different than 1: %v", rAffect))
	}

	auxID, err := result.LastInsertId()
	if err != nil {
		return service.Task{}, err
	}

	task.ID = int(auxID)

	steps, err := r.saveSteps(ctx, task.Steps, task.ID)
	if err != nil {
		return service.Task{}, err
	}
	task.Steps = steps

	return task, err
}

func (r repository) saveSteps(ctx context.Context, steps []service.Step, taskID int) ([]service.Step, error) {
	//Save each step with the Order field

	// Prepare the SQL statement
	stmt, err := r.db.PrepareContext(ctx, InsertStepQr)
	if err != nil {
		return []service.Step{}, fmt.Errorf("preparing insert steps stmt: %w", err)
	}
	defer stmt.Close()

	// Iterate over the rows and execute the prepared statement in a batch
	for position, step := range steps {
		var failureStep *service.Step = nil
		if step.FailureStep != nil {
			failureStep, err = r.insertFailureStep(ctx, *step.FailureStep, taskID)
			if err != nil {
				return []service.Step{}, err //TODO: proper error handling
			}
		}

		var result sql.Result = nil
		if failureStep != nil {
			result, err = stmt.Exec(taskID, step.Type, toJSON(step.Params), failureStep.ID, position)
		} else {
			result, err = stmt.Exec(taskID, step.Type, toJSON(step.Params), nil, position)
		}
		if err != nil {
			return []service.Step{}, err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return []service.Step{}, err //TODO: proper error handling
		}

		steps[position].ID = int(id)
		steps[position].FailureStep = failureStep
	}

	return steps, err
}

func (r repository) insertFailureStep(ctx context.Context, step service.Step, taskID int) (*service.Step, error) {
	step.FailureStep = nil //Only one failure step
	result, err := r.db.ExecContext(ctx, InsertStepQr, taskID, step.Type, toJSON(step.Params), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("inserting failure step: %w", err)
	}
	rAffect, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rAffect != 1 {
		return nil, errors.New(fmt.Sprintf("rows affected while inserting failure step different than 1: %v", rAffect))
	}

	auxID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	step.ID = int(auxID)

	return &step, nil
}

func toJSON(v any) string {
	jsonData, err := json.Marshal(v)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return ""
	}
	return string(jsonData)
}
