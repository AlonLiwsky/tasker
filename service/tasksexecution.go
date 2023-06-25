package service

import (
	"context"
	"fmt"
	"time"

	"github.com/tasker/entities"
)

const LastStepResultKey = "last_step_result"
const UseLastStepResultKey = "use_last_step_result"

// ExecuteTask
// TODO: We need to distinguish between system errors and execution errors:
// -check when to throw each one
// -when and what to save in DB
// -rollback and idempotency scenarios
// -finish implementation
func (s service) ExecuteTask(ctx context.Context, taskID int, scheduleID int) (entities.Execution, error) {
	//Get task from storage
	task, err := s.storage.GetTask(ctx, taskID)
	if err != nil {
		return entities.Execution{}, fmt.Errorf("getting task to execute: %w", err)
	}

	//Initialize execution values with success status
	exec := entities.Execution{
		Status:        entities.SuccessExecutionStatus,
		ScheduledTask: scheduleID,
		ExecutedTime:  time.Now(),
	}

	var stepResult string

	//Iterate steps one by one
	for i, step := range task.Steps {
		//Set the last step result as a param for this one
		if i > 0 {
			step.Params[LastStepResultKey] = stepResult
		}

		//Run step
		stepResult, err = s.runStep(ctx, step)
		if err != nil {
			//If it fails, check for failure steps
			if step.FailureStep != nil {
				step.FailureStep.Params[LastStepResultKey] = stepResult
				_, err = s.runStep(ctx, *step.FailureStep)
				if err == nil {
					//The failure step run successfully, we finish the execution with a handled failure status
					exec.Status = entities.HandledFailureExecutionStatus
					break
				}
			}
			//If there's no failure step, or it also failed, we fail the execution
			exec.Status = entities.FailureExecutionStatus
			break
		}
	}

	//Save execution on DB
	exec, err = s.storage.SaveExecution(ctx, exec)
	if err != nil {
		return entities.Execution{}, fmt.Errorf("saving execution: %w", err)
	}

	return exec, nil
}
