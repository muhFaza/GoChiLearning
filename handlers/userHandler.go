package handlers

import (
	"encoding/json"
	"fmt"
	"gochi/config"
	"gochi/helper"
	"gochi/models"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userRegister models.User

	if err := json.NewDecoder(r.Body).Decode(&userRegister); err != nil {
		http.Error(w, "Invalid Request!", http.StatusBadRequest)
	}

	validate := validator.New()
	if err := validate.Struct(userRegister); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if models.FindDuplicates(userRegister.Email) {
		http.Error(w, "Email has been registered!", http.StatusConflict)
		return
	}

	config.DB.Create(&userRegister)

	json.NewEncoder(w).Encode(models.RegisterResponse{Message: "User Registered!", Name: userRegister.Name, Email: userRegister.Email})
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []models.User
	config.DB.Find(&users)
	json.NewEncoder(w).Encode(users)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userLogin models.UserLogin

	if err := json.NewDecoder(r.Body).Decode(&userLogin); err != nil {
		http.Error(w, "Invalid Request!", http.StatusBadRequest)
		return
	}

	_, err := models.VerifyUser(userLogin.Email, userLogin.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if token, err := helper.GenerateToken(userLogin.Email); err != nil {
		http.Error(w, "Something went wrong!", http.StatusBadRequest)
		return
	} else {
		json.NewEncoder(w).Encode(models.LoginResponse{Message: "Login Success!", Token: token})
		return
	}
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response, err := http.Get("https://jsonplaceholder.typicode.com/todos")

	if err != nil {
		http.Error(w, "Something went wrong!", http.StatusBadRequest)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Something went wrong!", http.StatusBadRequest)
	}

	type todos []struct {
		UserID    int    `json:"userId"`
		ID        int    `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}

	var todo todos

	if err := json.Unmarshal(responseData, &todo); err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(todo)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token == "" {
			http.Error(w, "Unauthorized!", http.StatusUnauthorized)
			return
		}

		claims, err := helper.ValidateToken(token)
		if err != nil {
			http.Error(w, "Unauthorized!", http.StatusUnauthorized)
			return
		}

		r.Header.Set("email", claims.Email)
		next.ServeHTTP(w, r)
	})
}

func UpdateName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type updateName struct {
		Name string `json:"name" validate:"required"`
	}
	var updateReq updateName
	email := r.Header.Get("email")

	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, "Invalid Request!", http.StatusBadRequest)
		return
	}

	if updateReq.Name == "" {
		http.Error(w, "Invalid Request!", http.StatusBadRequest)
		return
	}

	var user models.User
	fmt.Println(updateReq)
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	user.Name = updateReq.Name
	if err := config.DB.Save(&user).Error; err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(models.RegisterResponse{Message: "Name successfully updated!", Name: user.Name, Email: user.Email})
}
