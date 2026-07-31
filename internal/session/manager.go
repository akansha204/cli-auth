package session

import (
	"errors"
	"time"

	"github.com/akansha204/cli-auth/internal/models"
	"github.com/akansha204/cli-auth/internal/repository"
)

var (
	ErrNoSession      = errors.New("no active session")
	ErrSessionExpired = errors.New("session has expired")
)

type Manager struct {
	repo    *repository.SessionRepository
	timeout time.Duration
}

func NewManager(repo *repository.SessionRepository, timeout time.Duration) *Manager {
	return &Manager{
		repo:    repo,
		timeout: timeout,
	}
}

func (m *Manager) Create(userID uint) (*models.Session, error) {
	session := &models.Session{
		UserID:    userID,
		ExpiresAt: time.Now().Add(m.timeout),
		Active:    true,
	}

	if err := m.repo.Create(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (m *Manager) Active(userID uint) (*models.Session, error) {
	session, err := m.repo.FindActiveByUserID(userID)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, ErrNoSession
	}

	if !session.ExpiresAt.After(time.Now()) {
		if err := m.repo.Invalidate(session.ID); err != nil {
			return nil, err
		}

		return nil, ErrSessionExpired
	}

	return session, nil
}

func (m *Manager) Refresh(session *models.Session) error {
	session.ExpiresAt = time.Now().Add(m.timeout)

	return m.repo.Update(session)
}

func (m *Manager) Invalidate(userID uint) error {
	session, err := m.repo.FindActiveByUserID(userID)
	if err != nil {
		return err
	}

	if session == nil {
		return ErrNoSession
	}

	return m.repo.Invalidate(session.ID)
}
