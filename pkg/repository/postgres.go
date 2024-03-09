package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	ordersTable     = "orders"
	deliveriesTable = "deliveries"
	paymentsTable   = "payments"
	itemsTable      = "items"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	query := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.DBName,
		cfg.SSLMode,
		cfg.Password,
	)

	db, err := sqlx.Connect(
		"postgres",
		query,
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
