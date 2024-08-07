package entity

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hdkef/be-assignment/pkg/helper"
)

type User struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	Name        string
	Email       string
	DateOfBirth time.Time
	Job         string
}

type UserSignUpDto struct {
	ID                   uuid.UUID
	Name                 string
	Email                string
	DateOfBirth          time.Time
	Job                  string
	Address              string
	District             string
	City                 string
	Province             string
	Country              string
	ZIP                  uint32
	FirstAccountDesc     string
	FirstAccountCurrency string
}

func (dto *UserSignUpDto) Validate() error {

	if dto.Name == "" {
		return errors.New("name is required")
	}

	if dto.Email == "" || !helper.IsValidEmail(dto.Email) {
		return errors.New("invalid email address")
	}

	if dto.DateOfBirth.IsZero() || dto.DateOfBirth.After(time.Now()) {
		return errors.New("invalid date of birth")
	}

	if dto.Job == "" {
		return errors.New("job is required")
	}

	if dto.Address == "" {
		return errors.New("address is required")
	}

	if dto.District == "" {
		return errors.New("district is required")
	}

	if dto.City == "" {
		return errors.New("city is required")
	}

	if dto.Province == "" {
		return errors.New("province is required")
	}

	if dto.Country == "" {
		return errors.New("country is required")
	}

	if dto.FirstAccountDesc == "" {
		return errors.New("first account description is required")
	}

	if dto.FirstAccountCurrency == "" {
		return errors.New("first account currency is required")
	}

	if dto.ZIP == 0 {
		return errors.New("zip is required")
	}

	return nil
}

func (dto *UserSignUpDto) SetUserID(userId string) error {

	uuid, err := uuid.Parse(userId)
	if err != nil {
		return err
	}
	dto.ID = uuid
	return nil
}

func (dto *UserSignUpDto) SetDOB(dob string) error {
	// Parse the date string
	parsedDOB, err := time.Parse("2006-01-02", dob)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}

	// Set the DateOfBirth field
	dto.DateOfBirth = parsedDOB
	return nil
}
