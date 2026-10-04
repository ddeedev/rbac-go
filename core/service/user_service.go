package service

import (
	"context"

	"github.com/ddeedev/rbac-go/core/domain"
	"github.com/ddeedev/rbac-go/core/ports"
)

type userService struct {
	repo ports.UserRepository
}

// Create implements [ports.UserService].
func (u *userService) Create(ctx context.Context, name string, email string, password string) (*domain.User, error) {
	panic("unimplemented")
}

// Delete implements [ports.UserService].
func (u *userService) Delete(ctx context.Context, id string) (*domain.User, error) {
	panic("unimplemented")
}

// Get implements [ports.UserService].
func (u *userService) Get(ctx context.Context, id string) (*domain.User, error) {
	panic("unimplemented")
}

// Login implements [ports.UserService].
func (u *userService) Login(ctx context.Context, email string, password string) (token string, err error) {
	panic("unimplemented")
}

// Update implements [ports.UserService].
func (u *userService) Update(ctx context.Context, id string, name string, email string) (*domain.User, error) {
	panic("unimplemented")
}

func NewUserService(repo ports.UserRepository) ports.UserService {
	return &userService{repo: repo}
}
