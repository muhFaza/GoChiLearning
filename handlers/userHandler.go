package handlers

import (
	"encoding/json"
	"gochi/helper"
	"gochi/models"
	"io/ioutil"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var users []models.User

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userRegister models.User

	if err := json.NewDecoder(r.Body).Decode(&userRegister); err != nil {
		http.Error(w, "Invalid Request!", 400)
	}

	validate := validator.New()
	if err := validate.Struct(userRegister); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	users = append(users, userRegister)

	json.NewEncoder(w).Encode(models.RegisterResponse{Message: "User Registered!", Name: userRegister.Name, Email: userRegister.Email})
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userLogin models.UserLogin

	if err := json.NewDecoder(r.Body).Decode(&userLogin); err != nil {
		http.Error(w, "Invalid Request!", 400)
	}

	for _, user := range users {
		if user.Email == userLogin.Email && user.Password == userLogin.Password {
			if token, err := helper.GenerateToken(userLogin.Email); err != nil {
				http.Error(w, "Something went wrong!", 400)
				return
			} else {
				json.NewEncoder(w).Encode(models.LoginResponse{Message: "Login Success!", Token: token})
			}
			return
		}
	}

	http.Error(w, "Invalid Credentials!", 400)
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response, err := http.Get("https://jsonplaceholder.typicode.com/todos")

	if err != nil {
		http.Error(w, "Something went wrong!", 400)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Something went wrong!", 400)
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
			http.Error(w, "Unauthorized!", 401)
			return
		}

		claims, err := helper.ValidateToken(token)
		if err != nil {
			http.Error(w, "Unauthorized!", 401)
			return
		}

		r.Header.Set("email", claims.Email)
		next.ServeHTTP(w, r)
	})
}
