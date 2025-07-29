package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"gophkeeper/internal/models"
	"gophkeeper/internal/server/storage/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUploadHandler_Upload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем тестовый userID
	testUserID := uuid.New()

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
			name:           "Unauthorized - no auth context",
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
			if tt.method == http.MethodPost && tt.withAuth && tt.name != "Invalid JSON" && tt.name != "Unauthorized - no auth context" {
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

			// Добавляем userID в контекст если нужно
			if tt.withAuth {
				req = setUserIDInContext(req, testUserID)
			}

			rr := httptest.NewRecorder()
			handler.Upload(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
