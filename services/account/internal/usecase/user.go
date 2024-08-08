package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	entity2 "github.com/hdkef/be-assignment/pkg/domain/entity"
	"github.com/hdkef/be-assignment/pkg/logger"
	"github.com/hdkef/be-assignment/services/account/domain/entity"
	repository2 "github.com/hdkef/be-assignment/services/account/domain/repository"
	"github.com/hdkef/be-assignment/services/account/domain/service"
	"github.com/hdkef/be-assignment/services/account/internal/repository"
)

// create address and accounts

type UserUC struct {
	UoW             repository.UnitOfWorkImplementor
	UserRepo        repository2.UserRepository
	UserAddressRepo repository2.UserAddressRepository
	AccountRepo     repository2.AccountRepository
	Publisher       service.Publisher
}

// SignUp implements usecase.UserUsecase.
func (u *UserUC) SignUp(ctx context.Context, dto *entity.UserSignUpDto) error {

	err := dto.Validate()
	if err != nil {
		logger.LogError(ctx, err)
		return err
	}

	uow, err := u.UoW.NewUnitOfWork()
	if err != nil {
		logger.LogError(ctx, err)
		return err
	}

	// create user
	user := entity.User{
		ID:          dto.ID,
		CreatedAt:   time.Now(),
		Name:        dto.Name,
		Email:       dto.Email,
		DateOfBirth: dto.DateOfBirth,
		Job:         dto.Job,
	}

	err = u.UserRepo.CreateUser(ctx, &user, uow)
	if err != nil {
		logger.LogError(ctx, err)
		uow.Tx.Rollback()
		return err
	}

	// create address
	userAddress := entity.UserAddress{
		ID:       uuid.New(),
		UserID:   user.ID,
		Address:  dto.Address,
		City:     dto.City,
		Country:  dto.Country,
		ZIP:      dto.ZIP,
		Province: dto.Province,
		District: dto.District,
	}
	err = u.UserAddressRepo.CreateAddress(ctx, &userAddress, uow)
	if err != nil {
		logger.LogError(ctx, err)
		uow.Tx.Rollback()
		return err
	}

	// create new account if no field from dto then generated / set default
	acc := entity.Account{
		ID:       uuid.New(),
		UserID:   user.ID,
		AccType:  entity.ENUM_ACCOUNT_ACC_TYPE_DEBIT, // default first account debit
		AccDesc:  dto.FirstAccountDesc,
		Currency: dto.FirstAccountCurrency,
		Balance:  0.0,
	}
	err = u.AccountRepo.CreateAccount(ctx, &acc, uow)
	if err != nil {
		logger.LogError(ctx, err)
		uow.Tx.Rollback()
		return err
	}

	// commit trx
	err = uow.Tx.Commit()
	if err != nil {
		logger.LogError(ctx, err)
		uow.Tx.Rollback()
		return err
	}

	// publish user created event
	userCreatedEventDto := entity2.UserCreatedEventDto{
		UserID:          user.ID.String(),
		AccountID:       acc.ID.String(),
		AccountCurrency: acc.Currency,
	}
	err = u.Publisher.PublishCreateUserEvents(ctx, &userCreatedEventDto)
	if err != nil {
		logger.LogError(ctx, err)
		uow.Tx.Rollback()
		return err
	}

	return nil
}
