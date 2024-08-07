package usecase

import (
	"context"

	"github.com/google/uuid"
	et2 "github.com/hdkef/be-assignment/pkg/domain/entity"
	"github.com/hdkef/be-assignment/services/account/domain/entity"
)

// TransactionCreatedWithdraw implements usecase.AccountUC.
func (a *AccountUC) TransactionCreatedWithdraw(ctx context.Context, dto et2.TransactionCreatedEventDto) error {

	accId, err := uuid.Parse(dto.Detail.AccID)
	if err != nil {
		return err
	}

	reffNum, err := uuid.Parse(dto.ReffNum)
	if err != nil {
		return err
	}

	// begin trx
	uow, err := a.UoW.NewUnitOfWork()
	if err != nil {
		uow.Tx.Rollback()
		return err
	}

	// update balance
	err = a.AccountRepo.DecrementBalance(ctx, accId, dto.Detail.Amount, uow)
	if err != nil {
		uow.Tx.Rollback()
		return err
	}

	hist := entity.History{
		Id:        uuid.New(),
		ReffNum:   reffNum,
		CreatedAt: dto.Detail.CreatedAt,
		AccID:     accId,
		TrxType:   entity.ENUM_HISTORY_TRX_TYPE_NEGATIVE,
		Amount:    dto.Detail.Amount,
		Status:    entity.ENUM_HISTORY_STATUS(dto.Detail.Status),
		Desc:      dto.Detail.Desc,
	}

	// insert history
	err = a.HistoryRepo.CreateHistory(ctx, &hist, uow)
	if err != nil {
		uow.Tx.Rollback()
		return err
	}

	err = uow.Tx.Commit()
	if err != nil {
		uow.Tx.Rollback()
		return err
	}

	return nil
}
