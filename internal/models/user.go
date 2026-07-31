package models

import "time"

type User struct {
	ID uint `gorm:"primaryKey"`

	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`

	MFAEnabled bool
	TOTPSecret string

	FailedAttempts int
	LockedUntil    *time.Time

	RegisteredAt time.Time
	LastLogin    *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
