package models

import (
	"errors"
	"gochi/config"

	"gorm.io/gorm"
)

var Users []User

type User struct {
	gorm.Model
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email" gorm:"unique"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	Name    string `json:"name"`
	Email   string `json:"email"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

func VerifyUser(email string, password string) (*User, error) {
	var user User
	config.DB.Where("email = ?", email).First(&user)
		if user.Email != email || user.Password != password {
			return nil, errors.New("Invalid Credential!")
		}
		return &user, nil
}

func FindDuplicates(email string) bool {
	var user User
	config.DB.Where("email = ?", email).First(&user)
	if user.Email == email {
		return true
	}
	return false
}
