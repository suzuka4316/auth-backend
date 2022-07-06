package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/badoux/checkmail"
	"gorm.io/gorm"
)

type User struct {
	Id       uint32 `json:"id"       gorm:"primaryKey;auto_increment"`
	Name     string `json:"name"     gorm:"size:255;not null"`
	Email    string `json:"email"    gorm:"size:100;unique;not null"`
	Password string `json:"password" gorm:"not null"`
}

func (u *User) Validate(action string) error {
	if strings.ToLower(action) == "login" {
		if u.Password == "" {
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

	if u.Name == "" {
		return errors.New("name required")
	}
	if u.Password == "" {
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

func (u *User) SaveUser(db *gorm.DB, hashedPassword []byte) (*User, error) {
	var err error

	newUser := User {
		Name:     u.Name,
		Email:    u.Email,
		Password: string(hashedPassword),
	}
	if err = db.Debug().Create(&newUser).Error; err != nil {
		return &User{}, err
	}

	return u, nil
}

func GetUserByEmail(db *gorm.DB, email interface{}) (*User, error) {
	var user User
	if err := db.Debug().Where("email = ?", email).First(&user).Error; err != nil {
		fmt.Printf("GetUserByEmail err %v", err)
		return &User{}, err
	}

	fmt.Printf("GetUserByEmail user %v", user)
	return &user, nil
}