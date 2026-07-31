package database

import (
	"log"
	"os"

	"github.com/akansha204/cli-auth/internal/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Initialize() error {
	db, err := gorm.Open(sqlite.Open(config.AppConfig.DatabasePath), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			LogLevel:                  logger.Error,
			IgnoreRecordNotFoundError: true,
		}),
	})
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
