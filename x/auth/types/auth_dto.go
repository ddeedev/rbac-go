package types

import "github.com/ddeedev/rbac-go/core/domain"

type LoginBodyDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponseDTO struct {
	AccessToken string `json:"access_token"`
	// RefreshToken string       `json:"refresh_token,omitempty"`
	ExpiresIn int          `json:"expires_in"`
	TokenType string       `json:"token_type"`
	User      *domain.User `json:"user"`
}

