package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hdkef/be-assignment/services/account/domain/entity"
	"github.com/hdkef/be-assignment/services/account/domain/repository"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance for operation")
	ErrNoUnitOfWork        = errors.New("unit of work is required for this operation")
)

type AccountRepo struct {
	db *sql.DB
}

// DecrementBalance implements repository.AccountRepository.
func (a *AccountRepo) DecrementBalance(ctx context.Context, id uuid.UUID, decrAmount float64, uow *repository.UnitOfWork) error {
	if uow == nil || uow.Tx == nil {
		return ErrNoUnitOfWork
	}

	query := `
		UPDATE accounts.accounts
		SET balance = GREATEST(balance - $1, 0)
		WHERE id = $2
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

// IncrementBalance implements repository.AccountRepository.
func (a *AccountRepo) IncrementBalance(ctx context.Context, id uuid.UUID, incrAmount float64, uow *repository.UnitOfWork) error {
	if uow == nil || uow.Tx == nil {
		return ErrNoUnitOfWork
	}

	query := `
		UPDATE accounts.accounts
		SET balance = balance + $1
		WHERE id = $2
		RETURNING balance
	`

	var newBalance float64
	err := uow.Tx.QueryRowContext(ctx, query, incrAmount, id).Scan(&newBalance)
	if err != nil {
		return fmt.Errorf("failed to increment balance: %w", err)
	}

	return nil
}

// CreateAccount implements repository.AccountRepository.
func (a *AccountRepo) CreateAccount(ctx context.Context, acc *entity.Account, uow *repository.UnitOfWork) error {
	query := `INSERT INTO accounts.accounts (id, created_at, user_id, acc_type, acc_desc, currency, balance) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	if uow != nil && uow.Tx != nil {
		_, err := uow.Tx.ExecContext(ctx, query, acc.ID, time.Now(), acc.UserID, acc.AccType, acc.AccDesc, acc.Currency, acc.Balance)
		if err != nil {
			return fmt.Errorf("failed to execute insert query: %w", err)
		}
	} else {
		_, err := a.db.ExecContext(ctx, query, acc.ID, time.Now(), acc.UserID, acc.AccType, acc.AccDesc, acc.Currency, acc.Balance)
		if err != nil {
			return fmt.Errorf("failed to execute insert query: %w", err)
		}
	}
	return nil
}

func NewAccountRepo(db *sql.DB) repository.AccountRepository {
	return &AccountRepo{
		db: db,
	}
}
