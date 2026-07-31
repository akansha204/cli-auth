package auth

import "errors"

var (
	ErrUserAlreadyExists = errors.New("username already exists")
	ErrUserNotFound      = errors.New("user not found")

	ErrInvalidUsername = errors.New("username must be at least 3 characters")
	ErrInvalidPassword = errors.New("password must be at least 8 characters")

	ErrInvalidCredentials = errors.New("invalid username or password")

	ErrAccountLocked = errors.New("account is temporarily locked")
)
