package entity

import "github.com/google/uuid"

type AccountBalance struct {
	UserID   uuid.UUID
	AccID    uuid.UUID
	Balance  float64
	Currency string
}
