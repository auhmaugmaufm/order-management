package repository

import (
	"context"

	"github.com/auhmaugmaufm/event-driven-order/internal/domain"
	"github.com/auhmaugmaufm/event-driven-order/internal/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *domain.Order) error {
	if err := r.db.WithContext(ctx).Create(order).Error; err != nil {
		return err
	}
	return nil
}

func (r *orderRepository) GetAll(ctx context.Context, pagination *dto.PaginationRequest) ([]domain.Order, int64, error) {
	var orders []domain.Order
	page := pagination.Page
	limit := pagination.Limit

	offset := (page - 1) * limit

	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at DESC").Preload("Items").Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	var totalOrder int64
	fetchTotalError := r.db.WithContext(ctx).Model(&domain.Order{}).Count(&totalOrder).Error
	if fetchTotalError != nil {
		return nil, 0, err
	}

	return orders, totalOrder, nil
}

func (r *orderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	var order *domain.Order
	err := r.db.WithContext(ctx).Preload("Items").Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}
