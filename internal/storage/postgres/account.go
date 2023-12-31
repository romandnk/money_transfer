package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
				Columns("created_at", "status", "amount", "currency_code", "from_account", "to_account").
				Values(now, statusError, amount, currency.Code, nil, to).
				ToSql()

			_ = setTransactionStatusError(ctx, p.Pool, sql, args...)
		}

		_ = tx.Rollback(ctx)
	}()

	// create note about creating transaction
	sql, args, _ := p.builder.
		Insert(transactionsTable).
		Columns("created_at", "status", "amount", "currency_code", "from_account", "to_account").
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

func setTransactionStatusError(ctx context.Context, pool PgxPool, sql string, args ...any) error {
	_, err := pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) Transfer(ctx context.Context, currency models.Currency, amount decimal.Decimal, userID int, to string) error {
	var success bool
	var from string
	now := time.Now().UTC()

	tx, err := p.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		return fmt.Errorf("AccountStorage - Transfer - p.pool.BeginTx: %w", err)
	}
	defer func() {
		// if we return from function before successful operation,
		// that means we have to add a new note about error transaction
		if !success {
			sql, args, _ := p.builder.
				Insert(transactionsTable).
				Columns("created_at", "status", "amount", "currency_code", "from_account", "to_account").
				Values(now, statusError, amount, currency.Code, from, to).
				ToSql()

			_ = setTransactionStatusError(ctx, p.Pool, sql, args...)
		}

		_ = tx.Rollback(ctx)
	}()

	sql, args, _ := p.builder.
		Select("number").
		From(accountsTable).
		Where(squirrel.Eq{"user_id": userID}).
		Where(squirrel.Eq{"currency_code": currency.Code}).
		ToSql()

	err = tx.QueryRow(ctx, sql, args...).Scan(&from)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("user %d doesn't have an account in %s", userID, currency.Code)
		}
		return fmt.Errorf("AccountStorage - Transfer - tx.QueryRow.Scan(from): %w", err)
	}

	// create note about creating transaction
	sql, args, _ = p.builder.
		Insert(transactionsTable).
		Columns("created_at", "status", "amount", "currency_code", "from_account", "to_account").
		Values(now, statusCreated, amount, currency.Code, from, to).
		Suffix("RETURNING id").
		ToSql()

	transactionId, err := setTransactionStatusCreated(ctx, tx, sql, args...)
	if err != nil {
		return fmt.Errorf("AccountStorage - Transfer - setTransactionStatusCreated: %w", err)
	}

	sql, args, _ = p.builder.
		Update(accountsTable).
		Set("balance", squirrel.Expr("balance - ?", amount)).
		Where(squirrel.Eq{"number": from}).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23514" {
				return fmt.Errorf("not enough money")
			}
		}
		return fmt.Errorf("AccountStorage - Transfer - tx.Exec: %v", err)
	}

	sql, args, _ = p.builder.
		Update(accountsTable).
		Set("balance", squirrel.Expr("balance + ?", amount)).
		Where(squirrel.Eq{"number": to}).
		ToSql()

	result, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AccountStorage - Transfer - tx.Exec: %v", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s doesn't exist", to)
	}

	// update transaction note about successful transaction
	sql, args, _ = p.builder.
		Update(transactionsTable).
		Set("status", statusSuccess).
		Where(squirrel.Eq{"id": transactionId}).
		ToSql()

	err = setTransactionStatusSuccess(ctx, tx, sql, args...)
	if err != nil {
		return fmt.Errorf("AccountStorage - Transfer - setTransactionStatusSuccess: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("AccountStorage - Transfer - tx.Commit: %w", err)
	}

	success = true

	return nil
}

func (p *Postgres) GetBalanceByUserID(ctx context.Context, userID int) ([]models.UserBalance, error) {
	var result []models.UserBalance

	query := fmt.Sprintf(`
		SELECT
   		 	a.currency_code,
    		(a.balance - COALESCE(SUM(CASE WHEN t.status = 'Created' THEN t.amount ELSE 0 END), 0)) AS actual_balance,
    		COALESCE(SUM(CASE WHEN t.status = 'Created' THEN t.amount ELSE 0 END), 0) AS frozen_balance
		FROM
   		 	%s AS a
    			LEFT JOIN
    		%s AS t ON a.number = t.from_account
		WHERE user_id = $1
		GROUP BY
    		a.currency_code, a.balance;`, accountsTable, transactionsTable)

	rows, err := p.Pool.Query(ctx, query, userID)
	if err != nil {
		return result, fmt.Errorf("AccountStorage - GetBalanceByCurrency - p.Pool.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var byCurrency models.UserBalance

		err := rows.Scan(&byCurrency.CurrencyCode, &byCurrency.Actual, &byCurrency.Reserved)
		if err != nil {
			return result, fmt.Errorf("AccountStorage - GetBalanceByCurrency - rows.Scan: %w", err)
		}

		result = append(result, byCurrency)
	}

	return result, nil
}
