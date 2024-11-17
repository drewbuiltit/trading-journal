package store

import (
	"github.com/drewbuiltit/trading-journal/backend/internal/models"
	"gorm.io/gorm"
)

type PostgresStore struct {
	DB *gorm.DB
}

func NewPostgresStore(db *gorm.DB) *PostgresStore {
	return &PostgresStore{
		DB: db,
	}
}

func (s *PostgresStore) CreateUser(user *models.User) error {
	return s.DB.Create(user).Error
}

func (s *PostgresStore) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
