package models

import (
	"errors"
)

var Users []User

type User struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
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
	for _, user := range Users {
		if user.Email == email && user.Password == password {
			return &user, nil
		}
	}
	return nil, errors.New("Invalid Credentials!")
}

func FindDuplicates (email string) bool {
	for _, user := range Users {
		if email == user.Email{
			return true
		}
	}
	return false
}