// @title           User Service API
// @version         1.0
// @description     API para gestión de usuarios
// @host            localhost:8082
// @BasePath        /

package main

import (
	"context"
	"log"
	"net/http"

	_ "github.com/SahidAyala/Event-Streaming-User-Service-Example/docs"
	appconfig "github.com/SahidAyala/Event-Streaming-User-Service-Example/internal/infrastructure/config"
	"github.com/SahidAyala/Event-Streaming-User-Service-Example/internal/infrastructure/persistence"
	"github.com/SahidAyala/Event-Streaming-User-Service-Example/internal/user"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	ctx := context.Background()
	cfg := appconfig.Load()

	pool, err := persistence.NewPool(ctx, persistence.Config{
		DSN: cfg.PostgresDSN,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	userModule := user.NewModule(pool, cfg.EventsBaseURL, cfg.EventsAPIKey)

	r := chi.NewRouter()
	r.Post("/users", userModule.Handler.CreateUser)
	r.Get("/users/{id}", userModule.Handler.GetUserById)
	r.Patch("/users/{id}/email", userModule.Handler.UpdateEmail)
	r.Patch("/users/{id}/password", userModule.Handler.UpdatePassword)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Println("Server running on: 8082")
	log.Fatal(http.ListenAndServe(":8082", r))
}
