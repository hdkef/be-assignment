package repository

import (
	"context"

	"github.com/hdkef/be-assignment/services/account/domain/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User, uow *UnitOfWork) error
}
