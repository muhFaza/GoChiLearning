package main

import (
	"fmt"
	"gochi/config"
	"gochi/handlers"
	"gochi/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error loading .env file")
	}

	config.DB = func() *gorm.DB {
		dbConfig := config.GetDBConfig()
		dsn := dbConfig.GetDBURL()
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("Failed to connect to database!")
		}
		err = db.AutoMigrate(&models.User{})
		if err != nil {
			panic("Failed to migrate database!")
		}
		fmt.Println("Database connected!")
		return db
	}()

}

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Test!"))
	})

	r.Group(func(r chi.Router) {
		r.Post("/register", handlers.RegisterUser)

		r.Get("/users", handlers.GetUsers)

		// Login with email and password
		// Returns JWT token
		r.Post("/login", handlers.LoginUser)
	})

	r.Group(func(r chi.Router) {
		r.Use(handlers.AuthMiddleware)

		r.Get("/todos", handlers.GetTodos)
		r.Patch("/profile/name", handlers.UpdateName)
	})

	http.ListenAndServe(":8080", r)
}
