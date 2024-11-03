package store

import "github.com/drewbuiltit/trading-journal/backend/internal/models"

type Store interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
}
