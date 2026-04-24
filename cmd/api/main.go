// @title           User Service API
// @version         1.0
// @description     API para gestión de usuarios
// @host            localhost:8081
// @BasePath        /

package main

import (
	"context"
	"log"
	"net/http"

	"github.com/SahidAyala/Event-Streaming-User-Service-Example/internal/infrastructure/persistence"
	"github.com/SahidAyala/Event-Streaming-User-Service-Example/internal/user"
	"github.com/go-chi/chi/v5"
)

func main() {
	ctx := context.Background()

	pool, err := persistence.NewPool(ctx, persistence.Config{
		DSN: "postgresql://postgres:root@localhost:5432/users", // TODO: mover a .env
	})
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	userModule := user.NewModule(pool)

	r := chi.NewRouter()
	r.Post("/users", userModule.Handler.CreateUser)
	r.Get("/users/{id}", userModule.Handler.GetUserById)
	r.Patch("/users/{id}/email", userModule.Handler.UpdateEmail)
	r.Patch("/users/{id}/password", userModule.Handler.UpdatePassword)

	log.Println("Server running on: 8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
