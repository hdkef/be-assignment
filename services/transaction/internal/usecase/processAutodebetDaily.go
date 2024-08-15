package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hdkef/be-assignment/pkg/logger"
	"github.com/hdkef/be-assignment/services/transaction/domain/entity"
)

// ProcessAutodebet implements usecase.TransactionUsecase.
func (t *TransactionUC) ProcessAutodebetDaily(ctx context.Context) error {

	// begin trx
	uow, err := t.UoW.NewUnitOfWork()
	if err != nil {
		logger.LogError(ctx, err)
		uow.Tx.Rollback()
		return err
	}

	// get schedule with status has not been checked and last checked < now and schedule daily
	scheds, err := t.ScheduleRepo.GetUnprocessed(ctx, time.Now(), uow)
	if err != nil {
		logger.LogError(ctx, err)
		uow.Tx.Rollback()
		return err
	}

	// insert into queue trx
	for _, sched := range scheds {
		newQ := entity.QueuedTrx{
			ID:            uuid.New(),
			CreatedAt:     time.Now(),
			ScheduleTrxID: sched.ID,
			Status:        entity.ENUM_QUEUED_TRX_STATUS_PENDING,
			Result:        "",
		}
		err = t.QueueRepo.Create(ctx, &newQ, uow)
		if err != nil {
			logger.LogError(ctx, err)
			uow.Tx.Rollback()
			return err
		}
	}

	// update has checked false & last checked now

	ids := make([]uuid.UUID, len(scheds))
	for i, sched := range scheds {
		ids[i] = sched.ID
	}

	err = t.ScheduleRepo.UpdateChecked(ctx, ids, time.Now(), uow)
	if err != nil {
		logger.LogError(ctx, err)
		uow.Tx.Rollback()
		return err
	}

	err = uow.Tx.Commit()
	if err != nil {
		logger.LogError(ctx, err)
		uow.Tx.Rollback()
		return err
	}

	return nil
}
