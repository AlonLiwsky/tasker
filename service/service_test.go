package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tasker/entities"
)

var emptyStepRunners = map[entities.StepType]StepRunner{
	entities.APICallStepType:      StepRunner(nil),
	entities.StorageReadStepType:  StepRunner(nil),
	entities.StorageWriteStepType: StepRunner(nil),
}

func TestNewService_InvalidStepRunners(t *testing.T) {
	defer func() {
		//Should panic for invalid step runners
		r := recover()
		assert.NotNil(t, r)
	}()
	NewService(nil, nil)
}

func TestNewService(t *testing.T) {
	defer func() {
		r := recover()
		assert.Nil(t, r)
	}()

	NewService(nil, emptyStepRunners)
}

func Test_service_CreateTask_SaveError(t *testing.T) {
	mockStorage := MockStorage{}
	srv := NewService(&mockStorage, emptyStepRunners)

	mockStorage.On("SaveTask", mock.Anything, entities.Task{}).Return(entities.Task{}, errors.New("mocked-err"))

	_, err := srv.CreateTask(context.Background(), entities.Task{})

	assert.ErrorContains(t, err, "mocked-err")
}

func Test_service_CreateTask(t *testing.T) {
	mockStorage := MockStorage{}
	srv := NewService(&mockStorage, emptyStepRunners)

	auxTask := entities.Task{
		Name: "test",
		Steps: []entities.Step{
			{
				ID:   1,
				Type: "test",
			},
		},
	}
	taskWithID := auxTask
	taskWithID.ID = 1
	mockStorage.On("SaveTask", mock.Anything, auxTask).Return(taskWithID, nil)

	task, err := srv.CreateTask(context.Background(), auxTask)

	assert.Nil(t, err)
	assert.Equal(t, taskWithID, task)
}

func Test_service_GetTask_Error(t *testing.T) {
	mockStorage := MockStorage{}

	srv := NewService(&mockStorage, emptyStepRunners)

	taskID := 1
	mockStorage.On("GetTask", mock.Anything, taskID).Return(entities.Task{}, errors.New("mocked-error"))

	retrievedTask, err := srv.GetTask(context.Background(), taskID)

	assert.ErrorContains(t, err, "mocked-error")
	assert.Equal(t, entities.Task{}, retrievedTask)
}

func Test_service_GetTask_Success(t *testing.T) {
	mockStorage := MockStorage{}
	srv := NewService(&mockStorage, emptyStepRunners)

	taskID := 1
	steps := []entities.Step{
		{
			ID:   1,
			Type: "API",
			Params: map[string]string{
				"url":     "https://example.com",
				"timeout": "5s",
			},
		},
		{
			ID:   2,
			Type: "StorageWrite",
			Params: map[string]string{
				"key":   "data",
				"value": "test",
			},
		},
	}
	expectedTask := entities.Task{
		ID:    taskID,
		Name:  "Test Task",
		Steps: steps,
	}
	mockStorage.On("GetTask", mock.Anything, taskID).Return(expectedTask, nil)

	retrievedTask, err := srv.GetTask(context.Background(), taskID)

	assert.Nil(t, err)
	assert.Equal(t, expectedTask, retrievedTask)
}
