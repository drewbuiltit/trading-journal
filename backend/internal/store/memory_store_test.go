package store

import (
	"github.com/drewbuiltit/trading-journal/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoryStore_CreateUser(t *testing.T) {
	store := NewMemoryStore()

	user := &models.User{
		Username: "john_doe",
		Email:    "john@example.com",
		Password: "hashedpassword",
	}

	err := store.CreateUser(user)
	assert.NoError(t, err, "CreateUser should not return an error for a new user")

	err = store.CreateUser(user)
	assert.Error(t, err, "CreateUser should return an error for an existing user")
}

func TestMemoryStore_GetUserByEmail(t *testing.T) {
	store := NewMemoryStore()

	user := &models.User{
		Username: "john_doe",
		Email:    "john@example.com",
		Password: "hashedpassword",
	}

	_, err := store.GetUserByEmail(user.Email)
	assert.Error(t, err, "GetUserByEmail should return an error for non-existent user")

	err = store.CreateUser(user)
	assert.NoError(t, err, "CreateUser should not return an error for a new user")

	retrievedUser, err := store.GetUserByEmail(user.Email)
	assert.NoError(t, err, "GetUserByEmail should not return an error for an existing user")
	assert.Equal(t, user.Username, retrievedUser.Username, "Usernames should match")
	assert.Equal(t, user.Email, retrievedUser.Email, "Emails should match")
	assert.Equal(t, user.Password, retrievedUser.Password, "Passwords should match")
}
