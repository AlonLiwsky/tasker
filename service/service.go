package service

import "context"

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
	//Check logical integrity?
	panic("implement me")
}

func NewService(str Storage) Service {
	return service{str}
}
