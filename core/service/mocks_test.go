package service

import (
	"context"
	"time"

	"github.com/ddeedev/rbac-go/core/domain"
)

// mockRepo to mimic interface and signature of ports.UserRepository.
type mockRepo struct {
	createFn     func(ctx context.Context, u *domain.User) error
	getByIDFn    func(ctx context.Context, id string) (*domain.User, error)
	getByEmailFn func(ctx context.Context, email string) (*domain.User, error)
	getAllFn     func(ctx context.Context) ([]*domain.User, error)
	updateFn     func(ctx context.Context, u *domain.User) error
	deleteFn     func(ctx context.Context, id string) error
}

func (m *mockRepo) Create(ctx context.Context, u *domain.User) error {
	return m.createFn(ctx, u)
}

func (m *mockRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return m.getByIDFn(ctx, id)
}

func (m *mockRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.getByEmailFn(ctx, email)
}

func (m *mockRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return m.getByEmailFn(ctx, username)
}

func (m *mockRepo) GetAll(ctx context.Context) ([]*domain.User, error) {
	return m.getAllFn(ctx)
}

func (m *mockRepo) Update(ctx context.Context, u *domain.User) error {
	return m.updateFn(ctx, u)
}

func (m *mockRepo) Delete(ctx context.Context, id string) error {
	return m.deleteFn(ctx, id)
}

func (m *mockRepo) UpdatePassword(ctx context.Context, id, hashedPassword string) error {
	return nil
}

func (m *mockRepo) Count(ctx context.Context) (int64, error) {
	return 0, nil
}

// mockHasher mimic ports.PasswordHasher.
type mockHasher struct {
	hashFn    func(password string, c *domain.Config) (string, error)
	compareFn func(password, encodedHash string) (bool, error)
}

func (m *mockHasher) Hash(password string, c *domain.Config) (string, error) {
	return m.hashFn(password, c)
}

func (m *mockHasher) Compare(password, encodedHash string) (bool, error) {
	return m.compareFn(password, encodedHash)
}

// mockTokens mimic ports.TokenManager.
type mockTokens struct {
	issueFn  func(u *domain.User) (string, time.Duration, error)
	verifyFn func(token string) (*domain.TokenClaims, error)
}

func (m *mockTokens) Issue(u *domain.User) (string, time.Duration, error) {
	return m.issueFn(u)
}

func (m *mockTokens) Verify(token string) (*domain.TokenClaims, error) {
	return m.verifyFn(token)
}
