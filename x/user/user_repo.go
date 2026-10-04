package user

import (
	"context"
	"time"

	"github.com/ddeedev/rbac-go/core/domain"
	"github.com/ddeedev/rbac-go/core/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const UserCollectionName = "users"

// userSchema is the MongoDB storage model for a user.
type userSchema struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password_hash"` // already hashed by the service
	CreatedAt time.Time          `bson:"created_at"`
}

// unmarshal data from mongo to domain for better usage
func (s userSchema) toDomain() *domain.User {
	return &domain.User{
		ID:        s.ID.Hex(),
		Name:      s.Name,
		Email:     s.Email,
		Password:  s.Password,
		CreatedAt: s.CreatedAt,
	}
}

// compile-time check that UserRepo satisfies the port.
var _ ports.UserRepository = (*UserRepo)(nil)

// UserRepo is the MongoDB adapter for ports.UserRepository.
type UserRepo struct {
	collection *mongo.Collection
}

func NewUserRepo(ctx context.Context, db *mongo.Database) (*UserRepo, error) {
	col := db.Collection(UserCollectionName)
	_, err := col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}
	return &UserRepo{collection: col}, nil
}

// Create implements [ports.UserRepository].
func (*UserRepo) Create(ctx context.Context, u *domain.User) error {
	panic("unimplemented")
}

// Delete implements [ports.UserRepository].
func (u *UserRepo) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// GetByEmail implements [ports.UserRepository].
func (u *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	panic("unimplemented")
}

// GetByID implements [ports.UserRepository].
func (u *UserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	panic("unimplemented")
}

// Update implements [ports.UserRepository].
func (*UserRepo) Update(ctx context.Context, u *domain.User) error {
	panic("unimplemented")
}

// GetAll implements [ports.UserRepository].
func (u *UserRepo) GetAll(ctx context.Context) ([]*domain.User, error) {
	panic("unimplemented")
}

// GetByUsername implements [ports.UserRepository].
func (u *UserRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	panic("unimplemented")
}

// UpdatePassword implements [ports.UserRepository].
func (u *UserRepo) UpdatePassword(ctx context.Context, id string, hashedPassword string) error {
	panic("unimplemented")
}