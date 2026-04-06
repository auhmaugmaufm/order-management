package service

import (
	"context"
	"errors"

	"github.com/auhmaugmaufm/event-driven-order/internal/domain"
	"github.com/auhmaugmaufm/event-driven-order/internal/dto"
	"github.com/google/uuid"
)

type OrderService struct {
	repo        domain.OrderRepository
	productRepo domain.ProductRepository
}

func NewOrderService(repo domain.OrderRepository, productRepo domain.ProductRepository) *OrderService {
	return &OrderService{repo: repo, productRepo: productRepo}
}

func (s *OrderService) Create(ctx context.Context, req *dto.OrderRequest) error {

	productsIDs := make([]uuid.UUID, len(req.Items))
	for i, item := range req.Items {
		productsIDs[i] = item.ProductID
	}
	products, err := s.productRepo.GetByIDs(ctx, productsIDs)
	if err != nil {
		return err
	}

	productMap := make(map[uuid.UUID]*domain.Product, len(products))
	for i := range products {
		productMap[products[i].ID] = &products[i]
	}

	var totalAmount uint
	items := make([]domain.OrderItem, len(req.Items))

	for i, item := range req.Items {
		product, ok := productMap[item.ProductID]
		if !ok {
			return errors.New("Products not found")
		}

		totalAmount += product.ProductPrice * uint(item.Quantity)
		items[i] = domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.ProductPrice,
		}
	}

	order := &domain.Order{
		UserID:      req.UserID,
		TotalAmount: totalAmount,
		Items:       items,
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *OrderService) GetByID(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error) {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	items := make([]dto.OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		items[i] = dto.OrderItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			OrderID:   item.OrderID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	return &dto.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		TotalAmount: order.TotalAmount,
		Items:       items,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}, nil
}

func (s *OrderService) GetAll(ctx context.Context) ([]dto.OrderResponse, error) {
	orders, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]dto.OrderResponse, len(orders))
	for i, order := range orders {
		items := make([]dto.OrderItemResponse, len(order.Items))
		for j, item := range order.Items {
			items[j] = dto.OrderItemResponse{
				ID:        item.ID,
				ProductID: item.ProductID,
				OrderID:   item.OrderID,
				Quantity:  item.Quantity,
				Price:     item.Price,
			}
		}

		results[i] = dto.OrderResponse{
			ID:          order.ID,
			UserID:      order.UserID,
			TotalAmount: order.TotalAmount,
			Items:       items,
			CreatedAt:   order.CreatedAt,
			UpdatedAt:   order.UpdatedAt,
		}
	}

	return results, nil
}
