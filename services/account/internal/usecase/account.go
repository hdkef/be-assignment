package usecase

import (
	repository2 "github.com/hdkef/be-assignment/services/account/domain/repository"
	"github.com/hdkef/be-assignment/services/account/domain/usecase"
	"github.com/hdkef/be-assignment/services/account/internal/repository"
)

type AccountUC struct {
	UoW         repository.UnitOfWorkImplementor
	HistoryRepo repository2.HistoryRepository
	AccountRepo repository2.AccountRepository
}

func NewAccountUC() usecase.AccountUC {
	return &AccountUC{}
}
