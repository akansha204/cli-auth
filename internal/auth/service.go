package auth

import (
	"strings"
	"time"

	"github.com/akansha204/cli-auth/internal/config"
	"github.com/akansha204/cli-auth/internal/models"
	"github.com/akansha204/cli-auth/internal/repository"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (a *AuthService) Register(username, password string) error {

	username = strings.ToLower(strings.TrimSpace(username))
	password = strings.TrimSpace(password)

	if len(username) < 3 {
		return ErrInvalidUsername
	}
	if len(password) < 8 {
		return ErrInvalidPassword
	}

	exists, err := a.userRepo.Exists(username)
	if err != nil {
		return err
	}

	if exists {
		return ErrUserAlreadyExists
	}

	hash, err := HashPassword(password)
	if err != nil {
		return err
	}

	user := &models.User{
		Username: username,

		PasswordHash: hash,

		MFAEnabled: false,

		FailedAttempts: 0,

		RegisteredAt: time.Now(),
	}

	return a.userRepo.Create(user)
}

func (a *AuthService) Login(username, password string) (*models.User, error) {

	username = strings.ToLower(strings.TrimSpace(username))
	password = strings.TrimSpace(password)

	user, err := a.userRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if IsLocked(user) {
		return nil, ErrAccountLocked
	}

	if err := VerifyPassword(user.PasswordHash, password); err != nil {

		IncrementFailedAttempts(
			user,
			config.AppConfig.MaxLoginAttempts,
			config.AppConfig.LockoutDuration,
		)

		if err := a.userRepo.Update(user); err != nil {
			return nil, err
		}

		return nil, ErrInvalidCredentials
	}

	ResetFailedAttempts(user)

	now := time.Now()
	user.LastLogin = &now

	if err := a.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
