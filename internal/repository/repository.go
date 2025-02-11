package repository

import (
	"github.com/senyabanana/avito-shop-service/internal/entity"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user entity.User) (int, error)
	GetUser(username string) (entity.User, error)
}

type Transaction interface {
}

type Inventory interface {
}

type Repository struct {
	Authorization
	Transaction
	Inventory
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
