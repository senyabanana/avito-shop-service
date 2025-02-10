package repository

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

func NewRepository() *Repository {
	return &Repository{}
}
