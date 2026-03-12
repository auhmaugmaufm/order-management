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

func (r *productRepository) GetAll(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.WithContext(ctx).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
