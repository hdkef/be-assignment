package repository

import (
	"context"

	"github.com/hdkef/be-assignment/services/account/domain/entity"
)

type HistoryRepository interface {
	CreateHistory(ctx context.Context, history *entity.History, uow *UnitOfWork) error
	GetHistory(ctx context.Context, f *entity.GetHistoryFilter, opt *entity.GetHistoryOptions) ([]*entity.History, error)
}
