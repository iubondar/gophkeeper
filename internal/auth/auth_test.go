package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func init() {
	// Initialize logger for tests
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}

func TestGenerateAccessToken(t *testing.T) {
	// Create a test user ID
	userID := uuid.New().String()

	// Call the function
	tokenString, err := GenerateAccessToken(userID)
	require.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Parse the token to verify its contents
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	require.NoError(t, err)
	assert.True(t, token.Valid)

	// Verify claims
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	require.True(t, ok)
	assert.Equal(t, userID, claims.Subject)
}

func TestGenerateRefreshToken(t *testing.T) {
	// Create a test user ID
	userID := uuid.New().String()

	// Call the function
	tokenString, err := GenerateRefreshToken(userID)
	require.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Parse the token to verify its contents
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	require.NoError(t, err)
	assert.True(t, token.Valid)

	// Verify claims
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	require.True(t, ok)
	assert.Equal(t, userID, claims.Subject)
}

func TestGetUserIDFromReq(t *testing.T) {
	// Create a test user ID
	userID := uuid.New()

	// Create a valid JWT token
	tokenString, err := GenerateAccessToken(userID.String())
	require.NoError(t, err)

	// Create a request with the cookie
	req := httptest.NewRequest("GET", "/", nil)
	cookie := &http.Cookie{
		Name:  AuthCookieName,
		Value: tokenString,
	}
	req.AddCookie(cookie)

	// Test valid cookie
	extractedUserID, err := GetUserIDFromReq(req)
	require.NoError(t, err)
	assert.Equal(t, userID, extractedUserID)

	// Test missing cookie
	req = httptest.NewRequest("GET", "/", nil)
	_, err = GetUserIDFromReq(req)
	assert.Error(t, err)
}

func TestGetUserID(t *testing.T) {
	// Create a test user ID
	userID := uuid.New()

	// Test valid token
	tokenString, err := GenerateAccessToken(userID.String())
	require.NoError(t, err)

	extractedUserID, err := getUserID(tokenString)
	require.NoError(t, err)
	assert.Equal(t, userID, extractedUserID)

	// Test invalid token
	_, err = getUserID("invalid.token.string")
	assert.Error(t, err)

	// Test expired token
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
	})
	expiredTokenString, err := expiredToken.SignedString([]byte(secretKey))
	require.NoError(t, err)
	_, err = getUserID(expiredTokenString)
	assert.Error(t, err)
}
