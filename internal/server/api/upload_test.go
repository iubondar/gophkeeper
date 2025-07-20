package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"gophkeeper/internal/auth"
	"gophkeeper/internal/models"
	"gophkeeper/internal/server/storage/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUploadHandler_Upload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем тестовый userID и токен
	testUserID := uuid.New()
	token, err := auth.GenerateAccessToken(testUserID.String())
	assert.NoError(t, err)

	tests := []struct {
		name           string
		method         string
		body           []byte
		withAuth       bool
		ucResult       models.UploadSecretOut
		ucError        error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Successful upload",
			method:         http.MethodPost,
			body:           mustMarshal(t, models.UploadSecretIn{Label: "test", Type: "note", Metadata: "meta", EncryptedData: []byte("data"), FileKey: "key"}),
			withAuth:       true,
			ucResult:       models.UploadSecretOut{ID: "id-123", Version: 1},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Usecase error",
			method:         http.MethodPost,
			body:           mustMarshal(t, models.UploadSecretIn{Label: "test", Type: "note"}),
			withAuth:       true,
			ucError:        assert.AnError,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Unauthorized - no auth cookie",
			method:         http.MethodPost,
			body:           mustMarshal(t, models.UploadSecretIn{Label: "test", Type: "note"}),
			withAuth:       false,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Wrong HTTP method",
			method:         http.MethodGet,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Invalid JSON",
			method:         http.MethodPost,
			body:           []byte("invalid json"),
			withAuth:       true,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUc := mocks.NewMockUploadSecretUsecase(ctrl)
			if tt.method == http.MethodPost && tt.withAuth && tt.name != "Invalid JSON" && tt.name != "Unauthorized - no auth cookie" {
				mockUc.EXPECT().
					UploadSecret(gomock.Any(), gomock.Any(), testUserID).
					Return(tt.ucResult, tt.ucError)
			}

			handler := NewUploadHandler(mockUc)

			var req *http.Request
			if tt.method == http.MethodPost {
				req = httptest.NewRequest(tt.method, "/api/upload", bytes.NewBuffer(tt.body))
			} else {
				req = httptest.NewRequest(tt.method, "/api/upload", nil)
			}

			// Добавляем аутентификацию если нужно
			if tt.withAuth {
				req.AddCookie(&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: token,
				})
			}

			rr := httptest.NewRecorder()
			handler.Upload(rr, req)

			resp := rr.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
