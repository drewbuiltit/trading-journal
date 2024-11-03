package auth

import (
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestAuthMiddleWare(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "test_secret_key")
	defer os.Unsetenv("JWT_SECRET_KEY")
	Init()

	validToken, err := GenerateJWT(1)
	assert.NoError(t, err)

	invalidToken := "invalid.token.string"

	expiredClaims := &Claims{
		UserID: 1,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(),
			IssuedAt:  time.Now().Add(-2 * time.Hour).Unix(),
			Issuer:    "trading-journal-backend",
		},
	}
	expiredTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredToken, err := expiredTokenObj.SignedString(jwtKey)
	assert.NoError(t, err)

	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(UserContextKey).(int)
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User ID: " + strconv.Itoa(userID)))
	})

	handler := AuthMiddleWare(protectedHandler)

	t.Run("Access with Valid Token", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/protected", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+validToken)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "User ID: 1", rr.Body.String())
	})

	t.Run("Access with Invalid Token", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/protected", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+invalidToken)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Equal(t, "Invalid token\n", rr.Body.String())
	})

	t.Run("Access with Expired Token", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/protected", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+expiredToken)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Equal(t, "Invalid token\n", rr.Body.String())
	})

	t.Run("Access without Authorization Header", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/protected", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Equal(t, "Authorization header missing\n", rr.Body.String())
	})

	t.Run("Access with Malformed Authorization Header", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/protected", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Equal(t, "Invalid Authorization header format\n", rr.Body.String())
	})
}
