package auth

import (
	"bytes"
	"encoding/json"
	"github.com/drewbuiltit/trading-journal/backend/internal/models"
	"github.com/drewbuiltit/trading-journal/backend/internal/store"
	"github.com/drewbuiltit/trading-journal/backend/pkg/utils"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func setupTestAuthHandler() *AuthHandler {
	return &AuthHandler{
		Store: store.NewMemoryStore(),
	}
}

func TestRegisterHandler(t *testing.T) {
	authHandler := setupTestAuthHandler()

	t.Run("Successful Registration", func(t *testing.T) {
		reqBody := RegisterRequest{
			Username: "jane_doe",
			Email:    "jane@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)

		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(authHandler.Register)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var user models.User
		err = json.Unmarshal(rr.Body.Bytes(), &user)
		assert.NoError(t, err)
		assert.Equal(t, reqBody.Username, user.Username)
		assert.Equal(t, reqBody.Email, user.Email)
		assert.Empty(t, user.Password, "Password should not be returned in the response")
	})

	t.Run("Registration with Existing Email", func(t *testing.T) {
		reqBody := RegisterRequest{
			Username: "jane_doe",
			Email:    "jane@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(authHandler.Register)
		handler.ServeHTTP(rr, req)

		req2, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
		rr2 := httptest.NewRecorder()
		handler.ServeHTTP(rr2, req2)

		assert.Equal(t, http.StatusBadRequest, rr2.Code)
		assert.Equal(t, "User already exists\n", rr2.Body.String())
	})

	t.Run("Registration with Invalid Payload", func(t *testing.T) {
		invalidJSON := `{"username": "jane_doe", "email": "jane@example.com", "password":}`

		req, err := http.NewRequest("POST", "/register", bytes.NewBufferString(invalidJSON))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(authHandler.Register)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "Invalid request payload\n", rr.Body.String())
	})
}

func TestLoginHandler(t *testing.T) {
	authHandler := setupTestAuthHandler()

	hashedPassword, _ := utils.HashPassword("password123")
	user := &models.User{
		Username: "jane_doe",
		Email:    "jane@example.com",
		Password: hashedPassword,
	}
	authHandler.Store.CreateUser(user)

	os.Setenv("JWT_SECRET_KEY", "test_secret_key")
	t.Cleanup(func() {
		os.Unsetenv("JWT_SECRET_KEY")
	})
	Init()

	t.Run("Successful Login", func(t *testing.T) {
		reqBody := LoginRequest{
			Email:    "jane@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(authHandler.Login)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var tokens TokenResponse
		err = json.Unmarshal(rr.Body.Bytes(), &tokens)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokens.AccessToken, "Access token should not be empty")
		assert.NotEmpty(t, tokens.RefreshToken, "Refresh token should not be empty")
	})

	t.Run("Login with Incorrect Password", func(t *testing.T) {
		reqBody := LoginRequest{
			Email:    "jane@example.com",
			Password: "wrongpassword",
		}
		body, _ := json.Marshal(reqBody)

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(authHandler.Login)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Equal(t, "Invalid email or password\n", rr.Body.String())
	})

	t.Run("Login with Non-Existent User", func(t *testing.T) {
		reqBody := LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(authHandler.Login)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Equal(t, "Invalid email or password\n", rr.Body.String())
	})

	t.Run("Login with Invalid Payload", func(t *testing.T) {
		invalidJSON := `{"email": "jane@example.com", "password":}`

		req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(invalidJSON))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(authHandler.Login)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "Invalid request payload\n", rr.Body.String())
	})
}

func TestRefreshTokenHandler(t *testing.T) {
	authHandler := setupTestAuthHandler()

	hashedPassword, _ := utils.HashPassword("password123")
	user := &models.User{
		Username: "jane_doe",
		Email:    "jane@example.com",
		Password: hashedPassword,
	}
	authHandler.Store.CreateUser(user)

	os.Setenv("JWT_SECRET_KEY", "test_secret_key")
	defer os.Unsetenv("JWT_SECRET_KEY")
	Init()

	refreshToken, err := GenerateRefreshToken(user.ID)
	assert.NoError(t, err)

	t.Run("Successful Token Refresh", func(t *testing.T) {
		reqBody := RefreshRequest{
			RefreshToken: refreshToken,
		}
		body, _ := json.Marshal(reqBody)

		req, err := http.NewRequest("POST", "/refresh", bytes.NewBuffer(body))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(authHandler.RefreshToken)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var tokens TokenResponse
		err = json.Unmarshal(rr.Body.Bytes(), &tokens)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokens.AccessToken, "New access token should not be empty")
		assert.NotEmpty(t, tokens.RefreshToken, "New refresh token should not be empty")

		newClaims, err := ParseJWT(tokens.AccessToken)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, newClaims.UserID, "UserID in new access claims should match")
	})

	t.Run("Refresh with Invalid Token", func(t *testing.T) {
		reqBody := RefreshRequest{
			RefreshToken: "invalid.token.string",
		}
		body, _ := json.Marshal(reqBody)

		req, err := http.NewRequest("POST", "/refresh", bytes.NewBuffer(body))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(authHandler.RefreshToken)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Equal(t, "Invalid refresh token\n", rr.Body.String())
	})

	t.Run("Refresh with Expired Token", func(t *testing.T) {
		expiredClaims := &RefreshClaims{
			UserID: user.ID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(),
				IssuedAt:  time.Now().Add(-2 * time.Hour).Unix(),
				Issuer:    "trading-journal-backend",
			},
		}
		expiredTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
		expiredRefreshToken, err := expiredTokenObj.SignedString(jwtKey)
		assert.NoError(t, err)

		reqBody := RefreshRequest{
			RefreshToken: expiredRefreshToken,
		}
		body, _ := json.Marshal(reqBody)

		req, err := http.NewRequest("POST", "/refresh", bytes.NewBuffer(body))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(authHandler.RefreshToken)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Equal(t, "Invalid refresh token\n", rr.Body.String())
	})

	t.Run("Refresh with Invalid Payload", func(t *testing.T) {
		invalidJSON := `{"refresh_token":}`

		req, err := http.NewRequest("POST", "/refresh", bytes.NewBufferString(invalidJSON))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(authHandler.RefreshToken)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "Invalid request payload\n", rr.Body.String())
	})
}
