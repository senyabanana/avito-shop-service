package service

import (
	"errors"

	"github.com/senyabanana/avito-shop-service/internal/entity"
	"github.com/senyabanana/avito-shop-service/internal/repository"
)

type TransactionService struct {
	authRepo repository.Authorization
	userRepo repository.UserTransaction
}

func NewTransactionService(authRepo repository.Authorization, userRepo repository.UserTransaction) *TransactionService {
	return &TransactionService{authRepo: authRepo, userRepo: userRepo}
}

func (s *TransactionService) GetUserInfo(userID int) (entity.InfoResponse, error) {
	var info entity.InfoResponse
	var err error

	info.Coins, err = s.userRepo.GetUserBalance(userID)
	if err != nil {
		return entity.InfoResponse{}, err
	}

	info.Inventory, err = s.userRepo.GetUserInventory(userID)
	if err != nil {
		return entity.InfoResponse{}, err
	}

	info.CoinHistory.Received, err = s.userRepo.GetReceivedTransactions(userID)
	if err != nil {
		return entity.InfoResponse{}, err
	}

	info.CoinHistory.Sent, err = s.userRepo.GetSentTransactions(userID)
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

func (s *TransactionService) SendCoin(fromUserID int, toUsername string, amount int) error {
	toUser, err := s.authRepo.GetUser(toUsername)
	if err != nil {
		return errors.New("recipient not found")
	}

	toUserId := toUser.ID
	if fromUserID == toUserId {
		return errors.New("cannot send coins to yourself")
	}

	balance, err := s.userRepo.GetUserBalance(fromUserID)
	if err != nil {
		return err
	}
	if balance < amount {
		return errors.New("insufficient balance")
	}

	err = s.userRepo.TransferCoins(fromUserID, toUserId, amount)
	if err != nil {
		return err
	}

	return nil
}
