package repository

import (
	"github.com/Fizik05/L0"
	"github.com/jmoiron/sqlx"
)

type Order interface {
	SaveOrder(order L0.Order) error
	RecoverCache() ([]L0.Order, error)
}

type Repository struct {
	Order
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Order: NewOrderPostgres(db),
	}
}
