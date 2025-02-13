package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/senyabanana/avito-shop-service/internal/entity"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestInventoryPostgres_GetItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewInventoryPostgres(sqlxDB)

	tests := []struct {
		name         string
		itemName     string
		mockBehavior func()
		wantError    error
		wantData     entity.MerchItems
	}{
		{
			name:     "Success",
			itemName: "t-shirt",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "item_type", "price"}).
					AddRow(1, "t-shirt", 80)

				mock.ExpectQuery(`SELECT id, item_type, price FROM merch_items WHERE item_type = \$1`).
					WithArgs("t-shirt").
					WillReturnRows(rows)
			},
			wantError: nil,
			wantData:  entity.MerchItems{ID: 1, ItemType: "t-shirt", Price: 80},
		},
		{
			name:     "Query Error",
			itemName: "t-shirt",
			mockBehavior: func() {
				mock.ExpectQuery(`SELECT id, item_type, price FROM merch_items WHERE item_type = \$1`).
					WithArgs("t-shirt").
					WillReturnError(errors.New("query error"))
			},
			wantError: errors.New("query error"),
			wantData:  entity.MerchItems{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			ctx := context.Background()
			data, err := repo.GetItem(ctx, tt.itemName)

			assert.Equal(t, tt.wantError, err)
			assert.Equal(t, tt.wantData, data)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestInventoryPostgres_GetUserInventory(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewInventoryPostgres(sqlxDB)

	tests := []struct {
		name         string
		userID       int
		mockBehavior func()
		wantError    error
		wantData     []entity.InventoryItem
	}{
		{
			name:   "Success",
			userID: 1,
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"type", "quantity"}).
					AddRow("t-shirt", 2).
					AddRow("cup", 1)

				mock.ExpectQuery(`
						SELECT mi.item_type AS type, COALESCE\(SUM\(i.quantity\), 0\) AS quantity
						FROM inventory AS i
						JOIN merch_items AS mi ON i.merch_id = mi.id
						WHERE i.user_id = \$1 GROUP BY mi.item_type`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			wantError: nil,
			wantData: []entity.InventoryItem{
				{Type: "t-shirt", Quantity: 2},
				{Type: "cup", Quantity: 1},
			},
		},
		{
			name:   "Query Error",
			userID: 1,
			mockBehavior: func() {
				mock.ExpectQuery(`
						SELECT mi.item_type AS type, COALESCE\(SUM\(i.quantity\), 0\) AS quantity
						FROM inventory AS i
						JOIN merch_items AS mi ON i.merch_id = mi.id
						WHERE i.user_id = \$1 GROUP BY mi.item_type`).
					WithArgs(1).
					WillReturnError(errors.New("query error"))
			},
			wantError: errors.New("query error"),
			wantData:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			ctx := context.Background()
			data, err := repo.GetUserInventory(ctx, tt.userID)

			assert.Equal(t, tt.wantError, err)
			assert.Equal(t, tt.wantData, data)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestInventoryPostgres_GetInventoryItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewInventoryPostgres(sqlxDB)

	tests := []struct {
		name         string
		userID       int
		merchID      int
		mockBehavior func()
		wantError    error
		wantData     int
	}{
		{
			name:    "Success",
			userID:  1,
			merchID: 2,
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"quantity"}).AddRow(3)

				mock.ExpectQuery(`SELECT quantity FROM inventory WHERE user_id = \$1 AND merch_id = \$2`).
					WithArgs(1, 2).
					WillReturnRows(rows)
			},
			wantError: nil,
			wantData:  3,
		},
		{
			name:    "Query Error",
			userID:  1,
			merchID: 2,
			mockBehavior: func() {
				mock.ExpectQuery(`SELECT quantity FROM inventory WHERE user_id = \$1 AND merch_id = \$2`).
					WithArgs(1, 2).
					WillReturnError(errors.New("query error"))
			},
			wantError: errors.New("query error"),
			wantData:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			ctx := context.Background()
			data, err := repo.GetInventoryItem(ctx, tt.userID, tt.merchID)

			assert.Equal(t, tt.wantError, err)
			assert.Equal(t, tt.wantData, data)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestInventoryPostgres_UpdateInventoryItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewInventoryPostgres(sqlxDB)

	tests := []struct {
		name         string
		userID       int
		merchID      int
		mockBehavior func()
		wantError    error
	}{
		{
			name:    "Success",
			userID:  1,
			merchID: 2,
			mockBehavior: func() {
				mock.ExpectExec(`UPDATE inventory SET quantity = quantity \+ 1 WHERE user_id = \$1 AND merch_id = \$2`).
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1)) // Успешное обновление одной строки
			},
			wantError: nil,
		},
		{
			name:    "Query Error",
			userID:  1,
			merchID: 2,
			mockBehavior: func() {
				mock.ExpectExec(`UPDATE inventory SET quantity = quantity \+ 1 WHERE user_id = \$1 AND merch_id = \$2`).
					WithArgs(1, 2).
					WillReturnError(errors.New("query error"))
			},
			wantError: errors.New("query error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			ctx := context.Background()
			err := repo.UpdateInventoryItem(ctx, tt.userID, tt.merchID)

			assert.Equal(t, tt.wantError, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestInventoryPostgres_InsertInventoryItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewInventoryPostgres(sqlxDB)

	tests := []struct {
		name         string
		userID       int
		merchID      int
		mockBehavior func()
		wantError    error
	}{
		{
			name:    "Success",
			userID:  1,
			merchID: 2,
			mockBehavior: func() {
				mock.ExpectExec(`INSERT INTO inventory \(user_id, merch_id, quantity\) VALUES \(\$1, \$2, 1\)`).
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1)) // Успешная вставка одной строки
			},
			wantError: nil,
		},
		{
			name:    "Query Error",
			userID:  1,
			merchID: 2,
			mockBehavior: func() {
				mock.ExpectExec(`INSERT INTO inventory \(user_id, merch_id, quantity\) VALUES \(\$1, \$2, 1\)`).
					WithArgs(1, 2).
					WillReturnError(errors.New("query error"))
			},
			wantError: errors.New("query error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			ctx := context.Background()
			err := repo.InsertInventoryItem(ctx, tt.userID, tt.merchID)

			assert.Equal(t, tt.wantError, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
