package usecase

import (
	"context"

	"github.com/google/uuid"
	entity2 "github.com/hdkef/be-assignment/pkg/domain/entity"
	"github.com/hdkef/be-assignment/services/transaction/domain/entity"
	repository2 "github.com/hdkef/be-assignment/services/transaction/domain/repository"
	"github.com/hdkef/be-assignment/services/transaction/internal/repository"
)

type AccountBalanceUC struct {
	UoW            repository.UnitOfWorkImplementor
	AccBalanceRepo repository2.AccountBalanceRepository
}

// CreateAccountBalance implements usecase.AccountBalanceUC.
func (a *AccountBalanceUC) CreateAccountBalance(ctx context.Context, dto *entity2.UserCreatedEventDto) error {

	userId, err := uuid.Parse(dto.UserID)
	if err != nil {
		return err
	}

	accId, err := uuid.Parse(dto.AccountID)
	if err != nil {
		return err
	}

	accBalance := entity.AccountBalance{
		UserID:   userId,
		Balance:  0,
		AccID:    accId,
		Currency: dto.AccountCurrency,
	}

	uow, err := a.UoW.NewUnitOfWork()
	if err != nil {
		uow.Tx.Rollback()
		return err
	}

	err = a.AccBalanceRepo.CreateAccountBalance(ctx, &accBalance, uow)
	if err != nil {
		uow.Tx.Rollback()
		return err
	}

	// commit trx
	err = uow.Tx.Commit()
	if err != nil {
		uow.Tx.Rollback()
		return err
	}

	return nil
}
