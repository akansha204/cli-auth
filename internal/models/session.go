package models

import "time"

type Session struct {
	ID uint `gorm:"primaryKey"`

	UserID uint `gorm:"index"`

	ExpiresAt time.Time
	Active    bool

	CreatedAt time.Time
}
