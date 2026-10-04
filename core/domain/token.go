package domain

import (
	"github.com/golang-jwt/jwt/v5"
)

// TokenType represents the type of JWT token issued by the system.
type TokenType string

const (
	TokenTypeAccess     TokenType = "access_token"
	// TODO: too much for poc
	// TokenTypeFirstLogin TokenType = "first_login"
	// TokenTypeRefresh    TokenType = "refresh_token"
)

type TokenClaims struct {
	UserID    string     `json:"user_id"`
	Username  string     `json:"username"` // email
	TokenType *TokenType `json:"token_type,omitempty"`
	jwt.RegisteredClaims
}
