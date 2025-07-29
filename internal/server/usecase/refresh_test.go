package usecase

import (
	"context"
	"testing"

	"gophkeeper/internal/auth"
	"gophkeeper/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefreshUsecase_Refresh(t *testing.T) {
	uc := NewRefreshUsecase()

	tests := []struct {
		name         string
		refreshToken string
		wantErr      error
		wantResult   models.RefreshOut
	}{
		{
			name:         "empty refresh token",
			refreshToken: "",
			wantErr:      models.ErrRefreshTokenInvalid,
		},
		{
			name:         "invalid refresh token",
			refreshToken: "invalid_token",
			wantErr:      models.ErrRefreshTokenExpired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := uc.Refresh(context.Background(), tt.refreshToken)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, models.RefreshOut{}, result)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
			}
		})
	}
}

func TestRefreshUsecase_Refresh_ValidToken(t *testing.T) {
	uc := NewRefreshUsecase()

	// Создаем валидный refresh token
	userID := uuid.New()
	refreshToken, err := auth.GenerateRefreshToken(userID.String())
	require.NoError(t, err)

	result, err := uc.Refresh(context.Background(), refreshToken)

	assert.NoError(t, err)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)

	// Проверяем, что новые токены валидны
	validatedUserID, err := auth.ValidateRefreshToken(result.RefreshToken)
	assert.NoError(t, err)
	assert.Equal(t, userID, validatedUserID)
}

func TestMakeRefreshOut(t *testing.T) {
	userID := uuid.New()

	result, err := makeRefreshOut(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)

	// Проверяем, что токены валидны
	validatedUserID, err := auth.ValidateRefreshToken(result.RefreshToken)
	assert.NoError(t, err)
	assert.Equal(t, userID, validatedUserID)

	// Проверяем access token через GetUserIDFromReq (использует ту же логику)
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(result.AccessToken, claims,
		func(t *jwt.Token) (any, error) {
			return []byte("supersecretkey"), nil
		})
	assert.NoError(t, err)
	assert.True(t, token.Valid)

	parsedUserID, err := uuid.Parse(claims.Subject)
	assert.NoError(t, err)
	assert.Equal(t, userID, parsedUserID)
}
