package store

import (
	"context"
	"database/sql"
	"order-processing/internal/models"
)

type Storage struct {
	Orders interface {
		CreateOrder(context.Context, *models.Order) error
		GetPendingOrders(context.Context) ([]models.Order, error)
		GetOrder(context.Context, string) (models.Order, error)
		UpdateOrder(context.Context, models.Order) error
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Orders: &OrderStorage{db: db},
	}
}
