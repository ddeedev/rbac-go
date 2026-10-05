package jwt

import (
	"errors"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"

	"github.com/ddeedev/rbac-go/core/domain"
	"github.com/ddeedev/rbac-go/core/ports"
)

var _ ports.TokenManager = (*Manager)(nil)

// claims is the wire format. It exists only inside this package.
type claims struct {
	UserID    string           `json:"user_id"`
	Email     string           `json:"email"`
	TokenType domain.TokenType `json:"token_type"`
	gojwt.RegisteredClaims
}

type Manager struct {
	secret []byte
	ttl    time.Duration
}

func NewManager(secret string, ttl time.Duration) (*Manager, error) {
	if len(secret) < 32 {
		return nil, errors.New("jwt secret must be at least 32 bytes")
	}
	if ttl <= 0 {
		return nil, errors.New("jwt ttl must be positive")
	}
	return &Manager{secret: []byte(secret), ttl: ttl}, nil
}

func (m *Manager) Issue(u *domain.User) (string, time.Duration, error) {
	now := time.Now()
	c := claims{
		UserID:    u.ID,
		Email:     u.Email,
		TokenType: domain.TokenTypeAccess,
		Subject:   u.ID,
		IssuedAt:  gojwt.NewNumericDate(now),
		ExpiresAt: gojwt.NewNumericDate(now.Add(m.ttl)),
	}
	
	signed, err := gojwt.NewWithClaims(gojwt.SigningMethodHS256, c).SignedString(m.secret)
	if err != nil {
		return "", 0, err
	}
	return signed, m.ttl, nil
}

func (m *Manager) Verify(tokenStr string) (*domain.TokenClaims, error) {
	var c claims
	token, err := gojwt.ParseWithClaims(tokenStr, &c,
		func(*gojwt.Token) (any, error) { return m.secret, nil },
		gojwt.WithValidMethods([]string{"HS256"}),
		gojwt.WithExpirationRequired(),
	)
	if err != nil || !token.Valid {
		return nil, domain.ErrInvalidToken
	}
	return &domain.TokenClaims{
		UserID:    c.UserID,
		Email:     c.Email,
		TokenType: c.TokenType,
		ExpiresAt: c.ExpiresAt.Time,
	}, nil
}
