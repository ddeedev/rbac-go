package service

import (
	"context"
	"errors"
	"time"

	"github.com/ddeedev/rbac-go/core/domain"
	"github.com/ddeedev/rbac-go/core/ports"
)

type authService struct {
	repo   ports.UserRepository
	hasher ports.PasswordHasher
	tokens ports.TokenManager
}

func NewAuthService(repo ports.UserRepository, hasher ports.PasswordHasher, tokens ports.TokenManager) ports.AuthService {
	return &authService{repo: repo, hasher: hasher, tokens: tokens}
}

func (a *authService) Register(ctx context.Context, name string, email string, password string) (*domain.RegisterResult, error) {
	// validate user email and duplication
	if _, err := a.repo.GetByEmail(ctx, email); err == nil {
		return nil, domain.ErrEmailTaken
	} else if !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}

	// hash user password before store value
	hash, err := a.hasher.Hash(password, nil)
	if err != nil {
		return nil, err
	}

	// create domain user
	du := &domain.User{Name: name, Email: email, Password: hash, CreatedAt: time.Now().UTC()}
	if err := a.repo.Create(ctx, du); err != nil {
		return nil, err
	}

	return &domain.RegisterResult{Success: true}, nil
}

// loigin verify credentials and issues an access token.
func (a *authService) Login(ctx context.Context, email, password string) (*domain.AuthResult, error) {
	u, err := a.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrInvalidLogin
		}
		return nil, err
	}

	// compare input hashed input password with hashed on stored
	ok, err := a.hasher.Compare(password, u.Password)
	if err != nil || !ok {
		return nil, domain.ErrInvalidLogin
	}

	// create jwt token
	token, ttl, err := a.tokens.Issue(u)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResult{AccessToken: token, ExpiresIn: ttl, User: u}, nil
}

// ValidateToken verifies an access token and returns its claims.
func (a *authService) ValidateToken(_ context.Context, token string) (*domain.TokenClaims, error) {
	return a.tokens.Verify(token)
}
