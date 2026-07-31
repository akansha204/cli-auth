package database

import (
	"github.com/akansha204/cli-auth/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Session{},
	)
}
