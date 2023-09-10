package service

import (
	"context"
	"github.com/romandnk/money_transfer/internal/storage"
)

type AccountDepositInput struct {
	CurrencyCode string
	Amount       float64
	To           string
}

type Account interface {
	Deposit(ctx context.Context, input AccountDepositInput) error
}

type Services struct {
	Account
}

func NewServices(storage storage.Storage) *Services {
	return &Services{
		Account: newAccountService(storage),
	}
}
