package repository

import (
	"database/sql"

	"github.com/hdkef/be-assignment/services/transaction/domain/repository"
)

type UnitOfWorkImplementor struct {
	Db *sql.DB
}

func (u *UnitOfWorkImplementor) NewUnitOfWork() (*repository.UnitOfWork, error) {

	tx, err := u.Db.Begin()
	if err != nil {
		return nil, err
	}

	return &repository.UnitOfWork{
		Tx: tx,
	}, nil
}
