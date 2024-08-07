package usecase

import (
	"context"

	"github.com/hdkef/be-assignment/services/account/domain/entity"
)

type UserUsecase interface {
	SignUp(ctx context.Context, dto *entity.UserSignUpDto) error
}
