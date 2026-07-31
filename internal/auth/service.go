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

	if user.MFAEnabled {
		if err := a.userRepo.Update(user); err != nil {
			return nil, err
		}

		return nil, ErrMFARequired
	}

	now := time.Now()
	user.LastLogin = &now

	if err := a.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (a *AuthService) VerifyMFA(username, code string) (*models.User, error) {
	user, err := a.userRepo.FindByUsername(strings.ToLower(strings.TrimSpace(username)))
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if !user.MFAEnabled {
		return nil, ErrInvalidTOTP
	}

	if !ValidateTOTPCode(user.TOTPSecret, strings.TrimSpace(code)) {
		return nil, ErrInvalidTOTP
	}

	now := time.Now()
	user.LastLogin = &now

	if err := a.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (a *AuthService) EnableMFA(username string) (string, string, error) {
	user, err := a.userRepo.FindByUsername(strings.ToLower(strings.TrimSpace(username)))
	if err != nil {
		return "", "", err
	}

	if user == nil {
		return "", "", ErrUserNotFound
	}

	if user.MFAEnabled {
		return "", "", ErrMFAEnabled
	}

	secret, uri, err := GenerateTOTPSecret(user.Username)
	if err != nil {
		return "", "", err
	}

	user.TOTPSecret = secret
	user.MFAEnabled = true

	if err := a.userRepo.Update(user); err != nil {
		return "", "", err
	}

	return secret, uri, nil
}

func (a *AuthService) DisableMFA(username string) error {
	user, err := a.userRepo.FindByUsername(strings.ToLower(strings.TrimSpace(username)))
	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	if !user.MFAEnabled {
		return ErrMFANotEnabled
	}

	user.TOTPSecret = ""
	user.MFAEnabled = false

	return a.userRepo.Update(user)
}

func (a *AuthService) HasMFA(username string) (bool, error) {
	user, err := a.userRepo.FindByUsername(strings.ToLower(strings.TrimSpace(username)))
	if err != nil {
		return false, err
	}

	if user == nil {
		return false, ErrUserNotFound
	}

	return user.MFAEnabled, nil
}
