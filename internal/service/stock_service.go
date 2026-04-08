package service

import (
	"context"
	"errors"

	"github.com/auhmaugmaufm/event-driven-order/internal/domain"
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

func (s *StockService) GetProductStock(ctx context.Context, productId uuid.UUID) (*domain.Stock, error) {
	stock, err := s.repo.GetProductStock(ctx, productId)
	if err != nil {
		return nil, err
	}
	return stock, nil
}

func (s *StockService) GetAll(ctx context.Context, pagination *domain.Pagination) ([]domain.Stock, int64, error) {
	stocks, total, err := s.repo.GetStocks(ctx, pagination)
	if err != nil {
		return nil, 0, errors.New("Stocks not found")
	}

	return stocks, total, nil
}
