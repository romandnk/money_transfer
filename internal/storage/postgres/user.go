package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/romandnk/money_transfer/internal/models"
)

func (p *Postgres) CreateUser(ctx context.Context, user models.User) (int, error) {
	var id int

	sql, args, _ := p.builder.
		Insert(usersTable).
		Columns("email", "password").
		Values(user.Email, user.Password).
		Suffix("RETURNING id").
		ToSql()

	err := p.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("UserStorage - CreateUser - p.Pool.QueryRow.Scan(id): %w", err)
	}

	return id, nil
}

func (p *Postgres) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User

	sql, args, _ := p.builder.
		Select("id", "email", "password").
		From(usersTable).
		Where(squirrel.Eq{"email": email}).
		ToSql()

	err := p.Pool.QueryRow(ctx, sql, args...).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, fmt.Errorf("user not found")
		}
		return user, fmt.Errorf("UserStorage - GetUserByEmail - p.Pool.QueryRow: %w", err)
	}

	return user, nil
}
