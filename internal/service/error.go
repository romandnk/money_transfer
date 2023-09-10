package service

import "errors"

var (
	ErrEmptyCurrencyCode          = errors.New("currency code cannot be empty")
	ErrNotPositiveAmount          = errors.New("amount must be positive")
	ErrEmptyNumber                = errors.New("number cannot be empty")
	ErrCurrencyCodeRepresentation = errors.New("currency code must contain only uppercase letters")
	ErrNumberTooLong              = errors.New("number cannot be empty")
)
