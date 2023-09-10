package service

import "errors"

var (
	ErrEmptyCurrencyCode          = errors.New("currency code cannot be empty")
	ErrNotPositiveAmount          = errors.New("amount must be positive")
	ErrEmptyNumber                = errors.New("number cannot be empty")
	ErrCurrencyCodeRepresentation = errors.New("currency code must contain only uppercase letters")
	ErrNumberTooLong              = errors.New("number cannot be empty")
)

var (
	ErrPasswordTooShort               = errors.New("min password length is 6")
	ErrPasswordTooShortLong           = errors.New("min password length is 18")
	ErrPasswordWithoutDigit           = errors.New("password must contain a digit")
	ErrPasswordContainSpace           = errors.New("password cannot contain a space")
	ErrPasswordWithoutUpperCaseSymbol = errors.New("password must contain at least one uppercase symbol")
	ErrPasswordWithoutLowerCaseSymbol = errors.New("password must contain at least one lowercase symbol")
	ErrInvalidPassword                = errors.New("invalid password")
)
