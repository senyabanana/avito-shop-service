package service

import (
	"github.com/senyabanana/avito-shop-service/internal/entity"
	"github.com/senyabanana/avito-shop-service/internal/repository"
	
	"github.com/sirupsen/logrus"
)

type Authorization interface {
	GetUser(username string) (entity.User, error)
	CreateUser(username, password string) error
	GenerateToken(username, password string) (string, error)
	ParseToken(accessToken string) (int, error)
}

type Transaction interface {
	GetUserInfo(userID int) (entity.InfoResponse, error)
	SendCoin(fromUserID int, toUsername string, amount int) error
}

type Inventory interface {
	BuyItem(userID int, itemName string) error
}

type Service struct {
	Authorization
	Transaction
	Inventory
}

func NewService(repos *repository.Repository, log *logrus.Logger) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, log),
		Transaction:   NewTransactionService(repos.Authorization, repos.UserTransaction, log),
		Inventory:     NewInventoryService(repos.UserTransaction, repos.Inventory, log),
	}
}
