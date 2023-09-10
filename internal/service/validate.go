package service

import (
	"github.com/romandnk/money_transfer/internal/models"
	"github.com/shopspring/decimal"
	"strings"
	"unicode/utf8"
)

func validateCurrencyCode(code string) (models.Currency, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return models.Currency{}, ErrEmptyCurrencyCode
	}

	if strings.ToUpper(code) != code {
		return models.Currency{}, ErrCurrencyCodeRepresentation
	}

	return models.Currency{
		Code: code,
	}, nil
}

func validateAmount(amount float64) (decimal.Decimal, error) {
	if amount <= 0 {
		return decimal.Decimal{}, ErrNotPositiveAmount
	}

	decimalAmount := decimal.NewFromFloat(amount)

	return decimalAmount, nil
}

func validateNumber(number string) (string, error) {
	number = strings.TrimSpace(number)
	if number == "" {
		return "", ErrEmptyNumber
	}

	if utf8.RuneCountInString(number) > 42 {
		return "", ErrNumberTooLong
	}

	return number, nil
}
