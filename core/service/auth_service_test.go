package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ddeedev/rbac-go/core/domain"
)

func TestAuthServiceRegisterSuccess(t *testing.T) {
	var created *domain.User

	repo := &mockRepo{
		getByEmailFn: func(context.Context, string) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
		createFn: func(_ context.Context, u *domain.User) error { created = u; return nil },
	}

	hasher := &mockHasher{
		hashFn: func(string, *domain.Config) (string, error) { return "hashed-pw", nil },
	}

	svc := NewAuthService(repo, hasher, &mockTokens{})

	res, err := svc.Register(context.Background(), "Alice", "alice@example.com", "supersecret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Success {
		t.Error("Success = false, want true")
	}
	if created == nil || created.Password != "hashed-pw" {
		t.Errorf("stored user should carry the hash, got %+v", created)
	}
}

func TestAuthServiceRegisterDuplicateEmail(t *testing.T) {
	repo := &mockRepo{
		getByEmailFn: func(context.Context, string) (*domain.User, error) {
			return &domain.User{}, nil
		},
	}
	svc := NewAuthService(repo, &mockHasher{}, &mockTokens{})

	_, err := svc.Register(context.Background(), "Alie", "alice@example.com", "supersecret")
	if !errors.Is(err, domain.ErrEmailTaken) {
		t.Fatalf("err = %v, want ErrEmailTaken", err)
	}
}

func TestAuthServiceLoginSuccess(t *testing.T) {
	stored := &domain.User{ID: "u1", Email: "alice@example.com", Password: "hashed-pw"}
	repo := &mockRepo{
		getByEmailFn: func(context.Context, string) (*domain.User, error) { return stored, nil },
	}
	hasher := &mockHasher{
		compareFn: func(pw, hash string) (bool, error) { return pw == "supersecret", nil },
	}
	tokens := &mockTokens{
		issueFn: func(u *domain.User) (string, time.Duration, error) {
			return "signed.jwt.token", time.Hour, nil
		},
	}
	svc := NewAuthService(repo, hasher, tokens)

	res, err := svc.Login(context.Background(), "alice@example.com", "supersecret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.AccessToken != "signed.jwt.token" {
		t.Errorf("AccessToken = %q, want signed.jwt.token", res.AccessToken)
	}
	if res.User != stored {
		t.Error("AuthResult should carry the authenticated user")
	}
}

func TestAuthServiceLoginUserNotFound(t *testing.T) {
	repo := &mockRepo{
		getByEmailFn: func(context.Context, string) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
	}
	svc := NewAuthService(repo, &mockHasher{}, &mockTokens{})

	// A missing user must look identical to a wrong password.
	_, err := svc.Login(context.Background(), "bob@example.com", "whatever")
	if !errors.Is(err, domain.ErrInvalidLogin) {
		t.Fatalf("err = %v, want ErrInvalidLogin", err)
	}
}

func TestAuthServiceLoginWrongPassword(t *testing.T) {
	repo := &mockRepo{
		getByEmailFn: func(context.Context, string) (*domain.User, error) {
			return &domain.User{Password: "hashed-pw"}, nil
		},
	}
	hasher := &mockHasher{
		compareFn: func(string, string) (bool, error) { return false, nil },
	}
	svc := NewAuthService(repo, hasher, &mockTokens{})

	_, err := svc.Login(context.Background(), "alice@example.com", "wrong")
	if !errors.Is(err, domain.ErrInvalidLogin) {
		t.Fatalf("err = %v, want ErrInvalidLogin", err)
	}
}

func TestAuthServiceValidateToken(t *testing.T) {
	want := &domain.TokenClaims{UserID: "u1", Email: "alice@example.com"}
	tokens := &mockTokens{
		verifyFn: func(tok string) (*domain.TokenClaims, error) {
			if tok != "good-token" {
				return nil, domain.ErrInvalidToken
			}
			return want, nil
		},
	}
	svc := NewAuthService(&mockRepo{}, &mockHasher{}, tokens)

	got, err := svc.ValidateToken(context.Background(), "good-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.UserID != "u1" {
		t.Errorf("UserID = %q, want u1", got.UserID)
	}
}
