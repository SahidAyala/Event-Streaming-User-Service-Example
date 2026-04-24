package domain

import (
	"errors"
	"time"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

type User struct {
	ID           string
	Email        string
	Username     string
	PasswordHash string
	Status       Status
	Version      int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}

var (
	ErrNotFound           = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)
