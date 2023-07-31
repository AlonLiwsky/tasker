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

//TODO: add test to check that last step result is being setted

func Test_service_ExecuteTask_AlreadyExecuted(t *testing.T) {
	mockStorage := MockStorage{}
	mockExecution := entities.Execution{ID: 1}

	mockStorage.On("GetExecutionIdempotency", mock.Anything, "idemp-token").Return(mockExecution, nil)

	srv := NewService(&mockStorage, emptyStepRunners)

	execution, err := srv.ExecuteTask(context.Background(), 1, 1, "idemp-token")

	assert.Nil(t, err)
	assert.Equal(t, mockExecution, execution)
	mockStorage.AssertExpectations(t)
}

func Test_service_ExecuteTask_GetExecutionIdempotencyError(t *testing.T) {
	mockStorage := MockStorage{}
	mockStorage.On("GetExecutionIdempotency", mock.Anything, "idemp-token").Return(entities.Execution{}, errors.New("mocked-error"))

	srv := NewService(&mockStorage, emptyStepRunners)

	execution, err := srv.ExecuteTask(context.Background(), 1, 1, "idemp-token")

	assert.ErrorContains(t, err, "checking idempotency: mocked-error")
	assert.Equal(t, entities.Execution{}, execution)
	mockStorage.AssertExpectations(t)
}

func Test_service_ExecuteTask_GetTaskError(t *testing.T) {
	mockStorage := MockStorage{}
	mockStorage.On("GetExecutionIdempotency", mock.Anything, "idemp-token").Return(entities.Execution{}, nil)
	mockStorage.On("GetTask", mock.Anything, 1).Return(entities.Task{}, errors.New("mocked-error"))

	srv := NewService(&mockStorage, emptyStepRunners)

	execution, err := srv.ExecuteTask(context.Background(), 1, 1, "idemp-token")

	assert.ErrorContains(t, err, "getting task to execute: mocked-error")
	assert.Equal(t, entities.Execution{}, execution)
	mockStorage.AssertExpectations(t)
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
		Status:           entities.SuccessExecutionStatus,
		ScheduledTask:    1,
		TaskID:           1,
		IdempotencyToken: "idemp-token",
		ExecutedTime:     time.Time{},
	}
	//SHOULD FAIL for MISSING ID IN EXPECTED
	mockStorage.On("SaveExecution", mock.Anything, expectedExecution).Return(expectedExecution, nil)

	mockStepRunner := MockStepRunner{}
	var expectedParams map[string]string
	mockStepRunner.On("RunStep", mock.Anything, expectedParams).Return("step-result", nil)
	emptyStepRunners["test"] = mockStepRunner

	srv := NewService(&mockStorage, emptyStepRunners)

	execution, err := srv.ExecuteTask(context.Background(), 1, 1, "idemp-token")

	assert.Nil(t, err)
	assert.Equal(t, expectedExecution, execution)
	mockStorage.AssertExpectations(t)
	mockStepRunner.AssertExpectations(t)
}

