package database

import (
	"log"

	"github.com/akansha204/cli-auth/internal/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Initialize() error {
	db, err := gorm.Open(sqlite.Open(config.AppConfig.DatabasePath), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db

	log.Println("Connected to SQLite")

	if err := Migrate(DB); err != nil {
		return err
	}

	log.Println("Database migrated")
	return nil
}
