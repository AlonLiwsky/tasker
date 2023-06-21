package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tasker/entities"
	"github.com/tasker/http"
)

const (
	InsertTaskQr = "INSERT INTO task (name) VALUES (?)"
	InsertStepQr = "INSERT INTO step (task_id, step_type, params, failure_step, position) VALUES (?, ?, ?, ?, ?)"
	GetTaskQr    = "SELECT * FROM task WHERE id = ?"
	GetStepsQr   = "SELECT id, step_type, params, failure_step, position FROM step WHERE task_id = ? ORDER BY position"
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
	switch {
	case err != nil:
		return
	case rAffect != 1:
		return entities.Task{}, fmt.Errorf("inserting task: should affect 1 and affected #%d rows", rAffect)
	}

	taskID, err := result.LastInsertId()
	if err != nil {
		return
	}

	steps, err := r.saveSteps(ctx, task.Steps, int(taskID))
	if err != nil {
		return
	}

	if err = r.db.Commit(ctx); err != nil {
		return
	}

	task.ID = int(taskID)
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

	// Iterate over the steps and insert one by one with their failure steps and positions
	for position, step := range steps {
		// if the step has a failure step, we insert it first to then link them through foreign key
		var failureStep *entities.Step = nil
		if step.FailureStep != nil {
			failureStep, err = r.insertFailureStep(ctx, *step.FailureStep, taskID)
			if err != nil {
				return []entities.Step{}, fmt.Errorf("inserting failure step: %w", err)
			}
		}

		// we insert it with a foreign key to a failure step depending on it existence
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
	step.FailureStep = nil //Only one failure step, nested failure steps are not allowed

	//Failure steps has a position NULL to differentiate them from normal steps
	result, err := r.db.ExecContext(ctx, InsertStepQr, taskID, step.Type, toJSON(step.Params), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("inserting failure step: %w", err)
	}

	rAffect, err := result.RowsAffected()
	switch {
	case err != nil:
		return nil, err
	case rAffect != 1:
		return nil, fmt.Errorf("inserting failure step: should affect 1 and affected #%d rows", rAffect)
	}

	auxID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	step.ID = int(auxID)
	return &step, nil
}

func (r repository) GetTask(ctx context.Context, taskID int) (entities.Task, error) {
	task := entities.Task{}
	err := r.db.QueryRowContext(ctx, GetTaskQr, taskID).Scan(&task.ID, &task.Name)
	switch {
	case err == sql.ErrNoRows:
		return entities.Task{}, http.WrapError(err, http.ErrNotFound.WithMessage("task not found"))
	case err != nil:
		return entities.Task{}, fmt.Errorf("getting task: %w", err)
	}

	steps, err := r.getSteps(ctx, taskID)
	if err != nil {
		return entities.Task{}, fmt.Errorf("getting steps: %w", err)
	}
	task.Steps = steps

	if err := task.IsValid(); err != nil {
		return entities.Task{}, fmt.Errorf("validating read task: %w", err)
	}

	return task, nil
}

func (r repository) getSteps(ctx context.Context, taskID int) ([]entities.Step, error) {
	rows, err := r.db.QueryContext(ctx, GetStepsQr, taskID)
	if err != nil {
		return nil, fmt.Errorf("getting steps from DB: %w", err)
	}

	failureSteps := map[int]dbStep{}
	var steps []entities.Step
	for rows.Next() {
		DBStep := dbStep{}
		var jsonParams []byte
		if err := rows.Scan(&DBStep.ID, &DBStep.Type, &jsonParams, &DBStep.FailureStep, &DBStep.Position); err != nil {
			return nil, fmt.Errorf("scanning step: %w", err)
		}

		if err = json.Unmarshal(jsonParams, &DBStep.Params); err != nil {
			log.Printf("Error unmarshalling JSON: %s. The steps params got corrupted on the DB", err)
		}

		//check if it is a failure step of the task
		if DBStep.Position == nil {
			failureSteps[DBStep.ID] = DBStep
			continue
		}

		step := DBStep.toStep()
		//check if it has a failure step
		if DBStep.FailureStep != nil {
			//search it from the map of failure steps of the task
			fDBStep, found := failureSteps[*DBStep.FailureStep]
			if !found {
				return nil, fmt.Errorf("connecting failure steps, could not find failure step %d from step %d", *DBStep.FailureStep, DBStep.ID)
			}
			//parse it to a real step
			fStep := fDBStep.toStep()
			//link it to its parent step
			step.FailureStep = &fStep
		}

		steps = append(steps, step)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return steps, nil
}

func toJSON(v any) string {
	jsonData, err := json.Marshal(v)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return ""
	}
	return string(jsonData)
}
