package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/ddeedev/rbac-go/core/domain"
	"github.com/ddeedev/rbac-go/config"
)

type JWTService struct {
	secretKey string
}

func NewJWTService(cfg *config.Config) *JWTService {
	return &JWTService{
		secretKey: cfg.JWTSecret,
	}
}

func (s *JWTService) GenerateToken(claims *domain.TokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (uc *JWTService) ValidateToken(tokenString string, jwtSecret string) (*domain.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidLogin
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, domain.ErrInvalidLogin
	}
	claims, ok := token.Claims.(*domain.TokenClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidLogin
	}
	return claims, nil
}
