package entity

import "github.com/google/uuid"

type UserAddress struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	ZIP      uint32
	Address  string
	District string
	City     string
	Province string
	Country  string
}
