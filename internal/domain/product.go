package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProductName  string    `json:"product_name" gorm:"not null"`
	ProductPrice uint      `json:"product_price" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]Product, error)
	GetAll(ctx context.Context) ([]Product, error)
}
