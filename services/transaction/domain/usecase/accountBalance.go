package usecase

import (
	"context"

	"github.com/hdkef/be-assignment/pkg/domain/entity"
)

type AccountBalanceUC interface {
	CreateAccountBalanceUserEvent(ctx context.Context, dto *entity.UserCreatedEventDto) error
	CreateAccountBalanceAccountEvent(ctx context.Context, dto *entity.AccountCreatedEventDto) error
}
