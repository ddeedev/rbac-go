package ports

import (
	"time"

	"github.com/ddeedev/rbac-go/x/auth/types"
	"github.com/ddeedev/rbac-go/core/domain"
)

type AuthUseCase interface {
	Login(input types.LoginBody, jwtSecret string, ttl time.Duration) (*types.LoginResponse, error)
	ValidateToken(tokenString string, jwtSecret string) (*domain.TokenClaims, error)
}
