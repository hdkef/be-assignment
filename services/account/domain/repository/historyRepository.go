package repository

import (
	"context"

	"github.com/hdkef/be-assignment/services/account/domain/entity"
)

type HistoryRepository interface {
	CreateHistory(ctx context.Context, history *entity.History, uow *UnitOfWork) error
}
