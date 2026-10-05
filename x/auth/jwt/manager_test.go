package jwt

import (
	"errors"
	"testing"
	"time"

	"github.com/ddeedev/rbac-go/core/domain"
)

const testSecret = "test-secret-at-least-32-bytes-long!!"

func TestNewManagerValidation(t *testing.T) {
	if _, err := NewManager("too-short", time.Hour); err == nil {
		t.Error("expected error for secret shorter than 32 bytes")
	}
	if _, err := NewManager(testSecret, 0); err == nil {
		t.Error("expected error for non-positive ttl")
	}
}

func TestManagerIssueVerifyRoundTrip(t *testing.T) {
	m, err := NewManager(testSecret, time.Hour)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	user := &domain.User{ID: "u1", Email: "alice@example.com"}
	token, ttl, err := m.Issue(user)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	if ttl != time.Hour {
		t.Errorf("ttl = %v, want 1h", ttl)
	}

	claims, err := m.Verify(token)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if claims.UserID != "u1" || claims.Email != "alice@example.com" {
		t.Errorf("claims = %+v, want UserID=u1 Email=alice@example.com", claims)
	}
	if claims.TokenType != domain.TokenTypeAccess {
		t.Errorf("TokenType = %q, want %q", claims.TokenType, domain.TokenTypeAccess)
	}
}

func TestManagerVerifyWrongSecret(t *testing.T) {
	issuer, _ := NewManager(testSecret, time.Hour)
	attacker, _ := NewManager("another-secret-also-32-bytes-long!!!", time.Hour)

	token, _, _ := issuer.Issue(&domain.User{ID: "u1"})

	if _, err := attacker.Verify(token); !errors.Is(err, domain.ErrInvalidToken) {
		t.Fatalf("err = %v, want ErrInvalidToken (signature must not verify)", err)
	}
}

func TestManagerVerifyExpired(t *testing.T) {
	m, _ := NewManager(testSecret, time.Nanosecond)

	token, _, _ := m.Issue(&domain.User{ID: "u1"})
	// make it expire
	time.Sleep(5 * time.Millisecond)

	if _, err := m.Verify(token); !errors.Is(err, domain.ErrInvalidToken) {
		t.Fatalf("err = %v, want ErrInvalidToken for expired token", err)
	}
}
