package user

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/SahidAyala/Event-Streaming-User-Service-Example/internal/user/application"
	userhttp "github.com/SahidAyala/Event-Streaming-User-Service-Example/internal/user/infrastructure/http"
	"github.com/SahidAyala/Event-Streaming-User-Service-Example/internal/user/infrastructure/persistence"
)

type Module struct {
	Handler *userhttp.Handler
	Service *application.Service
}

func NewModule(pool *pgxpool.Pool) *Module {
	repo := persistence.NewPostgresRepository(pool)
	service := application.NewService(repo)
	handler := userhttp.NewHandler(service)

	return &Module{
		Handler: handler,
		Service: service,
	}
}
