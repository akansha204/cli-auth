package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type App struct {
	DatabasePath     string
	SessionTimeout   int
	LockoutDuration  int
	MaxLoginAttempts int
}

var AppConfig = App{
	DatabasePath:     "data/app.db",
	SessionTimeout:   3600,
	LockoutDuration:  900,
	MaxLoginAttempts: 5,
}

func Load() error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("load .env: %w", err)
	}

	AppConfig = App{
		DatabasePath:     getEnv("DATABASE_PATH", AppConfig.DatabasePath),
		SessionTimeout:   getEnvInt("SESSION_TIMEOUT", AppConfig.SessionTimeout),
		LockoutDuration:  getEnvInt("LOCKOUT_DURATION", AppConfig.LockoutDuration),
		MaxLoginAttempts: getEnvInt("MAX_LOGIN_ATTEMPTS", AppConfig.MaxLoginAttempts),
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return fallback
}
