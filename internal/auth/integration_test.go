package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/akansha204/cli-auth/internal/models"
	"github.com/akansha204/cli-auth/internal/repository"
	"github.com/akansha204/cli-auth/internal/session"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setup(t *testing.T) (*AuthService, *session.Manager) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Session{}); err != nil {
		t.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	return NewAuthService(userRepo), session.NewManager(sessionRepo, time.Minute)
}

func TestFullAuthFlow(t *testing.T) {
	svc, sessions := setup(t)

	if err := svc.Register("Alice", "password123"); err != nil {
		t.Fatal(err)
	}

	if err := svc.Register("Alice", "password123"); !errors.Is(err, ErrUserAlreadyExists) {
		t.Fatalf("expected ErrUserAlreadyExists, got %v", err)
	}

	if _, err := svc.Login("alice", "wrongpass"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}

	user, err := svc.Login("alice", "password123")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := sessions.Create(user.ID); err != nil {
		t.Fatal(err)
	}

	secret, _, err := svc.EnableMFA("alice")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := svc.Login("alice", "password123"); !errors.Is(err, ErrMFARequired) {
		t.Fatalf("expected ErrMFARequired, got %v", err)
	}

	if _, err := svc.VerifyMFA("alice", "000000"); !errors.Is(err, ErrInvalidTOTP) {
		t.Fatalf("expected ErrInvalidTOTP, got %v", err)
	}

	code := currentCode(t, secret)
	if _, err := svc.VerifyMFA("alice", code); err != nil {
		t.Fatal(err)
	}

	if err := svc.DisableMFA("alice"); err != nil {
		t.Fatal(err)
	}

	if err := svc.DisableMFA("alice"); !errors.Is(err, ErrMFANotEnabled) {
		t.Fatalf("expected ErrMFANotEnabled, got %v", err)
	}
}

func TestAccountLockout(t *testing.T) {
	svc, _ := setup(t)

	if err := svc.Register("bob", "password123"); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		if _, err := svc.Login("bob", "wrongpass"); !errors.Is(err, ErrInvalidCredentials) {
			t.Fatalf("attempt %d: expected ErrInvalidCredentials, got %v", i+1, err)
		}
	}

	if _, err := svc.Login("bob", "password123"); !errors.Is(err, ErrAccountLocked) {
		t.Fatalf("expected ErrAccountLocked, got %v", err)
	}
}
