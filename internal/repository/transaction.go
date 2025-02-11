package repository

import (
	"fmt"

	"github.com/senyabanana/avito-shop-service/internal/database"
	"github.com/senyabanana/avito-shop-service/internal/entity"

	"github.com/jmoiron/sqlx"
)

type UserTransactionPostgres struct {
	db *sqlx.DB
}

func NewUserTransactionPostgres(db *sqlx.DB) *UserTransactionPostgres {
	return &UserTransactionPostgres{db: db}
}

func (r *UserTransactionPostgres) GetUserBalance(userID int) (int, error) {
	var balance int
	query := fmt.Sprintf("SELECT coins FROM %s WHERE id=$1", database.UsersTable)
	err := r.db.Get(&balance, query, userID)

	return balance, err
}

func (r *UserTransactionPostgres) GetUserInventory(userID int) ([]entity.InventoryItem, error) {
	var inventory []entity.InventoryItem
	query := fmt.Sprintf(`
		SELECT mi.item_type AS type, COUNT(i.id) AS quantity
		FROM %s AS i
		JOIN %s AS mi ON i.merch_id = mi.id
		WHERE i.user_id = $1
		GROUP BY mi.item_type`,
		database.InventoryTable, database.MerchItemsTable)
	err := r.db.Select(&inventory, query, userID)

	return inventory, err
}

func (r *UserTransactionPostgres) GetReceivedTransactions(userID int) ([]entity.TransactionDetail, error) {
	var received []entity.TransactionDetail
	query := fmt.Sprintf(`
		SELECT u.username AS from_user, t.amount
		FROM %s AS t
		JOIN %s AS u ON t.from_user = u.id
		WHERE t.to_user=$1`,
		database.TransactionsTable, database.UsersTable)
	err := r.db.Select(&received, query, userID)

	return received, err
}

func (r *UserTransactionPostgres) GetSentTransactions(userID int) ([]entity.TransactionDetail, error) {
	var sent []entity.TransactionDetail
	query := fmt.Sprintf(`
		SELECT u.username AS to_user, t.amount
		FROM %s AS t
		JOIN %s AS u ON t.to_user = u.id
		WHERE t.from_user=$1`,
		database.TransactionsTable, database.UsersTable)
	err := r.db.Select(&sent, query, userID)

	return sent, err
}
