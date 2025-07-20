package api

import (
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
					Type:          "text",
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
			expectedBody:   `{"id":"test-id","label":"test-secret","type":"text","metadata":"test metadata","encrypted_data":"ZW5jcnlwdGVkLWRhdGE=","file_key":"","version":1}`,
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := mocks.NewMockGetSecretUsecase(ctrl)
			tt.setupMock(mockUC)

			handler := NewGetHandler(mockUC)

			req := httptest.NewRequest(tt.method, "/api/get?name="+tt.secretName, nil)
			if tt.userID != uuid.Nil {
				req = addUserIDToRequest(req, tt.userID)
			}

			w := httptest.NewRecorder()

			handler.GetSecret(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

// addUserIDToRequest добавляет userID в контекст запроса для тестирования
func addUserIDToRequest(req *http.Request, userID uuid.UUID) *http.Request {
	// Создаем тестовый JWT токен
	token, err := auth.GenerateAccessToken(userID.String())
	if err != nil {
		panic(err)
	}

	// Добавляем cookie с токеном
	req.AddCookie(&http.Cookie{
		Name:  auth.AuthCookieName,
		Value: token,
	})

	return req
}
