package service

import (
	"github.com/Fizik05/L0"
	"github.com/Fizik05/L0/pkg/repository"
)

type Cache interface {
	NewCache() error
	GetOrder(order_uid string) (L0.Order, error)
	AddOrder(order_uid string, order L0.Order)
}

type Service struct {
	Cache
}

func NewService(repo *repository.Repository) (*Service, error) {
	cache, err := NewCacheService(repo)
	if err != nil {
		return nil, err
	}

	return &Service{Cache: cache}, nil
}
