package service

import (
	"context"
	"fmt"
	"github.com/romandnk/money_transfer/internal/storage"
)

type accountService struct {
	accountStorage storage.Account
}

func newAccountService(accountStorage storage.Account) *accountService {
	return &accountService{accountStorage: accountStorage}
}

func (a *accountService) Deposit(ctx context.Context, input AccountDepositInput) error {
	currency, err := validateCurrencyCode(input.CurrencyCode)
	if err != nil {
		return err
	}

	amount, err := validateAmount(input.Amount)
	if err != nil {
		return err
	}

	to, err := validateNumber(input.To)
	if err != nil {
		return fmt.Errorf("receiver: %w", err)
	}

	return a.accountStorage.Deposit(ctx, currency, amount, to)
}
