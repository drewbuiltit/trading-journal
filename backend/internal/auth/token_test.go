package auth

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestGenerateAndParseJWT(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "test_secret_key")
	defer os.Unsetenv("JWT_SECRET_KEY")

	Init()

	userID := 1
	token, err := GenerateJWT(userID)
	assert.NoError(t, err, "GenerateJWT should not return an error")
	assert.NotEmpty(t, token, "Generated JWT should not be empty")

	claims, err := ParseJWT(token)
	assert.NoError(t, err, "ParseJWT should not return an error for a valid token")
	assert.Equal(t, userID, claims.UserID, "UserID in claims should match")

	expirationTime := time.Unix(claims.ExpiresAt, 0).UTC()
	expectedExpiration := time.Now().UTC().Add(15 * time.Minute)

	assert.WithinDuration(t, expectedExpiration, expirationTime, time.Minute, "Token expiration should be around 15 minutes from generation")
}

func TestParseJWT_InvalidToken(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "test_secret_key")
	defer os.Unsetenv("JWT_SECRET_KEY")

	Init()

	invalidToken := "invalid.token.string"
	_, err := ParseJWT(invalidToken)
	assert.Error(t, err, "ParseJWT should return an error for an invalid token")
}

func TestGenerateAndParseRefreshToken(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "test_secret_key")
	defer os.Unsetenv("JWT_SECRET_KEY")

	Init()

	userID := 2
	refreshToken, err := GenerateRefreshToken(userID)
	assert.NoError(t, err, "GenerateRefreshToken should not return an error")
	assert.NotEmpty(t, userID, refreshToken, "Generated Refresh Token should not be empty")

	refreshClaims, err := ParseRefreshToken(refreshToken)
	assert.NoError(t, err, "ParseRefreshToken should not return an error for a valid refresh token")
	assert.Equal(t, userID, refreshClaims.UserID, "UserID in refresh claims should match")

	expirationTime := time.Unix(refreshClaims.ExpiresAt, 0).UTC()
	expectedExpiration := time.Now().UTC().Add(7 * 24 * time.Hour)
	assert.WithinDuration(t, expectedExpiration, expirationTime, 5*time.Minute, "Refresh token expiration should be around 7 days from generation")
}

func TestParseRefreshToken(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "test_secret_key")
	defer os.Unsetenv("JWT_SECRET_KEY")

	Init()

	invalidRefreshToken := "invalid.refresh.token"
	_, err := ParseRefreshToken(invalidRefreshToken)
	assert.Error(t, err, "ParseRefreshToken should return an error for an invalid refresh token")
}
