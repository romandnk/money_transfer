package storage

import (
	"context"
	"github.com/romandnk/money_transfer/internal/models"
	"github.com/shopspring/decimal"
)

type Account interface {
	Deposit(ctx context.Context, currency models.Currency, amount decimal.Decimal, to string) error
}

type Storage interface {
	Account
}
