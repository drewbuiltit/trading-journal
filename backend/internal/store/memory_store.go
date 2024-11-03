package store

import (
	"errors"
	"github.com/drewbuiltit/trading-journal/backend/internal/models"
	"sync"
)

type MemoryStore struct {
	users map[string]*models.User
	mu    sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users: make(map[string]*models.User),
	}
}

func (m *MemoryStore) CreateUser(user *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.users[user.Email]; exists {
		return errors.New("user already exists")
	}

	user.ID = len(m.users) + 1
	m.users[user.Email] = user
	return nil
}

func (m *MemoryStore) GetUserByEmail(email string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, exists := m.users[email]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}
