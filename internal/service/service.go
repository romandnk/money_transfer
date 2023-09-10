package service

import (
	"context"
	"github.com/romandnk/money_transfer/internal/models"
	"github.com/romandnk/money_transfer/internal/storage"
)

type AccountDepositInput struct {
	CurrencyCode string
	Amount       float64
	To           string
}

type User interface {
	SignUp(ctx context.Context, user models.User) (int, error)
	SignIn(ctx context.Context, email, password string) (string, error)
}

type Account interface {
	Deposit(ctx context.Context, input AccountDepositInput) error
}

type Services struct {
	Account
	User
}

func NewServices(storage storage.Storage, salt string) *Services {
	return &Services{
		Account: newAccountService(storage),
		User:    newUserService(storage, salt),
	}
}
