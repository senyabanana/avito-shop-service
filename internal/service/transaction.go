package service

import (
	"errors"

	"github.com/senyabanana/avito-shop-service/internal/entity"
	"github.com/senyabanana/avito-shop-service/internal/repository"
	
	"github.com/sirupsen/logrus"
)

type TransactionService struct {
	authRepo repository.Authorization
	userRepo repository.UserTransaction
	log      *logrus.Logger
}

func NewTransactionService(authRepo repository.Authorization, userRepo repository.UserTransaction, log *logrus.Logger) *TransactionService {
	return &TransactionService{
		authRepo: authRepo,
		userRepo: userRepo,
		log:      log,
	}
}

func (s *TransactionService) GetUserInfo(userID int) (entity.InfoResponse, error) {
	s.log.Infof("Fetching user info for userID: %d", userID)

	var info entity.InfoResponse
	var err error

	info.Coins, err = s.userRepo.GetUserBalance(userID)
	if err != nil {
		s.log.Errorf("Failed to get user balance for userID %d: %v", userID, err)
		return entity.InfoResponse{}, err
	}

	info.Inventory, err = s.userRepo.GetUserInventory(userID)
	if err != nil {
		s.log.Errorf("Failed to get user inventory for userID %d: %v", userID, err)
		return entity.InfoResponse{}, err
	}

	info.CoinHistory.Received, err = s.userRepo.GetReceivedTransactions(userID)
	if err != nil {
		s.log.Errorf("Failed to get received transactions for userID %d: %v", userID, err)
		return entity.InfoResponse{}, err
	}

	info.CoinHistory.Sent, err = s.userRepo.GetSentTransactions(userID)
	if err != nil {
		s.log.Errorf("Failed to get sent transactions for userID %d: %v", userID, err)
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

	s.log.Infof("Successfully fetched user info for userID: %d", userID)
	return info, nil
}

func (s *TransactionService) SendCoin(fromUserID int, toUsername string, amount int) error {
	s.log.Infof("User %d is sending %d coins to %s", fromUserID, amount, toUsername) // ✅ Логируем начало операции

	toUser, err := s.authRepo.GetUser(toUsername)
	if err != nil {
		s.log.Warnf("SendCoin failed: recipient %s not found", toUsername)
		return errors.New("recipient not found")
	}

	toUserId := toUser.ID
	if fromUserID == toUserId {
		s.log.Warnf("SendCoin failed: user %d tried to send coins to themselves", fromUserID)
		return errors.New("cannot send coins to yourself")
	}

	balance, err := s.userRepo.GetUserBalance(fromUserID)
	if err != nil {
		s.log.Errorf("SendCoin failed: failed to fetch balance for user %d: %v", fromUserID, err)
		return err
	}
	if balance < amount {
		s.log.Warnf("SendCoin failed: insufficient balance for user %d", fromUserID)
		return errors.New("insufficient balance")
	}

	err = s.userRepo.TransferCoins(fromUserID, toUserId, amount)
	if err != nil {
		s.log.Errorf("SendCoin failed: error transferring coins from user %d to user %d: %v", fromUserID, toUserId, err)
		return err
	}

	s.log.Infof("Transaction successful: %d coins from %d to %s", amount, fromUserID, toUsername)
	return nil
}
