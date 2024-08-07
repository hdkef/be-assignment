package service

import (
	"context"

	et2 "github.com/hdkef/be-assignment/pkg/domain/entity"
)

type Publisher interface {
	PublishCreateUserEvents(ctx context.Context, dto *et2.UserCreatedEventDto) error
	PublishCreateAccountEvents(ctx context.Context, dto *et2.AccountCreatedEventDto) error
}
