package service

import (
	"context"
	"fmt"

	"github.com/tasker/entities"
)

type Storage interface {
	SaveTask(ctx context.Context, task entities.Task) (entities.Task, error)
}

type Service interface {
	CreateTask(ctx context.Context, task entities.Task) (entities.Task, error)
}

type service struct {
	storage Storage
}

func (s service) CreateTask(ctx context.Context, task entities.Task) (entities.Task, error) {
	task, err := s.storage.SaveTask(ctx, task)
	if err != nil {
		return entities.Task{}, fmt.Errorf("saving task: %w", err)
	}

	return task, nil
}

func NewService(str Storage) Service {
	return service{str}
}
