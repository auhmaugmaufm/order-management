package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	TotalAmount uint      `json:"total_amount" gorm:"not null;default:0"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Items []OrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
	User  *User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type OrderItem struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	OrderID   uuid.UUID `json:"order_id" gorm:"type:uuid;not null;index"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	Price     uint      `json:"price" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Order   *Order   `json:"-" gorm:"foreignKey:OrderID"`
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*Order, error)
	GetAll(ctx context.Context) ([]Order, error)
}
