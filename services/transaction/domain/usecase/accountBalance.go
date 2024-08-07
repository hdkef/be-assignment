package usecase

import (
	"context"

	"github.com/hdkef/be-assignment/pkg/domain/entity"
)

type AccountBalanceUC interface {
	CreateAccountBalance(ctx context.Context, dto *entity.UserCreatedEventDto) error
}
