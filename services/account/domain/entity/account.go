package entity

import "github.com/google/uuid"

type ENUM_ACCOUNT_ACC_TYPE string

const (
	ENUM_ACCOUNT_ACC_TYPE_DEBIT  ENUM_ACCOUNT_ACC_TYPE = "DEBIT"
	ENUM_ACCOUNT_ACC_TYPE_CREDIT ENUM_ACCOUNT_ACC_TYPE = "CREDIT"
	ENUM_ACCOUNT_ACC_TYPE_LOAN   ENUM_ACCOUNT_ACC_TYPE = "LOAN"
)

type Account struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	AccType  ENUM_ACCOUNT_ACC_TYPE
	AccDesc  string
	Currency string
	Balance  float64
}
