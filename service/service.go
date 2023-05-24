package service

type Storage interface {
}

type Service interface {
}

type service struct {
	storage Storage
}

func NewService(str Storage) Service {
	return service{str}
}
