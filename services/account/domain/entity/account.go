package entity

import (
	"errors"

	"github.com/google/uuid"
)

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

type CreateAccountDto struct {
	userID   uuid.UUID
	AccType  ENUM_ACCOUNT_ACC_TYPE `json:"accType"`
	AccDesc  string                `json:"accDesc"`
	Currency string                `json:"currency"`
}

func (dto *CreateAccountDto) SetUserID(userID uuid.UUID) {
	dto.userID = userID
}

func (dto *CreateAccountDto) GetUserID() uuid.UUID {
	return dto.userID
}

func (dto *CreateAccountDto) Validate() error {
	if dto.userID == uuid.Nil {
		return errors.New("userID is required")
	}

	if dto.AccType != ENUM_ACCOUNT_ACC_TYPE_DEBIT && dto.AccType != ENUM_ACCOUNT_ACC_TYPE_CREDIT && dto.AccType != ENUM_ACCOUNT_ACC_TYPE_LOAN {
		return errors.New("accType is invalid")
	}

	if dto.AccDesc == "" {
		return errors.New("accDesc is required")
	}

	if len(dto.Currency) != 3 {
		return errors.New("currency is invalid")
	}

	return nil
}
