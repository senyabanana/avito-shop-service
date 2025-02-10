package repository

import "github.com/jmoiron/sqlx"

type Authorization interface {
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
	return &Repository{}
}
