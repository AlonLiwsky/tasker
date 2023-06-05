package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/tasker/entities"
)

const (
	InsertTaskQr = "INSERT INTO task (name) VALUES (?)"
	InsertStepQr = "INSERT INTO step (task_id, step_type, params, failure_step, position) VALUES (?, ?, ?, ?, ?)"
)

func (r repository) SaveTask(ctx context.Context, task entities.Task) (savedTask entities.Task, err error) {
	ctx, err = r.db.Begin(ctx)
	if err != nil {
		return entities.Task{}, fmt.Errorf("starting task saving transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if err := r.db.Rollback(ctx); err != nil {
				log.Println("TRANSACTION ERROR: rollbacking save task tx")
			}
		}
	}()

	result, err := r.db.ExecContext(ctx, InsertTaskQr, task.Name)
	if err != nil {
		return
	}
	rAffect, err := result.RowsAffected()
	if err != nil {
		return
	}
	if rAffect != 1 {
		return
	}

	auxID, err := result.LastInsertId()
	if err != nil {
		return
	}

	task.ID = int(auxID)

	steps, err := r.saveSteps(ctx, task.Steps, task.ID)
	if err != nil {
		return
	}
	task.Steps = steps

	return task, nil
}

// saveSteps saves each steps with their order field and failure steps
func (r repository) saveSteps(ctx context.Context, steps []entities.Step, taskID int) ([]entities.Step, error) {
	// Prepare the SQL statement
	stmt, err := r.db.PrepareContext(ctx, InsertStepQr)
	if err != nil {
		return []entities.Step{}, fmt.Errorf("preparing insert steps stmt: %w", err)
	}
	defer stmt.Close()

	// Iterate over the rows and execute the prepared statement in a batch
	for position, step := range steps {
		var failureStep *entities.Step = nil
		if step.FailureStep != nil {
			failureStep, err = r.insertFailureStep(ctx, *step.FailureStep, taskID)
			if err != nil {
				return []entities.Step{}, fmt.Errorf("inserting failure step: %w", err)
			}
		}

		var result sql.Result = nil
		if failureStep != nil {
			result, err = stmt.Exec(taskID, step.Type, toJSON(step.Params), failureStep.ID, position)
		} else {
			result, err = stmt.Exec(taskID, step.Type, toJSON(step.Params), nil, position)
		}
		if err != nil {
			return []entities.Step{}, err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return []entities.Step{}, err
		}

		steps[position].ID = int(id)
		steps[position].FailureStep = failureStep
	}

	return steps, nil
}

func (r repository) insertFailureStep(ctx context.Context, step entities.Step, taskID int) (*entities.Step, error) {
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
