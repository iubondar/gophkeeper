package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gophkeeper/internal/auth"
	"gophkeeper/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefreshHandler_Refresh(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		requestBody    interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "wrong method",
			method:         http.MethodGet,
			requestBody:    models.RefreshIn{RefreshToken: "test"},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "Only POST requests are allowed!",
		},
		{
			name:           "empty request body",
			method:         http.MethodPost,
			requestBody:    "",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "unexpected end of JSON input",
		},
		{
			name:           "invalid JSON",
			method:         http.MethodPost,
			requestBody:    `{"invalid": json}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid character 'j' looking for beginning of value",
		},
		{
			name:           "empty refresh token",
			method:         http.MethodPost,
			requestBody:    models.RefreshIn{RefreshToken: ""},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid refresh token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем mock usecase
			mockUsecase := &mockRefreshUsecase{
				refreshFunc: func(ctx context.Context, token string) (models.RefreshOut, error) {
					if token == "" {
						return models.RefreshOut{}, models.ErrRefreshTokenInvalid
					}
					return models.RefreshOut{}, errors.New("unknown error")
				},
			}
			handler := NewRefreshHandler(mockUsecase)

			// Создаем request body
			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				require.NoError(t, err)
			}

			// Создаем HTTP request
			req := httptest.NewRequest(tt.method, "/api/refresh", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			// Вызываем handler
			handler.Refresh(w, req)

			// Проверяем статус код
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Проверяем ошибку
			if tt.expectedError != "" {
				var jsonError models.JSONError
				err = json.Unmarshal(w.Body.Bytes(), &jsonError)
				assert.NoError(t, err)
				assert.Contains(t, jsonError.Message, tt.expectedError)
			}
		})
	}
}

func TestRefreshHandler_Refresh_Success(t *testing.T) {
	// Создаем валидный refresh token
	userID := uuid.New()
	refreshToken, err := auth.GenerateRefreshToken(userID.String())
	require.NoError(t, err)

	// Создаем mock usecase
	mockUsecase := &mockRefreshUsecase{
		refreshFunc: func(ctx context.Context, token string) (models.RefreshOut, error) {
			return models.RefreshOut{
				AccessToken:  "new_access_token",
				RefreshToken: "new_refresh_token",
			}, nil
		},
	}

	handler := NewRefreshHandler(mockUsecase)

	// Создаем request body
	requestBody := models.RefreshIn{RefreshToken: refreshToken}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	// Создаем HTTP request
	req := httptest.NewRequest(http.MethodPost, "/api/refresh", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Вызываем handler
	handler.Refresh(w, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем response
	var response models.RefreshOut
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "new_access_token", response.AccessToken)
	assert.Equal(t, "new_refresh_token", response.RefreshToken)
}

func TestRefreshHandler_Refresh_UsecaseError(t *testing.T) {
	tests := []struct {
		name           string
		usecaseError   error
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "refresh token invalid error",
			usecaseError:   models.ErrRefreshTokenInvalid,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid refresh token",
		},
		{
			name:           "refresh token expired error",
			usecaseError:   models.ErrRefreshTokenExpired,
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "refresh token expired",
		},
		{
			name:           "unknown error",
			usecaseError:   errors.New("unknown error"),
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to refresh tokens",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем mock usecase
			mockUsecase := &mockRefreshUsecase{
				refreshFunc: func(ctx context.Context, token string) (models.RefreshOut, error) {
					return models.RefreshOut{}, tt.usecaseError
				},
			}

			handler := NewRefreshHandler(mockUsecase)

			// Создаем request body
			requestBody := models.RefreshIn{RefreshToken: "test_token"}
			body, err := json.Marshal(requestBody)
			require.NoError(t, err)

			// Создаем HTTP request
			req := httptest.NewRequest(http.MethodPost, "/api/refresh", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			// Вызываем handler
			handler.Refresh(w, req)

			// Проверяем статус код
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Проверяем ошибку
			var jsonError models.JSONError
			err = json.Unmarshal(w.Body.Bytes(), &jsonError)
			assert.NoError(t, err)
			assert.Contains(t, jsonError.Message, tt.expectedError)
		})
	}
}

// mockRefreshUsecase - мок для RefreshUsecase
type mockRefreshUsecase struct {
	refreshFunc func(context.Context, string) (models.RefreshOut, error)
}

func (m *mockRefreshUsecase) Refresh(ctx context.Context, refreshToken string) (models.RefreshOut, error) {
	if m.refreshFunc != nil {
		return m.refreshFunc(ctx, refreshToken)
	}
	return models.RefreshOut{}, errors.New("mock not configured")
}
