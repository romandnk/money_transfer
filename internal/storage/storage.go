package storage

import (
	"context"
	"github.com/romandnk/money_transfer/internal/models"
	"github.com/shopspring/decimal"
)

type User interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type AccountGetBalanceByUserIDOutput struct {
	CurrencyCode string
	Actual       float64
	Reserved     float64
}

type Account interface {
	Deposit(ctx context.Context, currency models.Currency, amount decimal.Decimal, to string) error
	Transfer(ctx context.Context, currency models.Currency, amount decimal.Decimal, userID int, to string) error
	GetBalanceByUserID(ctx context.Context, userID int) ([]models.UserBalance, error)
}

type Storage interface {
	Account
	User
}
