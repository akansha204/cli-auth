package auth

import (
	"time"

	"github.com/akansha204/cli-auth/internal/models"
)

func IsLocked(user *models.User) bool {
	if user.LockedUntil == nil {
		return false
	}

	return time.Now().Before(*user.LockedUntil)
}

func ResetFailedAttempts(user *models.User) {
	user.FailedAttempts = 0
	user.LockedUntil = nil
}

func IncrementFailedAttempts(
	user *models.User,
	maxAttempts int,
	lockoutDuration time.Duration,
) {
	user.FailedAttempts++

	if user.FailedAttempts >= maxAttempts {
		lockUntil := time.Now().Add(lockoutDuration)
		user.LockedUntil = &lockUntil
	}
}
