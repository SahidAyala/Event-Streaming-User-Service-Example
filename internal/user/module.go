package user

import (
	"github.com/jackc/pgx/v5/pgxpool"

	eventsinfra "github.com/SahidAyala/Event-Streaming-User-Service-Example/internal/shared/events/infrastructure"
	"github.com/SahidAyala/Event-Streaming-User-Service-Example/internal/user/application"
	userhttp "github.com/SahidAyala/Event-Streaming-User-Service-Example/internal/user/infrastructure/http"
	"github.com/SahidAyala/Event-Streaming-User-Service-Example/internal/user/infrastructure/persistence"
)

type Module struct {
	Handler *userhttp.Handler
	Service *application.Service
}

func NewModule(pool *pgxpool.Pool, eventsBaseURL, eventsAPIKey string) *Module {
	repo := persistence.NewPostgresRepository(pool)
	service := application.NewService(repo)
	service.SetPublisher(eventsinfra.NewHTTPEventPublisher(eventsBaseURL, eventsAPIKey))
	handler := userhttp.NewHandler(service)

	return &Module{
		Handler: handler,
		Service: service,
	}
}
