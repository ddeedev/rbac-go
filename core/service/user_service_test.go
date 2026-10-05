package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ddeedev/rbac-go/core/domain"
)

func TestUserServiceCreateSuccess(t *testing.T) {
	repo := &mockRepo{
		getByEmailFn: func(context.Context, string) (*domain.User, error) {
			// no existing user
			return nil, domain.ErrNotFound
		},
		createFn: func(_ context.Context, u *domain.User) error {
			u.ID = "generated-id"
			return nil
		},
	}
	hasher := &mockHasher{
		hashFn: func(string, *domain.Config) (string, error) { return "hashed-pw", nil },
	}

	svc := NewUserService(repo, hasher)
	u, err := svc.Create(context.Background(), "Alice", "alice@example.com", "supersecret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.ID != "generated-id" {
		t.Errorf("ID = %q, want generated-id", u.ID)
	}
	if u.Password != "hashed-pw" {
		t.Errorf("Password = %q, want the hashed value (never the plaintext)", u.Password)
	}
	if u.Email != "alice@example.com" || u.Name != "Alice" {
		t.Errorf("unexpected user fields: %+v", u)
	}
}

func TestUserServiceCreateDuplicateEmail(t *testing.T) {
	repo := &mockRepo{
		getByEmailFn: func(context.Context, string) (*domain.User, error) {
			// already exists
			return &domain.User{Email: "alice@example.com"}, nil
		},
	}
	svc := NewUserService(repo, &mockHasher{})

	_, err := svc.Create(context.Background(), "Alice", "alice@example.com", "supersecret")
	if !errors.Is(err, domain.ErrEmailTaken) {
		t.Fatalf("err = %v, want ErrEmailTaken", err)
	}
}

func TestUserServiceCreateRepoLookupError(t *testing.T) {
	boom := errors.New("db down")
	repo := &mockRepo{
		getByEmailFn: func(context.Context, string) (*domain.User, error) { return nil, boom },
	}
	svc := NewUserService(repo, &mockHasher{})

	_, err := svc.Create(context.Background(), "Alice", "alice@example.com", "supersecret")
	if !errors.Is(err, boom) {
		t.Fatalf("err = %v, want the propagated repo error", err)
	}
}

func TestUserServiceUpdateSuccess(t *testing.T) {
	repo := &mockRepo{
		getByIDFn: func(_ context.Context, id string) (*domain.User, error) {
			return &domain.User{ID: id, Name: "Alice_Old", Email: "alice_old@example.com"}, nil
		},
		updateFn: func(context.Context, *domain.User) error { return nil },
	}
	svc := NewUserService(repo, &mockHasher{})

	u, err := svc.Update(context.Background(), "id1", "Alice_new", "alice_new@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Name != "Alie_new" || u.Email != "alice_new@example.com" {
		t.Errorf("fields not updated: %+v", u)
	}
}

func TestUserServiceUpdateNotFound(t *testing.T) {
	repo := &mockRepo{
		getByIDFn: func(context.Context, string) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
	}
	svc := NewUserService(repo, &mockHasher{})

	_, err := svc.Update(context.Background(), "bob_id", "New", "new@example.com")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestUserServiceDeleteSuccess(t *testing.T) {
	deleted := ""
	repo := &mockRepo{
		getByIDFn: func(_ context.Context, id string) (*domain.User, error) {
			return &domain.User{ID: id}, nil
		},
		deleteFn: func(_ context.Context, id string) error { deleted = id; return nil },
	}
	svc := NewUserService(repo, &mockHasher{})

	u, err := svc.Delete(context.Background(), "id1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.ID != "id1" || deleted != "id1" {
		t.Errorf("delete did not target id1: returned=%+v deleted=%q", u, deleted)
	}
}

func TestUserServiceList(t *testing.T) {
	want := []*domain.User{{ID: "1", Name: "Alice", Email: "alice@example.com"}, {ID: "2", Name: "Bob", Email: "bob@exmapele.com"}}
	repo := &mockRepo{
		getAllFn: func(context.Context) ([]*domain.User, error) { return want, nil },
	}
	svc := NewUserService(repo, &mockHasher{})

	got, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
}
