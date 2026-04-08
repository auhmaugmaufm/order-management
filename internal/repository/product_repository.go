package repository

import (
	"context"

	"github.com/auhmaugmaufm/event-driven-order/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func NewProductRepositoryWithTx(tx *gorm.DB) domain.ProductRepository {
	return &productRepository{db: tx}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		return err
	}
	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetAll(ctx context.Context, pagination *domain.Pagination) ([]domain.Product, int64, error) {
	page := pagination.Page
	limit := pagination.Limit

	offset := (page - 1) * limit

	var products []domain.Product
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	var totalOrder int64
	fetchTotalError := r.db.WithContext(ctx).Model(&domain.Product{}).Count(&totalOrder).Error
	if fetchTotalError != nil {
		return nil, 0, fetchTotalError
	}

	return products, totalOrder, nil
}

func (r *productRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
