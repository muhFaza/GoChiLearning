package main

import (
	"gochi/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

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
	})
	
	http.ListenAndServe(":8080", r)
}
