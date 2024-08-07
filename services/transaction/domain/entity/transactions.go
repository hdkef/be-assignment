package entity

import (
	"errors"

	"github.com/google/uuid"
)

type SendTransactionDto struct {
	AccID   uuid.UUID
	Amount  float64
	Desc    string
	ToAccID uuid.UUID
}

type WithdrawTransactionDto struct {
	AccID  uuid.UUID
	Amount float64
	Desc   string
}

func (s *SendTransactionDto) Validate() error {
	if s.AccID == uuid.Nil {
		return errors.New("accID is required")
	}
	if s.ToAccID == uuid.Nil {
		return errors.New("toAccID is required")
	}

	if s.AccID == s.ToAccID {
		return errors.New("cannot send to same account")
	}

	if s.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if len(s.Desc) == 0 {
		return errors.New("description is required")
	}
	return nil
}

func (w *WithdrawTransactionDto) Validate() error {
	if w.AccID == uuid.Nil {
		return errors.New("accID is required")
	}
	if w.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if len(w.Desc) == 0 {
		return errors.New("description is required")
	}
	return nil
}
