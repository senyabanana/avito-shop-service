package repository

import (
	"errors"
	"fmt"

	"github.com/senyabanana/avito-shop-service/internal/database"
	"github.com/senyabanana/avito-shop-service/internal/entity"
	
	"github.com/jmoiron/sqlx"
)

type InventoryPostgres struct {
	db *sqlx.DB
}

func NewInventoryPostgres(db *sqlx.DB) *InventoryPostgres {
	return &InventoryPostgres{db: db}
}

func (r *InventoryPostgres) GetItem(itemName string) (entity.MerchItems, error) {
	var item entity.MerchItems
	query := fmt.Sprintf("SELECT id, item_type, price FROM %s WHERE item_type=$1", database.MerchItemsTable)
	err := r.db.Get(&item, query, itemName)

	return item, err
}

func (r *InventoryPostgres) PurchaseItem(userID, merchID, price int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	createSenderQuery := fmt.Sprintf("UPDATE %s SET coins = coins-$1 WHERE id=$2 AND coins>=$1", database.UsersTable)
	res, err := tx.Exec(createSenderQuery, price, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		tx.Rollback()
		return errors.New("insufficient balance")
	}

	createInventoryQuery := fmt.Sprintf("INSERT INTO %s (user_id, merch_id) VALUES ($1, $2)", database.InventoryTable)
	_, err = tx.Exec(createInventoryQuery, userID, merchID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
