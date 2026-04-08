package service

import (
	"context"
	"errors"

	"github.com/auhmaugmaufm/event-driven-order/internal/domain"
	"github.com/auhmaugmaufm/event-driven-order/internal/dto"
	"github.com/google/uuid"
)

type StockMovementService struct {
	repo      domain.StockMovementRepository
	txm       domain.TxRepository
	stockRepo domain.StockRepository
}

func NewStockMovementService(
	repo domain.StockMovementRepository,
	txm domain.TxRepository,
	stockRepo domain.StockRepository) *StockMovementService {
	return &StockMovementService{repo: repo, txm: txm, stockRepo: stockRepo}
}

func (s *StockMovementService) Create(ctx context.Context, req *dto.StockMovementRequest) error {
	stock, err := s.stockRepo.GetProductStock(ctx, req.ProductID)
	if err != nil {
		return err
	}
	return s.txm.ExecStockMovementTx(ctx, func(repo domain.StockMovementRepository, stockRepo domain.StockRepository) error {
		stockMovement := &domain.StockMovement{
			StockID:      stock.ID,
			MovementType: req.MovementType,
			Quantity:     req.Quantity,
		}
		if err := s.repo.Create(ctx, stockMovement); err != nil {
			return err
		}
		// stock := &domain.Stock{
		// 	ProductID: req.ProductID,
		// 	Quantity:  req.Quantity,
		// } --> useless
		if req.MovementType == "IN" {
			return s.stockRepo.IncreaseStockWithTx(ctx, req.ProductID, req.Quantity)
		} else {
			return s.stockRepo.DecreaseStockWithTx(ctx, req.ProductID, req.Quantity)
		}
	})
}

func (s *StockMovementService) GetByMovementID(ctx context.Context, id uuid.UUID) (*dto.StockMovementResponse, error) {
	movement, err := s.repo.FindByStockMovementID(ctx, id)
	if err != nil {
		return nil, errors.New("Stock not found")
	}
	return &dto.StockMovementResponse{
		ID:           movement.ID,
		ProductID:    movement.Stock.ProductID,
		MovementType: movement.MovementType,
		Quantity:     movement.Quantity,
		CreatedAt:    movement.CreatedAt,
		UpdatedAt:    movement.UpdatedAt,
	}, nil
}

func (s *StockMovementService) GetAllMovement(ctx context.Context, pagination *domain.Pagination) ([]domain.StockMovement, int64, error) {
	movements, total, err := s.repo.GetStockMovement(ctx, pagination)
	if err != nil {
		return nil, 0, errors.New("Movements not found")
	}
	return movements, total, nil
}

func (s *StockMovementService) GetAllMovementType(ctx context.Context, movementType string, pagination *domain.Pagination) ([]domain.StockMovement, int64, error) {
	movements, total, err := s.repo.FindByMovementType(ctx, movementType, pagination)
	if err != nil {
		return nil, 0, errors.New("Movements not found")
	}
	return movements, total, nil
}
