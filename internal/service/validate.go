package service

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/romandnk/money_transfer/internal/models"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
	"unicode"
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

func validatePassword(password string) error {
	if utf8.RuneCountInString(password) < 6 {
		return ErrPasswordTooShort
	}

	if utf8.RuneCountInString(password) > 18 {
		return ErrPasswordTooShortLong
	}

	if !strings.ContainsAny(password, "0123456789") {
		return ErrPasswordWithoutDigit
	}

	if strings.ContainsAny(password, " ") {
		return ErrPasswordContainSpace
	}

	if !containUpperCaseSymbol(password) {
		return ErrPasswordWithoutUpperCaseSymbol
	}

	if !containLowerCaseSymbol(password) {
		return ErrPasswordWithoutLowerCaseSymbol
	}

	return nil
}

func containUpperCaseSymbol(str string) bool {
	for _, i := range str {
		if unicode.IsUpper(i) {
			return true
		}
	}
	return false
}

func containLowerCaseSymbol(str string) bool {
	for _, i := range str {
		if unicode.IsLower(i) {
			return true
		}
	}
	return false
}

func hashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func comparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func createJWT(signKey []byte, userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Hour).Unix()
	claims["user_id"] = userID

	tokenStr, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
