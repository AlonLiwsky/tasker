package service

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tasker/entities"
)

// MockStorage is a mock implementation of the Storage interface
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) SaveTask(ctx context.Context, task entities.Task) (entities.Task, error) {
	args := m.Called(ctx, task)
	return args.Get(0).(entities.Task), args.Error(1)
}

func (m *MockStorage) GetTask(ctx context.Context, taskID int) (entities.Task, error) {
	args := m.Called(ctx, taskID)
	return args.Get(0).(entities.Task), args.Error(1)
}

func (m *MockStorage) SaveExecution(ctx context.Context, exec entities.Execution) (entities.Execution, error) {
	args := m.Called(ctx, exec)
	return args.Get(0).(entities.Execution), args.Error(1)
}

func (m *MockStorage) GetExecutionIdempotency(ctx context.Context, idempToken string) (entities.Execution, error) {
	args := m.Called(ctx, idempToken)
	return args.Get(0).(entities.Execution), args.Error(1)
}

func (m *MockStorage) SaveSchedule(ctx context.Context, sch entities.ScheduledTask) (entities.ScheduledTask, error) {
	args := m.Called(ctx, sch)
	return args.Get(0).(entities.ScheduledTask), args.Error(1)
}

func (m *MockStorage) GetEnabledSchedules(ctx context.Context) ([]entities.ScheduledTask, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entities.ScheduledTask), args.Error(1)
}

func (m *MockStorage) SetScheduleLastRun(ctx context.Context, schID int, time time.Time) error {
	args := m.Called(ctx, schID, time)
	return args.Error(0)
}
