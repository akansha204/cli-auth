package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type App struct {
	DatabasePath     string
	SessionTimeout   time.Duration
	LockoutDuration  time.Duration
	MaxLoginAttempts int
}

var AppConfig = App{
	DatabasePath:     "data/app.db",
	SessionTimeout:   time.Hour,
	LockoutDuration:  15 * time.Minute,
	MaxLoginAttempts: 5,
}

func Load() error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("load .env: %w", err)
	}

	AppConfig = App{
		DatabasePath:     getEnv("DATABASE_PATH", AppConfig.DatabasePath),
		SessionTimeout:   getEnvDuration("SESSION_TIMEOUT", AppConfig.SessionTimeout),
		LockoutDuration:  getEnvDuration("LOCKOUT_DURATION", AppConfig.LockoutDuration),
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

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return fallback
}
