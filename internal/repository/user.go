package repository

import (
	"context"

	"github.com/senyabanana/avito-shop-service/internal/entity"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db     *sqlx.DB
	getter *trmsqlx.CtxGetter
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{
		db:     db,
		getter: trmsqlx.DefaultCtxGetter,
	}
}

func (r *UserPostgres) CreateUser(ctx context.Context, user entity.User) (int, error) {
	var id int
	query := `INSERT INTO users (username, password_hash, coins) VALUES ($1, $2, $3) RETURNING id`
	row := r.getter.DefaultTrOrDB(ctx, r.db).QueryRowContext(ctx, query, user.Username, user.Password, user.Coins)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserPostgres) GetUser(ctx context.Context, username string) (entity.User, error) {
	var user entity.User
	query := `SELECT id, username, users.password_hash, coins FROM users WHERE username = $1`

	return user, r.getter.DefaultTrOrDB(ctx, r.db).GetContext(ctx, &user, query, username)
}

func (r *UserPostgres) GetUserBalance(ctx context.Context, userID int) (int, error) {
	var balance int
	query := `SELECT coins FROM users WHERE id = $1`

	return balance, r.getter.DefaultTrOrDB(ctx, r.db).GetContext(ctx, &balance, query, userID)
}

func (r *UserPostgres) UpdateCoins(ctx context.Context, userID, amount int) error {
	query := `UPDATE users SET coins = coins + $1 WHERE id = $2 AND coins >= $1`

	res, err := r.getter.DefaultTrOrDB(ctx, r.db).ExecContext(ctx, query, amount, userID)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return nil
	}

	return nil
}
