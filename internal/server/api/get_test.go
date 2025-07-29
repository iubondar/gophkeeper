package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gophkeeper/internal/models"
	"gophkeeper/internal/server/storage/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetHandler_GetSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		method         string
		secretName     string
		userID         uuid.UUID
		setupMock      func(*mocks.MockGetSecretUsecase)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:       "successful get secret",
			method:     http.MethodGet,
			secretName: "test-secret",
			userID:     uuid.New(),
			setupMock: func(mockUC *mocks.MockGetSecretUsecase) {
				expectedOut := &models.GetSecretOut{
					ID:            "test-id",
					Label:         "test-secret",
					Type:          models.SecretTypeText,
					Metadata:      "test metadata",
					EncryptedData: []byte("encrypted-data"),
					FileKey:       "",
					Version:       1,
				}
				mockUC.EXPECT().
					GetSecret(gomock.Any(), "test-secret", gomock.Any()).
					Return(expectedOut, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"test-id","label":"test-secret","type":"text","metadata":"test metadata","encrypted_data":"ZW5jcnlwdGVkLWRhdGE=","file_key":"","file_name":"","version":1}`,
		},
		{
			name:           "method not allowed",
			method:         http.MethodPost,
			secretName:     "test-secret",
			userID:         uuid.New(),
			setupMock:      func(mockUC *mocks.MockGetSecretUsecase) {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   `{"message":"Only GET requests are allowed!","code":405}`,
		},
		{
			name:           "missing secret name",
			method:         http.MethodGet,
			secretName:     "",
			userID:         uuid.New(),
			setupMock:      func(mockUC *mocks.MockGetSecretUsecase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Secret name is required","code":400}`,
		},
		{
			name:       "secret not found",
			method:     http.MethodGet,
			secretName: "non-existent-secret",
			userID:     uuid.New(),
			setupMock: func(mockUC *mocks.MockGetSecretUsecase) {
				mockUC.EXPECT().
					GetSecret(gomock.Any(), "non-existent-secret", gomock.Any()).
					Return(nil, models.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message":"Secret not found","code":404}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := mocks.NewMockGetSecretUsecase(ctrl)
			tt.setupMock(mockUC)

			handler := NewGetHandler(mockUC)

			req := httptest.NewRequest(tt.method, "/api/get?name="+tt.secretName, nil)
			if tt.userID != uuid.Nil {
				req = setUserIDInContext(req, tt.userID)
			}

			w := httptest.NewRecorder()

			handler.GetSecret(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
