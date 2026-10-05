package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/ddeedev/rbac-go/core/domain"
	"github.com/ddeedev/rbac-go/core/ports"
	"golang.org/x/crypto/argon2"
)

// DefaultConfig provides recommended parameters for Argon2id hashing.
var DefaultConfig = domain.Config{
	Memory:      64 * 1024, // 64 MB
	Iterations:  3,
	Parallelism: uint8(runtime.NumCPU()),
	SaltLength:  16,
	KeyLength:   32,
}

// config holds the parameters for Argon2id hashing.
type hasherImpl struct{}

// creates a new hasher
// if config is nil it uses DefaultConfig.
func NewHasher() ports.PasswordHasher {
	return &hasherImpl{}
}

func (h *hasherImpl) Hash(password string, c *domain.Config) (string, error) {
	if c == nil {
		c = &DefaultConfig
	}

	// gen a random salt
	salt := make([]byte, c.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// hashing the password by using Argon2id
	hash := argon2.IDKey([]byte(password), salt, c.Iterations, c.Memory, c.Parallelism, c.KeyLength)

	// encoding the salt and hash to base64 for storage
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// PHC
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, c.Memory, c.Iterations, c.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func (h *hasherImpl) Compare(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}

	var c domain.Config
	var version int
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &c.Memory, &c.Iterations, &c.Parallelism)
	if err != nil {
		return false, err
	}

	_, err = fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	c.KeyLength = uint32(len(decodedHash))

	comparisonHash := argon2.IDKey([]byte(password), salt, c.Iterations, c.Memory, c.Parallelism, c.KeyLength)

	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1, nil
}
