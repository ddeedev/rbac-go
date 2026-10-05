package domain

import (
	"context"
	"time"
)

type TokenType string

const (
	TokenTypeAccess TokenType = "access_token"
	// NOTE: too much for poc
	// TokenTypeFirstLogin TokenType = "first_login"
	// TokenTypeRefresh    TokenType = "refresh_token"
)

// TokenClaims is what a validated token tells us about the caller.
type TokenClaims struct {
	UserID    string
	Email     string
	TokenType TokenType
	ExpiresAt time.Time
}

type AuthResult struct {
	AccessToken string
	ExpiresIn   time.Duration
	User        *User
}

// token helper funtion
type ctxKey struct{}

func WithClaims(ctx context.Context, c *TokenClaims) context.Context {
	return context.WithValue(ctx, ctxKey{}, c)
}

func ClaimsFromContext(ctx context.Context) (*TokenClaims, bool) {
	c, ok := ctx.Value(ctxKey{}).(*TokenClaims)
	return c, ok
}

