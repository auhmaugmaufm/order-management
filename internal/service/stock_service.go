package service

import (
	"context"
	"errors"

	"github.com/auhmaugmaufm/event-driven-order/internal/domain"
	"github.com/auhmaugmaufm/event-driven-order/internal/dto"
	"github.com/google/uuid"
)

type StockService struct {
	repo domain.StockRepository
}

func NewStockService(repo domain.StockRepository) *StockService {
	return &StockService{repo: repo}
}

// TODO: Stock Adjustment by create Stock movement

// func (s *StockService) IncreaseStock(productId uuid.UUID, quantity int) error {
// 	return s.repo.IncreaseStockWithTx(productId, quantity)
// }

// func (s *StockService) DecreaseStock(productId uuid.UUID, quantity int) error {
// 	return s.repo.DecreaseStockWithTx(productId, quantity)
// }

func (s *StockService) GetProductStock(ctx context.Context, productId uuid.UUID) (*dto.StockResponse, error) {
	stock, err := s.repo.GetProductStock(ctx, productId)
	if err != nil {
		return nil, err
	}

	productName := ""
	if stock.Product != nil {
		productName = stock.Product.ProductName
	}

	return &dto.StockResponse{
		ID:          stock.ID,
		ProductID:   stock.ProductID,
		ProductName: productName,
		Quantity:    stock.Quantity,
		CreatedAt:   stock.CreatedAt,
		UpdatedAt:   stock.UpdatedAt,
	}, nil
}

func (s *StockService) GetAll(ctx context.Context) ([]dto.StockResponse, error) {
	stocks, err := s.repo.GetStocks(ctx)
	if err != nil {
		return nil, errors.New("Stocks not found")
	}
	res := make([]dto.StockResponse, 0, len(stocks))

	for _, stock := range stocks {
		productName := ""
		if stock.Product != nil {
			productName = stock.Product.ProductName
		}
		res = append(res, dto.StockResponse{
			ID:          stock.ID,
			ProductID:   stock.ProductID,
			ProductName: productName,
			Quantity:    stock.Quantity,
			CreatedAt:   stock.CreatedAt,
			UpdatedAt:   stock.UpdatedAt,
		})
	}
	return res, nil
}
