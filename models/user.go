package models

import (
	"errors"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/suzuka4316/auth-backend/utils"
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
			return errors.New(utils.PasswordRequired)
		}
		if u.Email == "" {
			return errors.New(utils.EmailRequired)
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New(utils.InvalidEmail)
		}
		return nil
	}

	if u.Name == "" {
		return errors.New(utils.NameRequired)
	}
	if u.Password == "" {
		return errors.New(utils.PasswordRequired)
	}
	if u.Email == "" {
		return errors.New(utils.EmailRequired)
	}
	if err := checkmail.ValidateFormat(u.Email); err != nil {
		return errors.New(utils.InvalidEmail)
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

	return &newUser, nil
}

func GetUserByEmail(db *gorm.DB, email interface{}) (*User, error) {
	var user User
	if err := db.Debug().Where("email = ?", email).First(&user).Error; err != nil {
		return &User{}, err
	}

	return &user, nil
}