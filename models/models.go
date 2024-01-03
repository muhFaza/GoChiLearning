package models

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
	Name		string `json:"name"`
	Email		string `json:"email"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Token	 string `json:"token"`
}