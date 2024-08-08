package usecase

import (
	"context"

	"github.com/google/uuid"
	et2 "github.com/hdkef/be-assignment/pkg/domain/entity"
	"github.com/hdkef/be-assignment/pkg/logger"
	"github.com/hdkef/be-assignment/services/account/domain/entity"
	repository2 "github.com/hdkef/be-assignment/services/account/domain/repository"
	"github.com/hdkef/be-assignment/services/account/domain/service"
	"github.com/hdkef/be-assignment/services/account/domain/usecase"
	"github.com/hdkef/be-assignment/services/account/internal/repository"
)

type AccountUC struct {
	UoW         repository.UnitOfWorkImplementor
	HistoryRepo repository2.HistoryRepository
	AccountRepo repository2.AccountRepository
	Publisher   service.Publisher
}

// CreateAccount implements usecase.AccountUC.
func (a *AccountUC) CreateAccount(ctx context.Context, dto *entity.CreateAccountDto) error {

	err := dto.Validate()
	if err != nil {
		logger.LogError(ctx, err)
		return err
	}

	uow, err := a.UoW.NewUnitOfWork()
	if err != nil {
		logger.LogError(ctx, err)
		return err
	}

	// create account
	acc := entity.Account{
		ID:       uuid.New(),
		UserID:   dto.GetUserID(),
		AccType:  entity.ENUM_ACCOUNT_ACC_TYPE(dto.AccType),
		AccDesc:  dto.AccDesc,
		Currency: dto.Currency,
		Balance:  0.0,
	}
	err = a.AccountRepo.CreateAccount(ctx, &acc, uow)
	if err != nil {
		logger.LogError(ctx, err)
		uow.Tx.Rollback()
		return err
	}

	// publish event account created
	err = a.Publisher.PublishCreateAccountEvents(ctx, &et2.AccountCreatedEventDto{
		UserID:          acc.UserID.String(),
		AccountID:       acc.ID.String(),
		AccountCurrency: acc.Currency,
	})

	if err != nil {
		logger.LogError(ctx, err)
		uow.Tx.Rollback()
		return err
	}

	// commit trx
	err = uow.Tx.Commit()
	if err != nil {
		logger.LogError(ctx, err)
		uow.Tx.Rollback()
		return err
	}

	return nil
}

// GetHistory implements usecase.AccountUC.
func (a *AccountUC) GetHistory(ctx context.Context, dto *entity.GetHistoryDto) ([]*entity.History, error) {

	// query

	history, err := a.HistoryRepo.GetHistory(ctx, &entity.GetHistoryFilter{
		AccID:  dto.AccID,
		UserID: dto.UserID,
	}, &entity.GetHistoryOptions{
		Page:  dto.Page,
		Limit: dto.GetLimit(),
	})
	if err != nil {
		logger.LogError(ctx, err)
		return nil, err
	}

	return history, nil
}

func NewAccountUC() usecase.AccountUC {
	return &AccountUC{}
}
