package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hdkef/be-assignment/services/account/domain/entity"
	"github.com/hdkef/be-assignment/services/account/domain/repository"
)

type UserAddressRepo struct {
	db *sql.DB
}

// CreateAddress implements repository.UserAddressRepository.
func (u *UserAddressRepo) CreateAddress(ctx context.Context, address *entity.UserAddress, uow *repository.UnitOfWork) error {
	query := `INSERT INTO accounts.user_address (id, user_id, zip, address, district, city, province, country) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	if uow != nil && uow.Tx != nil {
		_, err := uow.Tx.ExecContext(ctx, query, address.ID, address.UserID, address.ZIP, address.Address, address.District, address.City, address.Province, address.Country)
		if err != nil {
			return fmt.Errorf("failed to execute insert query: %w", err)
		}
	} else {
		_, err := u.db.ExecContext(ctx, query, address.ID, address.UserID, address.ZIP, address.Address, address.District, address.City, address.Province, address.Country)
		if err != nil {
			return fmt.Errorf("failed to execute insert query: %w", err)
		}
	}
	return nil
}

func NewUserAddressRepo(db *sql.DB) repository.UserAddressRepository {
	return &UserAddressRepo{
		db: db,
	}
}
