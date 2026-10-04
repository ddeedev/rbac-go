package domain

import (
	"errors"
	"time"
)

var (
	ErrNotFound     = errors.New("user not found")
	ErrInvalidInput = errors.New("invalid input data")
	ErrEmailTaken   = errors.New("email alrady in use")
	ErrInvalidLogin = errors.New("invalid credentials")
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}
