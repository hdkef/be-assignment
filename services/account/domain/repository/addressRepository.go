package repository

import (
	"context"

	"github.com/hdkef/be-assignment/services/account/domain/entity"
)

type UserAddressRepository interface {
	CreateAddress(ctx context.Context, address *entity.UserAddress, uow *UnitOfWork) error
}
