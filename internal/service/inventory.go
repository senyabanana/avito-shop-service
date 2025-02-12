package service

import (
	"errors"

	"github.com/senyabanana/avito-shop-service/internal/repository"
)

type InventoryService struct {
	userRepo repository.UserTransaction
	itemRepo repository.Inventory
}

func NewInventoryService(userRepo repository.UserTransaction, itemRepo repository.Inventory) *InventoryService {
	return &InventoryService{userRepo: userRepo, itemRepo: itemRepo}
}

func (s *InventoryService) BuyItem(userID int, itemName string) error {
	item, err := s.itemRepo.GetItem(itemName)
	if err != nil {
		return errors.New("item not found")
	}

	balance, err := s.userRepo.GetUserBalance(userID)
	if err != nil {
		return err
	}
	if balance < item.Price {
		return errors.New("insufficient balance")
	}

	err = s.itemRepo.PurchaseItem(userID, item.ID, item.Price)
	if err != nil {
		return err
	}

	return nil
}
