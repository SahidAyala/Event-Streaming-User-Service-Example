package application

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/SahidAyala/Event-Streaming-User-Service-Example/internal/shared/events"
	"github.com/SahidAyala/Event-Streaming-User-Service-Example/internal/user/application/ports"
	"github.com/SahidAyala/Event-Streaming-User-Service-Example/internal/user/domain"
)

type Service struct {
	repo      ports.UserRepository
	publisher events.Publisher
}

func NewService(repo ports.UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SetPublisher(publisher events.Publisher) {
	s.publisher = publisher
}

func (s *Service) CreateUser(ctx context.Context, email, username, password string) (*domain.User, error) {
	if email == "" || username == "" || password == "" {
		return nil, errors.New("missing required fields")
	}

	existing, _ := s.repo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, domain.ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &domain.User{
		ID:           uuid.NewString(),
		Email:        email,
		Username:     username,
		PasswordHash: string(hashedPassword),
		Status:       domain.StatusActive,
		Version:      1,
	}

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}

	if s.publisher != nil {
		if err := s.publisher.Publish(ctx, events.Event{
			StreamID: "user:" + u.ID,
			Type:     "user.created",
			Source:   "user-service",
			Payload: map[string]interface{}{
				"user_id": u.ID,
			},
			Metadata: map[string]interface{}{
				"email": u.Email,
			},
		}); err != nil {
			return nil, err
		}
	}

	return u, nil
}

func (s *Service) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) UpdateEmail(ctx context.Context, id, newEmail string) (*domain.User, error) {
	if newEmail == "" {
		return nil, errors.New("email is required")
	}

	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	existing, _ := s.repo.GetByEmail(ctx, newEmail)
	if existing != nil {
		return nil, domain.ErrEmailAlreadyExists
	}

	u.Email = newEmail
	u.Version++

	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) UpdatePassword(ctx context.Context, id, currentPassword, newPassword string) error {
	if currentPassword == "" || newPassword == "" {
		return errors.New("passwords are required")
	}

	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.ErrNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(currentPassword)); err != nil {
		return errors.New("invalid current password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hashedPassword)
	u.Version++

	return s.repo.Update(ctx, u)
}
