package ports

import (
	"github.com/ddeedev/rbac-go/core/domain"
)

type PasswordHasher interface {
	Hash(password string, c *domain.Config) (string, error)
	Compare(password, encodedHash string) (bool, error)
}