package postgres

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/romandnk/money_transfer/internal/models"
	"github.com/shopspring/decimal"
	"time"
)

const (
	statusCreated = "Created"
	statusError   = "Error"
	statusSuccess = "Success"
)

func (p *Postgres) Deposit(ctx context.Context, currency models.Currency, amount decimal.Decimal, to string) error {
	var success bool
	now := time.Now().UTC()

	tx, err := p.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		return fmt.Errorf("AccountStorage - Deposit - p.pool.BeginTx: %w", err)
	}
	defer func() {
		// if we return from function before successful operation,
		// that means we have to add a new note about error transaction
		if !success {
			sql, args, _ := p.builder.
				Insert(transactionsTable).
				Columns("created_at", "status", "amount", "currency_code", "from_wallet", "to_wallet").
				Values(now, statusError, amount, currency.Code, nil, to).
				ToSql()

			_ = setTransactionStatusError(ctx, tx, sql, args...)
		}

		_ = tx.Rollback(ctx)
	}()

	// create note about creating transaction
	sql, args, _ := p.builder.
		Insert(transactionsTable).
		Columns("created_at", "status", "amount", "currency_code", "from_wallet", "to_wallet").
		Values(now, statusCreated, amount, currency.Code, nil, to).
		Suffix("RETURNING id").
		ToSql()

	transactionId, err := setTransactionStatusCreated(ctx, tx, sql, args...)
	if err != nil {
		return fmt.Errorf("AccountStorage - Deposit - setTransactionStatusCreated: %w", err)
	}

	// check if card/wallet number exists
	sql, args, _ = p.builder.Select("1").
		Prefix("SELECT EXISTS (").
		From(accountsTable).
		Where(squirrel.Eq{"number": to}).
		Suffix(")").ToSql()

	var exist bool
	err = tx.QueryRow(ctx, sql, args...).Scan(&exist)
	if err != nil {
		return fmt.Errorf("AccountStorage - Deposit - tx.QueryRow.Scan(exist): %w", err)
	}
	if !exist {
		return fmt.Errorf("%s doesn't exist", to)
	}

	// update balance of the account
	sql, args, _ = p.builder.
		Update(accountsTable).
		Set("balance", squirrel.Expr("balance + ?", amount)).
		Where(squirrel.Eq{"number": to}).
		Where(squirrel.Eq{"currency_code": currency.Code}).
		ToSql()

	result, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AccountStorage - Deposit - tx.Exec: %w", err)
	}
	// if no rows were updated, that means that the account doesn't have a balance in certain currency
	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s doesn't have a wallet in %s", to, currency.Code)
	}

	// update transaction note about successful transaction
	sql, args, _ = p.builder.
		Update(transactionsTable).
		Set("status", statusSuccess).
		Where(squirrel.Eq{"id": transactionId}).
		ToSql()

	err = setTransactionStatusSuccess(ctx, tx, sql, args...)
	if err != nil {
		return fmt.Errorf("AccountStorage - Deposit - setTransactionStatusSuccess: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("AccountStorage - Deposit - tx.Commit: %w", err)
	}

	success = true

	return nil
}

func setTransactionStatusCreated(ctx context.Context, tx pgx.Tx, sql string, args ...any) (int, error) {
	var transactionId int
	err := tx.QueryRow(ctx, sql, args...).Scan(&transactionId)
	if err != nil {
		return transactionId, err
	}
	return transactionId, nil
}

func setTransactionStatusSuccess(ctx context.Context, tx pgx.Tx, sql string, args ...any) error {
	_, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}

func setTransactionStatusError(ctx context.Context, tx pgx.Tx, sql string, args ...any) error {
	_, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}
