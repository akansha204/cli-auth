package repository

import (
	"github.com/akansha204/cli-auth/internal/models"

	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (r *SessionRepository) Create(session *models.Session) error {
	return r.db.Create(session).Error
}

func (r *SessionRepository) FindActiveByUserID(userID uint) (*models.Session, error) {
	var session models.Session

	err := r.db.
		Where("user_id = ? AND active = ?", userID, true).
		First(&session).Error

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *SessionRepository) Invalidate(sessionID uint) error {
	return r.db.
		Model(&models.Session{}).
		Where("id = ?", sessionID).
		Update("active", false).Error
}
