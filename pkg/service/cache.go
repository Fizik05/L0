package service

import (
	"errors"

	"github.com/Fizik05/L0"
	"github.com/Fizik05/L0/pkg/repository"
	"github.com/sirupsen/logrus"
)

type CacheService struct {
	repo  repository.Order
	cache map[string]L0.Order
}

func NewCacheService(repo *repository.Repository) (*CacheService, error) {
	cacheService := &CacheService{repo: repo}
	err := cacheService.NewCache()
	if err != nil {
		return nil, err
	}

	return cacheService, nil
}

func (c *CacheService) NewCache() error {
	cache := make(map[string]L0.Order)
	orders, err := c.repo.RecoverCache()
	if err != nil {
		return err
	}

	for _, order := range orders {
		cache[order.Order_uid] = order
	}

	c.cache = cache

	logrus.Info("Cache was UP")

	return nil
}

func (c *CacheService) GetOrder(order_uid string) (L0.Order, error) {
	order, ok := c.cache[order_uid]
	if !ok {
		return order, errors.New("order with this UID is not found")
	}
	return order, nil
}

func (c *CacheService) AddOrder(order_uid string, order L0.Order) {
	c.cache[order_uid] = order
}
