package dto

import (
	"time"

	"github.com/google/uuid"
)

type OrderRequest struct {
	UserID uuid.UUID          `json:"user_id" validate:"required"`
	Items  []OrderItemRequest `json:"items" validate:"required,min=1"`
}

type OrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
}

type OrderResponse struct {
	ID          uuid.UUID           `json:"id"`
	UserID      uuid.UUID           `json:"user_id"`
	TotalAmount uint                `json:"total_amount"`
	Items       []OrderItemResponse `json:"items"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

type OrderItemResponse struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	OrderID   uuid.UUID `json:"order_id"`
	Quantity  int       `json:"quantity"`
	Price     uint      `json:"price"`
}
