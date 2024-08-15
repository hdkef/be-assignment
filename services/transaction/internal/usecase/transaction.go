package usecase

import (
	"github.com/hdkef/be-assignment/services/transaction/domain/repository"
	"github.com/hdkef/be-assignment/services/transaction/domain/service"
	"github.com/hdkef/be-assignment/services/transaction/domain/usecase"
	repo2 "github.com/hdkef/be-assignment/services/transaction/internal/repository"
)

type TransactionUC struct {
	UoW            repo2.UnitOfWorkImplementor
	AccBalanceRepo repository.AccountBalanceRepository
	TrxLogsRepo    repository.TransactionLogsRepository
	ScheduleRepo   repository.ScheduleRepository
	QueueRepo      repository.QueueRepository
	Publisher      service.Publisher
}

func NewTransactionUC() usecase.TransactionUsecase {
	return &TransactionUC{}
}
