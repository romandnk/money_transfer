package service

import (
	"context"
	"fmt"
	"github.com/romandnk/money_transfer/internal/models"
	"github.com/romandnk/money_transfer/internal/storage"
	"net/mail"
)

type userService struct {
	userStorage storage.User
	signKey     string
}

func newUserService(userStorage storage.User, signKey string) *userService {
	return &userService{
		userStorage: userStorage,
		signKey:     signKey,
	}
}

func (u *userService) SignUp(ctx context.Context, user models.User) (int, error) {
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return 0, fmt.Errorf("invalid email: %w", err)
	}

	if err := validatePassword(user.Password); err != nil {
		return 0, err
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return 0, err
	}

	user.Password = hashedPassword

	return u.userStorage.CreateUser(ctx, user)
}

func (u *userService) SignIn(ctx context.Context, email, password string) (string, error) {
	user, err := u.userStorage.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if !comparePassword(password, user.Password) {
		return "", ErrInvalidPassword
	}

	token, err := createJWT([]byte(u.signKey), user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
