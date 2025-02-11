package service

import (
	"github.com/senyabanana/avito-shop-service/internal/entity"
	"github.com/senyabanana/avito-shop-service/internal/repository"
)

type UserTransactionService struct {
	repo repository.UserTransaction
}

func NewTransactionService(repo repository.UserTransaction) *UserTransactionService {
	return &UserTransactionService{repo: repo}
}

func (s *UserTransactionService) GetUserInfo(userID int) (entity.InfoResponse, error) {
	var info entity.InfoResponse
	var err error

	info.Coins, err = s.repo.GetUserBalance(userID)
	if err != nil {
		return entity.InfoResponse{}, err
	}

	info.Inventory, err = s.repo.GetUserInventory(userID)
	if err != nil {
		return entity.InfoResponse{}, err
	}

	info.CoinHistory.Received, err = s.repo.GetReceivedTransactions(userID)
	if err != nil {
		return entity.InfoResponse{}, err
	}

	info.CoinHistory.Sent, err = s.repo.GetSentTransactions(userID)
	if err != nil {
		return entity.InfoResponse{}, err
	}

	if info.Inventory == nil {
		info.Inventory = make([]entity.InventoryItem, 0)
	}
	if info.CoinHistory.Received == nil {
		info.CoinHistory.Received = make([]entity.TransactionDetail, 0)
	}

	if info.CoinHistory.Sent == nil {
		info.CoinHistory.Sent = make([]entity.TransactionDetail, 0)
	}

	return info, nil
}
