package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tasker/entities"
)

func Test_service_ExecuteTask_AlreadyExecuted(t *testing.T) {
	mockStorage := MockStorage{}
	mockExecution := entities.Execution{ID: 1}

	mockStorage.On("GetExecutionIdempotency", mock.Anything, "idemp-token").Return(mockExecution, nil)

	srv := NewService(&mockStorage, nil)

	execution, err := srv.ExecuteTask(context.Background(), 1, 1, "idemp-token")

	assert.Nil(t, err)
	assert.Equal(t, mockExecution, execution)
}

func Test_service_ExecuteTask_GetExecutionIdempotencyError(t *testing.T) {
	mockStorage := MockStorage{}
	mockStorage.On("GetExecutionIdempotency", mock.Anything, "idemp-token").Return(entities.Execution{}, errors.New("mocked-error"))

	srv := NewService(&mockStorage, nil)

	execution, err := srv.ExecuteTask(context.Background(), 1, 1, "idemp-token")

	assert.ErrorContains(t, err, "checking idempotency: mocked-error")
	assert.Equal(t, entities.Execution{}, execution)
}

func Test_service_ExecuteTask_GetTaskError(t *testing.T) {
	mockStorage := MockStorage{}
	mockStorage.On("GetExecutionIdempotency", mock.Anything, "idemp-token").Return(entities.Execution{}, nil)
	mockStorage.On("GetTask", mock.Anything, 1).Return(entities.Task{}, errors.New("mocked-error"))

	srv := NewService(&mockStorage, nil)

	execution, err := srv.ExecuteTask(context.Background(), 1, 1, "idemp-token")

	assert.ErrorContains(t, err, "getting task to execute: mocked-error")
	assert.Equal(t, entities.Execution{}, execution)
}

func Test_service_ExecuteTask_StepExecution_Success(t *testing.T) {
	mockStorage := MockStorage{}
	mockStorage.On("GetExecutionIdempotency", mock.Anything, "idemp-token").Return(entities.Execution{}, nil)

	task := entities.Task{
		ID: 1,
		Steps: []entities.Step{
			{
				ID:   1,
				Type: "test",
			},
		},
	}
	mockStorage.On("GetTask", mock.Anything, 1).Return(task, nil)

	expectedExecution := entities.Execution{
		ID:               1,
		Status:           entities.SuccessExecutionStatus,
		ScheduledTask:    1,
		TaskID:           1,
		IdempotencyToken: "idemp-token",
		ExecutedTime:     time.Time{},
	}
	//SHOULD FAIL for MISSING ID IN EXPECTED
	mockStorage.On("SaveExecution", mock.Anything, expectedExecution).Return(expectedExecution, nil)

	srv := NewService(&mockStorage, nil)

	execution, err := srv.ExecuteTask(context.Background(), 1, 1, "idemp-token")

	assert.Nil(t, err)
	assert.Equal(t, expectedExecution, execution)
}

func Test_service_ExecuteTask_StepExecution_Failure(t *testing.T) {
	mockStorage := MockStorage{}
	mockStorage.On("GetExecutionIdempotency", mock.Anything, "idemp-token").Return(entities.Execution{}, nil)

	task := entities.Task{
		ID: 1,
		Steps: []entities.Step{
			{
				ID:   1,
				Type: "test",
			},
		},
		FailureStep: &entities.Step{
			ID:   2,
			Type: "failure",
		},
	}
	mockStorage.On("GetTask", mock.Anything, 1).Return(task, nil)

	expectedExecution := entities.Execution{
		ID:                1,
		Status:            entities.HandledFailureExecutionStatus,
		ScheduledTask:     1,
		TaskID:            1,
		IdempotencyToken:  "idemp-token",
		ExecutedTime:      mock.AnythingOfType("time.Time"),
		StepResults:       []string{""},
		HandledFailure:    true,
		FailureStepResult: "",
	}
	mockStorage.On("SaveExecution", mock.Anything, expectedExecution).Return(expectedExecution, nil)

	srv := NewService(&mockStorage, nil)

	execution, err := srv.ExecuteTask(context.Background(), 1, 1, "idemp-token")

	assert.Nil(t, err)
	assert.Equal(t, expectedExecution, execution)
}

func Test_service_ExecuteTask_StepExecution_FailureStepExecutionError(t *testing.T) {
	mockStorage := MockStorage{}
	mockStorage.On("GetExecutionIdempotency", mock.Anything, "idemp-token").Return(entities.Execution{}, nil)

	task := entities.Task{
		ID: 1,
		Steps: []entities.Step{
			{
				ID:   1,
				Type: "test",
			},
		},
		FailureStep: &entities.Step{
			ID:   2,
			Type: "failure",
			Params: map[string]string{
				LastStepResultKey: "",
			},
		},
	}
	mockStorage.On("GetTask", mock.Anything, 1).Return(task, nil)

	// Mock failure step execution error
	mockStepRunner := MockStepRunner{}
	mockStepRunner.On("RunStep", mock.Anything, task.FailureStep).Return("", errors.New("failure step execution error"))

	expectedExecution := entities.Execution{
		ID:                1,
		Status:            entities.FailureExecutionStatus,
		ScheduledTask:     1,
		TaskID:            1,
		IdempotencyToken:  "idemp-token",
		ExecutedTime:      mock.AnythingOfType("time.Time"),
		StepResults:       []string{""},
		HandledFailure:    false,
		FailureStepResult: "",
	}
	mockStorage.On("SaveExecution", mock.Anything, expectedExecution).Return(expectedExecution, nil)

	srv := NewService(&mockStorage, map[entities.StepType]StepRunner{
		"failure": &mockStepRunner,
	})

	execution, err := srv.ExecuteTask(context.Background(), 1, 1, "idemp-token")

	assert.Nil(t, err)
	assert.Equal(t, expectedExecution, execution)
}

func Test_service_ExecuteTask_SaveExecutionError(t *testing.T) {
	mockStorage := MockStorage{}
	mockStorage.On("GetExecutionIdempotency", mock.Anything, "idemp-token").Return(entities.Execution{}, nil)

	task := entities.Task{
		ID: 1,
		Steps: []entities.Step{
			{
				ID:   1,
				Type: "test",
			},
		},
	}
	mockStorage.On("GetTask", mock.Anything, 1).Return(task, nil)
	mockStorage.On("SaveExecution", mock.Anything, mock.AnythingOfType("entities.Execution")).Return(entities.Execution{}, errors.New("mocked-error"))

	srv := NewService(&mockStorage, nil)

	execution, err := srv.ExecuteTask(context.Background(), 1, 1, "idemp-token")

	assert.ErrorContains(t, err, "saving execution")
	assert.Equal(t, entities.Execution{}, execution)
}
