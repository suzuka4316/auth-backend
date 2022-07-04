package models

import (
	"errors"

	"github.com/badoux/checkmail"
)

type User struct {
	Id       uint   `json:"id"    gorm:"primaryKey; not null"`
	Name     string `json:"name"  gorm:"not null"`
	Email    string `json:"email" gorm:"unique; not null"`
	Password []byte `json:"-"     gorm:"not null"` // do not return password to client as response
}

func (u *User) Validate() error {
	if u.Password == nil {
		return errors.New("password required")
	}
	if u.Email == "" {
		return errors.New("email required")
	}
	if err := checkmail.ValidateFormat(u.Email); err != nil {
		return errors.New("invalid email")
	}
	return nil	
}