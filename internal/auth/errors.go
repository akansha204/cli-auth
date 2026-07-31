package auth

import "errors"

var (
	ErrUserAlreadyExists = errors.New("username already exists")
	ErrUserNotFound      = errors.New("user not found")

	ErrInvalidUsername = errors.New("username must be at least 3 characters")
	ErrInvalidPassword = errors.New("password must be at least 8 characters")

	ErrInvalidCredentials = errors.New("invalid username or password")

	ErrAccountLocked = errors.New("account is temporarily locked")

	ErrMFARequired   = errors.New("MFA code required")
	ErrInvalidTOTP   = errors.New("invalid MFA code")
	ErrMFAEnabled    = errors.New("MFA is already enabled")
	ErrMFANotEnabled = errors.New("MFA is not enabled")
)
