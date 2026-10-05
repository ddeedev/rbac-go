// core/ports/auth.go
package ports

import (
	"context"
	"time"

	"github.com/ddeedev/rbac-go/core/domain"
)

// authservice => inbound port used by the gRPC handler and the auth interceptor.
type AuthService interface {
	Register(ctx context.Context, name, email, password string) (*domain.RegisterResult, error)
	Login(ctx context.Context, email, password string) (*domain.AuthResult, error)
	ValidateToken(ctx context.Context, token string) (*domain.TokenClaims, error)
}

// tokenmngr is an outbound port for signing and verifying access tokens.
type TokenManager interface {
	Issue(u *domain.User) (token string, expiresIn time.Duration, err error)
	Verify(token string) (*domain.TokenClaims, error)
}
