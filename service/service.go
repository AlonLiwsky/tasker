package service

import (
	"context"
	"fmt"
)

type Storage interface {
	SaveTask(ctx context.Context, task Task) (Task, error)
}

type Service interface {
	CreateTask(ctx context.Context, task Task) (Task, error)
}

type service struct {
	storage Storage
}

func (s service) CreateTask(ctx context.Context, task Task) (Task, error) {
	if err := task.IsValid(); err != nil {
		return Task{}, err
	}

	task, err := s.storage.SaveTask(ctx, task)
	if err != nil {
		return Task{}, fmt.Errorf("creating task: %w", err)
	}

	return task, nil
}

func NewService(str Storage) Service {
	return service{str}
}