func Test_service_ExecuteTask_StepExecution_Failure_NoFailureStep(t *testing.T) {
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
		Status:           entities.FailureExecutionStatus,
		ScheduledTask:    1,
		TaskID:           1,
		IdempotencyToken: "idemp-token",
		ExecutedTime:     time.Time{},
	}
	mockStorage.On("SaveExecution", mock.Anything, expectedExecution).Return(expectedExecution, nil)

	mockStepRunner := MockStepRunner{}
	var expectedParams map[string]string
	mockStepRunner.On("RunStep", mock.Anything, expectedParams).Return("", errors.New("mocked runstep error"))
	emptyStepRunners["test"] = mockStepRunner

	srv := NewService(&mockStorage, emptyStepRunners)

	execution, err := srv.ExecuteTask(context.Background(), 1, 1, "idemp-token")

	assert.Nil(t, err)
	assert.Equal(t, expectedExecution, execution)
	mockStorage.AssertExpectations(t)
	mockStepRunner.AssertExpectations(t)
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
				FailureStep: &entities.Step{
					ID:     2,
					Type:   "test",
					Params: map[string]string{},
				},
			},
		},
	}
	mockStorage.On("GetTask", mock.Anything, 1).Return(task, nil)

	expectedExecution := entities.Execution{
		Status:           entities.FailureExecutionStatus,
		ScheduledTask:    1,
		TaskID:           1,
		IdempotencyToken: "idemp-token",
		ExecutedTime:     time.Time{},
	}
	mockStorage.On("SaveExecution", mock.Anything, expectedExecution).Return(expectedExecution, nil)

	mockStepRunner := MockStepRunner{}
	var expectedParams1 map[string]string
	expectedParams2 := map[string]string{
		"last_step_result": "",
	}
	mockStepRunner.On("RunStep", mock.Anything, expectedParams1).Return("", errors.New("mocked runstep error"))
	mockStepRunner.On("RunStep", mock.Anything, expectedParams2).Return("", errors.New("mocked failure step runstep error"))
	emptyStepRunners["test"] = mockStepRunner

	srv := NewService(&mockStorage, emptyStepRunners)

	execution, err := srv.ExecuteTask(context.Background(), 1, 1, "idemp-token")

	assert.Nil(t, err)
	assert.Equal(t, expectedExecution, execution)
	mockStorage.AssertExpectations(t)
	mockStepRunner.AssertExpectations(t)
}

func Test_service_ExecuteTask_StepExecution_FailureStepExecution_Success(t *testing.T) {
	mockStorage := MockStorage{}
	mockStorage.On("GetExecutionIdempotency", mock.Anything, "idemp-token").Return(entities.Execution{}, nil)

	task := entities.Task{
		ID: 1,
		Steps: []entities.Step{
			{
				ID:   1,
				Type: "test",
				FailureStep: &entities.Step{
					ID:     2,
					Type:   "test",
					Params: map[string]string{},
				},
			},
		},
	}
	mockStorage.On("GetTask", mock.Anything, 1).Return(task, nil)

	expectedExecution := entities.Execution{
		Status:           entities.HandledFailureExecutionStatus,
		ScheduledTask:    1,
		TaskID:           1,
		IdempotencyToken: "idemp-token",
		ExecutedTime:     time.Time{},
	}
	mockStorage.On("SaveExecution", mock.Anything, expectedExecution).Return(expectedExecution, nil)

	mockStepRunner := MockStepRunner{}
	var expectedParams1 map[string]string
	expectedParams2 := map[string]string{
		"last_step_result": "",
	}
	mockStepRunner.On("RunStep", mock.Anything, expectedParams1).Return("", errors.New("mocked runstep error"))
	mockStepRunner.On("RunStep", mock.Anything, expectedParams2).Return("", nil)
	emptyStepRunners["test"] = mockStepRunner

	srv := NewService(&mockStorage, emptyStepRunners)

	execution, err := srv.ExecuteTask(context.Background(), 1, 1, "idemp-token")

	assert.Nil(t, err)
	assert.Equal(t, expectedExecution, execution)
	mockStorage.AssertExpectations(t)
	mockStepRunner.AssertExpectations(t)
}

func Test_service_ExecuteTask_StepExecution_FailureSavingExecution(t *testing.T) {
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
		Status:           entities.SuccessExecutionStatus,
		ScheduledTask:    1,
		TaskID:           1,
		IdempotencyToken: "idemp-token",
		ExecutedTime:     time.Time{},
	}
	mockStorage.On("SaveExecution", mock.Anything, expectedExecution).Return(entities.Execution{}, errors.New("mocked save exec error"))

	mockStepRunner := MockStepRunner{}
	var expectedParams map[string]string
	mockStepRunner.On("RunStep", mock.Anything, expectedParams).Return("step-result", nil)
	emptyStepRunners["test"] = mockStepRunner

	srv := NewService(&mockStorage, emptyStepRunners)

	execution, err := srv.ExecuteTask(context.Background(), 1, 1, "idemp-token")

	assert.Equal(t, execution, entities.Execution{})
	assert.Error(t, err, "mocked save exec error")
	mockStorage.AssertExpectations(t)
	mockStepRunner.AssertExpectations(t)
}
