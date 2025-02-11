package service

import (
	"github.com/senyabanana/avito-shop-service/internal/entity"
	"github.com/senyabanana/avito-shop-service/internal/repository"
)

type Authorization interface {
	GetUser(username string) (entity.User, error)
	CreateUser(username, password string) error
	GenerateToken(username, password string) (string, error)
	ParseToken(accessToken string) (int, error)
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
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
