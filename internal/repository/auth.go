package repository

import (
	"fmt"

	"github.com/senyabanana/avito-shop-service/internal/database"
	"github.com/senyabanana/avito-shop-service/internal/entity"

	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user entity.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, password_hash, coins) VALUES ($1, $2, $3) RETURNING id", database.UsersTable)
	row := r.db.QueryRow(query, user.Username, user.Password, user.Coins)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username string) (entity.User, error) {
	var user entity.User
	query := fmt.Sprintf("SELECT id, username, password_hash, coins FROM %s WHERE username=$1", database.UsersTable)
	err := r.db.Get(&user, query, username)

	return user, err
}
