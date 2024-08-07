package repository

import "database/sql"

type UnitOfWork struct {
	Tx *sql.Tx
}
