package ports

import (
	"context"

	"github.com/ddeedev/rbac-go/core/domain"
)

// outboud port required for loading data from mongo
type UserRepository interface {
	Create(ctx context.Context, u *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, u *domain.User) error
	Delete(ctx context.Context, id string) error
}

// inbound port served to clietn
type UserService interface {
	Create(ctx context.Context, name, email, password string) (*domain.User, error)
	Get(ctx context.Context, id string) (*domain.User, error)
	Update(ctx context.Context, id, name, email string) (*domain.User, error)
	Delete(ctx context.Context, id string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (token string, err error)
}
