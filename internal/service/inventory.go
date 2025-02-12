package service

import (
	"errors"

	"github.com/senyabanana/avito-shop-service/internal/repository"

	"github.com/sirupsen/logrus"
)

type InventoryService struct {
	userRepo repository.UserTransaction
	itemRepo repository.Inventory
	log      *logrus.Logger
}

func NewInventoryService(userRepo repository.UserTransaction, itemRepo repository.Inventory, log *logrus.Logger) *InventoryService {
	return &InventoryService{
		userRepo: userRepo,
		itemRepo: itemRepo,
		log:      log,
	}
}

func (s *InventoryService) BuyItem(userID int, itemName string) error {
	s.log.Infof("User %d is attempting to buy item: %s", userID, itemName)

	item, err := s.itemRepo.GetItem(itemName)
	if err != nil {
		s.log.Warnf("BuyItem failed: item %s not found", itemName)
		return errors.New("item not found")
	}

	balance, err := s.userRepo.GetUserBalance(userID)
	if err != nil {
		s.log.Errorf("BuyItem failed: failed to fetch balance for user %d: %v", userID, err)
		return err
	}
	if balance < item.Price {
		s.log.Warnf("BuyItem failed: insufficient balance for user %d", userID)
		return errors.New("insufficient balance")
	}

	err = s.itemRepo.PurchaseItem(userID, item.ID, item.Price)
	if err != nil {
		s.log.Errorf("BuyItem failed: error purchasing item %s for user %d: %v", itemName, userID, err)
		return err
	}

	s.log.Infof("User %d successfully purchased item: %s", userID, itemName)
	return nil
}
