# User Service

Microservice for user management. Exposes a REST API to create and manage users, and publishes domain events to an external event streaming service.

## Architecture

The project follows hexagonal architecture with explicit layer separation inside each module:

```
internal/
  user/
    domain/               → User entity, domain errors
    application/
      ports/              → driven port interfaces
      service.go          → use cases
    infrastructure/
      http/               → HTTP handler (driving adapter)
      persistence/        → Postgres repository (driven adapter)
    module.go             → dependency wiring
  shared/
    events/               → Publisher port + Event struct
      infrastructure/     → HTTP adapter toward the events service
  infrastructure/
    config/               → environment variable loading
    persistence/          → Postgres connection pool
cmd/
  api/
    main.go               → entry point
migrations/               → SQL up/down
```

## Requirements

- Go 1.21+
- PostgreSQL
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI
- [swag](https://github.com/swaggo/swag) CLI (installed automatically by `make setup`)

## Getting started

```bash
make setup
```

This creates `.env` from `.env.example`, runs `go mod tidy`, and installs `swag`.

Edit `.env` with your values:

```env
POSTGRES_DSN=postgresql://user:password@localhost:5432/users?sslmode=disable
EVENTS_BASE_URL=http://localhost:8082
EVENTS_API_KEY=your-api-key
```

Apply migrations:

```bash
make migrate-up
```

Start the server:

```bash
make run
```

The API will be available at `http://localhost:8082` and the interactive docs at `http://localhost:8082/swagger/index.html`.

## Available commands

| Command | Description |
|---------|-------------|
| `make setup` | Bootstrap the environment (`.env`, dependencies, swag) |
| `make run` | Start the server |
| `make swagger` | Regenerate Swagger docs from code annotations |
| `make migrate-up` | Apply pending migrations |
| `make migrate-down` | Revert the last migration |

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/users` | Create user |
| `GET` | `/users/{id}` | Get user by ID |
| `PATCH` | `/users/{id}/email` | Update email |
| `PATCH` | `/users/{id}/password` | Update password |
| `GET` | `/swagger/*` | Interactive API docs |

## Events service integration

When a user is created, the service publishes a domain event to an external streaming service. The integration is decoupled through the `events.Publisher` port defined in `internal/shared/events/publisher.go`:

```go
type Publisher interface {
    Publish(ctx context.Context, event Event) error
}
```

The current adapter (`HTTPEventPublisher`) sends the event via `POST {EVENTS_BASE_URL}/events` with API key authentication in the `X-API-Key` header.

### Event published on user creation

```json
{
  "stream_id": "user:<uuid>",
  "type":      "user.created",
  "source":    "user-service",
  "payload": {
    "user_id": "<uuid>"
  },
  "metadata": {
    "email": "user@example.com"
  }
}
```

### Required environment variables

| Variable | Description | Example |
|----------|-------------|---------|
| `EVENTS_BASE_URL` | Base URL of the events service | `http://localhost:8082` |
| `EVENTS_API_KEY` | API key to authenticate requests | `my-secret-key` |

### What the events service must expose

For the integration to work, the external service must:

- Accept `POST /events` with `Content-Type: application/json`
- Validate the `X-API-Key` header
- Accept a body with the fields `stream_id`, `type`, `source`, `payload`, and `metadata`

If the events service responds with a status `>= 300`, the user creation fails and the error is returned to the client.

### Swapping the events adapter

The `Publisher` port allows replacing the transport without touching business logic. To implement a different adapter (Kafka, RabbitMQ, etc.), implement the interface and inject it in `module.go`:

```go
type Publisher interface {
    Publish(ctx context.Context, event Event) error
}
```
