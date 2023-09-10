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

type Account interface {
	Deposit(ctx context.Context, currency models.Currency, amount decimal.Decimal, to string) error
	Transfer(ctx context.Context, currency models.Currency, amount decimal.Decimal, userID int, to string) error
}

type Storage interface {
	Account
	User
}
