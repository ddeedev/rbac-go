package service

import (
	"context"
	"errors"
	"time"

	"github.com/ddeedev/rbac-go/core/domain"
	"github.com/ddeedev/rbac-go/core/ports"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) ports.UserService {
	return &userService{repo: repo}
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
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	// create domain user
	du := &domain.User{Name: name, Email: email, Password: string(hash), CreatedAt: time.Now().UTC()}
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

// TODO:: Move to new module or implemnt here
// Login implements [ports.UserService].
func (u *userService) Login(ctx context.Context, email string, password string) (token string, err error) {
	panic("unimplemented")
}
