package service

import (
	"context"
	"errors"
	"time"

	"github.com/ddeedev/rbac-go/core/domain"
	"github.com/ddeedev/rbac-go/core/ports"
)

type userService struct {
	repo   ports.UserRepository
	hasher ports.PasswordHasher
}

func NewUserService(repo ports.UserRepository, hasher ports.PasswordHasher) ports.UserService {
	return &userService{repo: repo, hasher: hasher}
}

// create user service
func (u *userService) Create(ctx context.Context, name string, email string, password string) (*domain.User, error) {
	// validate user email and duplication
	if _, err := u.repo.GetByEmail(ctx, email); err == nil {
		return nil, domain.ErrEmailTaken
	} else if !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}

	// hash user password before store value
	hash, err := u.hasher.Hash(password, nil)
	if err != nil {
		return nil, err
	}

	// create domain user
	du := &domain.User{Name: name, Email: email, Password: hash, CreatedAt: time.Now().UTC()}
	if err := u.repo.Create(ctx, du); err != nil {
		return nil, err
	}

	return du, nil
}

// get user service
func (u *userService) Get(ctx context.Context, id string) (*domain.User, error) {
	return u.repo.GetByID(ctx, id)
}

// set user service
func (u *userService) Update(ctx context.Context, id string, name string, email string) (*domain.User, error) {
	du, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	du.Name, du.Email = name, email
	if err := u.repo.Update(ctx, du); err != nil {
		return nil, err
	}

	return du, nil
}

// delete user service
func (u *userService) Delete(ctx context.Context, id string) (*domain.User, error) {
	// find exiting user
	du, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return du, u.repo.Delete(ctx, id)
}
