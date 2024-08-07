package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hdkef/be-assignment/services/transaction/domain/entity"
	"github.com/hdkef/be-assignment/services/transaction/domain/repository"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance for operation")
	ErrNoUnitOfWork        = errors.New("unit of work is required for this operation")
)

type AccountBalanceRepo struct {
	db *sql.DB
}

// DecrementBalance implements repository.AccountBalanceRepository.
func (a *AccountBalanceRepo) DecrementBalance(ctx context.Context, id uuid.UUID, decrAmount float64, uow *repository.UnitOfWork) error {
	if uow == nil || uow.Tx == nil {
		return ErrNoUnitOfWork
	}

	query := `
		UPDATE transactions.accounts_balance
		SET balance = GREATEST(balance - $1, 0)
		WHERE acc_id = $2
		RETURNING balance
	`

	var newBalance float64
	err := uow.Tx.QueryRowContext(ctx, query, decrAmount, id).Scan(&newBalance)
	if err != nil {
		return fmt.Errorf("failed to decrement balance: %w", err)
	}

	if newBalance == 0 && decrAmount > 0 {
		return ErrInsufficientBalance
	}

	return nil
}

// IncrementBalance implements repository.AccountBalanceRepository.
func (a *AccountBalanceRepo) IncrementBalance(ctx context.Context, id uuid.UUID, incrAmount float64, uow *repository.UnitOfWork) error {
	if uow == nil || uow.Tx == nil {
		return ErrNoUnitOfWork
	}

	query := `
		UPDATE transactions.accounts_balance
		SET balance = balance + $1
		WHERE acc_id = $2
		RETURNING balance
	`

	var newBalance float64
	err := uow.Tx.QueryRowContext(ctx, query, incrAmount, id).Scan(&newBalance)
	if err != nil {
		return fmt.Errorf("failed to increment balance: %w", err)
	}

	return nil
}

// CreateAccountBalance implements repository.AccountBalanceRepository.
func (a *AccountBalanceRepo) CreateAccountBalance(ctx context.Context, acc *entity.AccountBalance, uow *repository.UnitOfWork) error {
	query := `
		INSERT INTO transactions.accounts_balance (user_id, acc_id, currency, balance)
		VALUES ($1, $2, $3, $4)
	`

	var err error
	if uow != nil && uow.Tx != nil {
		_, err = uow.Tx.ExecContext(ctx, query, acc.UserID, acc.AccID, acc.Currency, acc.Balance)
	} else {
		_, err = a.db.ExecContext(ctx, query, acc.UserID, acc.AccID, acc.Currency, acc.Balance)
	}

	if err != nil {
		return fmt.Errorf("failed to create account balance: %w", err)
	}

	return nil
}

func NewAccountBalanceRepo(db *sql.DB) repository.AccountBalanceRepository {
	return &AccountBalanceRepo{
		db: db,
	}
}
