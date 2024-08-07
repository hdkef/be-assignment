package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hdkef/be-assignment/services/account/domain/entity"
	"github.com/hdkef/be-assignment/services/account/domain/repository"
)

type UserRepo struct {
	db *sql.DB
}

// CreateUser implements repository.UserRepository.
func (u *UserRepo) CreateUser(ctx context.Context, user *entity.User, uow *repository.UnitOfWork) error {
	query := `INSERT INTO accounts.users (id, name, dob, job, created_at, email) VALUES ($1, $2, $3, $4, $5, $6)`

	if uow != nil && uow.Tx != nil {
		_, err := uow.Tx.ExecContext(ctx, query, user.ID, user.Name, user.DateOfBirth, user.Job, user.CreatedAt, user.Email)
		if err != nil {
			return fmt.Errorf("failed to execute insert query: %w", err)
		}
	} else {
		_, err := u.db.ExecContext(ctx, query, user.ID, user.Name, user.DateOfBirth, user.Job, user.CreatedAt, user.Email)
		if err != nil {
			return fmt.Errorf("failed to execute insert query: %w", err)
		}
	}
	return nil
}

func NewUserRepo(db *sql.DB) repository.UserRepository {
	return &UserRepo{
		db: db,
	}
}
