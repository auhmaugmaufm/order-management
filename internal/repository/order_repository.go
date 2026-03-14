package repository

import (
	"context"

	"github.com/auhmaugmaufm/event-driven-order/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{db: db}
}

// Create implements domain.OrderRepository.
func (r *orderRepository) Create(ctx context.Context, order *domain.Order) error {
	if err := r.db.WithContext(ctx).Create(order).Error; err != nil {
		return err
	}
	return nil
}

// GetAll implements domain.OrderRepository.
func (r *orderRepository) GetAll(ctx context.Context) ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.WithContext(ctx).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// GetByID implements domain.OrderRepository.
func (r *orderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	var order *domain.Order
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}
