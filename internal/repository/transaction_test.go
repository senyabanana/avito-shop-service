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

func TestTransactionPostgres_GetReceivedTransactions(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewTransactionPostgres(sqlxDB)

	tests := []struct {
		name         string
		userID       int
		mockBehavior func()
		wantError    error
		wantData     []entity.TransactionDetail
	}{
		{
			name:   "Success",
			userID: 1,
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"from_user", "amount"}).
					AddRow("user1", 100).
					AddRow("user2", 50)

				mock.ExpectQuery(`
						SELECT u.username AS from_user, t.amount FROM transactions AS t
						JOIN users AS u ON t.from_user = u.id WHERE t.to_user = \$1`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			wantError: nil,
			wantData: []entity.TransactionDetail{
				{FromUser: "user1", Amount: 100},
				{FromUser: "user2", Amount: 50},
			},
		},
		{
			name:   "Query Error",
			userID: 1,
			mockBehavior: func() {
				mock.ExpectQuery(`
						SELECT u.username AS from_user, t.amount FROM transactions AS t
						JOIN users AS u ON t.from_user = u.id WHERE t.to_user = \$1`).
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
			data, err := repo.GetReceivedTransactions(ctx, tt.userID)

			assert.Equal(t, tt.wantError, err)
			assert.Equal(t, tt.wantData, data)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTransactionPostgres_GetSentTransactions(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewTransactionPostgres(sqlxDB)

	tests := []struct {
		name         string
		userID       int
		mockBehavior func()
		wantError    error
		wantData     []entity.TransactionDetail
	}{
		{
			name:   "Success",
			userID: 1,
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"to_user", "amount"}).
					AddRow("user2", 100).
					AddRow("user3", 200)

				mock.ExpectQuery(`
						SELECT u.username AS to_user, t.amount FROM transactions AS t
						JOIN users AS u ON t.to_user = u.id WHERE t.from_user = \$1`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			wantError: nil,
			wantData: []entity.TransactionDetail{
				{ToUser: "user2", Amount: 100},
				{ToUser: "user3", Amount: 200},
			},
		},
		{
			name:   "Query Error",
			userID: 1,
			mockBehavior: func() {
				mock.ExpectQuery(`
						SELECT u.username AS to_user, t.amount FROM transactions AS t
						JOIN users AS u ON t.to_user = u.id WHERE t.from_user = \$1`).
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
			data, err := repo.GetSentTransactions(ctx, tt.userID)

			assert.Equal(t, tt.wantError, err)
			assert.Equal(t, tt.wantData, data)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTransactionPostgres_InsertTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewTransactionPostgres(sqlxDB)

	tests := []struct {
		name         string
		fromUserID   int
		toUserID     int
		amount       int
		mockBehavior func()
		wantError    error
	}{
		{
			name:       "Success",
			fromUserID: 1,
			toUserID:   2,
			amount:     100,
			mockBehavior: func() {
				mock.ExpectExec(`INSERT INTO transactions \(from_user, to_user, amount\) VALUES \(\$1, \$2, \$3\)`).
					WithArgs(1, 2, 100).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantError: nil,
		},
		{
			name:       "Query Error",
			fromUserID: 1,
			toUserID:   2,
			amount:     100,
			mockBehavior: func() {
				mock.ExpectExec(`INSERT INTO transactions \(from_user, to_user, amount\) VALUES \(\$1, \$2, \$3\)`).
					WithArgs(1, 2, 100).
					WillReturnError(errors.New("insert error"))
			},
			wantError: errors.New("insert error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			ctx := context.Background()
			err := repo.InsertTransaction(ctx, tt.fromUserID, tt.toUserID, tt.amount)

			assert.Equal(t, tt.wantError, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
