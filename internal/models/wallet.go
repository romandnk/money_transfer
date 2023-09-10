package models

import "github.com/shopspring/decimal"

type Wallet struct {
	Number   string
	Currency Currency
	Balance  decimal.Decimal
}
