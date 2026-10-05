package user

import (
	"context"
	"errors"
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

func objectID(id string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, domain.ErrNotFound
	}
	return oid, nil
}

func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	s := userSchema{
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
	}
	res, err := r.collection.InsertOne(ctx, s)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.ErrEmailTaken
		}
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		u.ID = oid.Hex()
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	oid, err := objectID(id)
	if err != nil {
		return err
	}
	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.findOne(ctx, bson.M{"email": email})
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	oid, err := objectID(id)
	if err != nil {
		return nil, err
	}
	return r.findOne(ctx, bson.M{"_id": oid})
}

func (r *UserRepo) Update(ctx context.Context, u *domain.User) error {
	oid, err := objectID(u.ID)
	if err != nil {
		return err
	}
	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": oid},
		bson.M{"$set": bson.M{"name": u.Name, "email": u.Email}})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.ErrEmailTaken
		}
		return err
	}
	if res.MatchedCount == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *UserRepo) GetAll(ctx context.Context) ([]*domain.User, error) {
	cur, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []*domain.User
	for cur.Next(ctx) {
		var s userSchema
		if err := cur.Decode(&s); err != nil {
			return nil, err
		}
		out = append(out, s.toDomain())
	}
	return out, cur.Err()
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return r.GetByEmail(ctx, username)
}

func (r *UserRepo) UpdatePassword(ctx context.Context, id string, hashedPassword string) error {
	oid, err := objectID(id)
	if err != nil {
		return err
	}
	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": oid},
		bson.M{"$set": bson.M{"password_hash": hashedPassword}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *UserRepo) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

// findOne runs a single-document query and maps a miss to ErrNotFound.
func (r *UserRepo) findOne(ctx context.Context, filter bson.M) (*domain.User, error) {
	var s userSchema
	err := r.collection.FindOne(ctx, filter).Decode(&s)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return s.toDomain(), nil
}
