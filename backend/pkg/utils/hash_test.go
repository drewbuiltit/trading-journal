package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "securepassword"
	hashed, err := HashPassword(password)

	assert.NoError(t, err, "HashPassword should not return an error")
	assert.NotEmptyf(t, hashed, "Hashed password would not be empty")
	assert.NotEqual(t, password, hashed, "Hashed password should not be equal to the plain password")
}

func TestCheckPasswordHash(t *testing.T) {
	password := "securepassword"
	hashed, err := HashPassword(password)
	assert.NoError(t, err, "HashPassword should not return an error")

	isValid := CheckPasswordHash(password, hashed)
	assert.True(t, isValid, "CheckPasswordHash should return true for correct password")

	isValid = CheckPasswordHash("wrongpassword", hashed)
	assert.False(t, isValid, "CheckPasswordHash should return false for incorrect password")
}
