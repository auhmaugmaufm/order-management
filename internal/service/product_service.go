package service

import (
	"context"
	"errors"

	"github.com/auhmaugmaufm/event-driven-order/internal/domain"
	"github.com/auhmaugmaufm/event-driven-order/internal/dto"
	"github.com/google/uuid"
)

type ProductService struct {
	repo      domain.ProductRepository
	txm       domain.TxRepository
	stockRepo domain.StockRepository
}

func NewProductService(repo domain.ProductRepository, txm domain.TxRepository, stockRepo domain.StockRepository) *ProductService {
	return &ProductService{repo: repo, txm: txm, stockRepo: stockRepo}
}

func (s *ProductService) Create(ctx context.Context, req *dto.ProductRequest) error {
	return s.txm.ExecTx(ctx, func(repo domain.ProductRepository,
		stockRepo domain.StockRepository) error {
		product := &domain.Product{
			ProductName:  req.ProductName,
			ProductPrice: req.ProductPrice,
		}
		if err := s.repo.Create(ctx, product); err != nil {
			return err
		}
		stock := &domain.Stock{
			ProductID: product.ID,
			Quantity:  0,
		}
		return s.stockRepo.Create(ctx, stock)
	})
}

func (s *ProductService) GetByID(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("Product not found")
	}
	return &dto.ProductResponse{
		ID:           product.ID,
		ProductName:  product.ProductName,
		ProductPrice: product.ProductPrice,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}, nil
}

func (s *ProductService) GetAll(ctx context.Context) ([]dto.ProductResponse, error) {
	products, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, errors.New("Products not found")
	}
	res := make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		res = append(res, dto.ProductResponse{
			ID:           product.ID,
			ProductName:  product.ProductName,
			ProductPrice: product.ProductPrice,
			CreatedAt:    product.CreatedAt,
			UpdatedAt:    product.UpdatedAt,
		})
	}
	return res, nil
}
