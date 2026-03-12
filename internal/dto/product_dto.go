package dto

import (
	"time"

	"github.com/google/uuid"
)

type ProductRequest struct {
	ProductName  string `json:"product_name" validate:"required"`
	ProductPrice uint   `json:"product_price" validate:"required"`
}

type ProductResponse struct {
	ID           uuid.UUID `json:"id"`
	ProductName  string    `json:"product_name"`
	ProductPrice uint      `json:"product_price"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
