package service

import "github.com/senyabanana/avito-shop-service/internal/repository"

type Authorization interface {
}

type Transaction interface {
}

type Inventory interface {
}

type Service struct {
	Authorization
	Transaction
	Inventory
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
