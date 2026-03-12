package domain

import (
	"context"
)

type TxRepository interface {
	ExecTx(ctx context.Context, fn func(
		productRepo ProductRepository,
		stockRepo StockRepository,
	) error) error
}
