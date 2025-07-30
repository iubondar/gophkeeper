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

func TestVersionHandler_GetSecretVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserID := uuid.New()

	tests := []struct {
		name           string
		method         string
		secretName     string
		withAuth       bool
		ucResult       *models.GetSecretOut
		ucError        error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Successful get version",
			method:         http.MethodGet,
			secretName:     "test-secret",
			withAuth:       true,
			ucResult:       &models.GetSecretOut{Version: 2},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Secret not found",
			method:         http.MethodGet,
			secretName:     "non-existent",
			withAuth:       true,
			ucError:        models.ErrRecordNotFound,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Unauthorized - no auth context",
			method:         http.MethodGet,
			secretName:     "test-secret",
			withAuth:       false,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Wrong HTTP method",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Missing secret name",
			method:         http.MethodGet,
			secretName:     "",
			withAuth:       true,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUc := mocks.NewMockGetSecretUsecase(ctrl)
			if tt.method == http.MethodGet && tt.withAuth && tt.secretName != "" {
				mockUc.EXPECT().
					GetSecret(gomock.Any(), tt.secretName, testUserID).
					Return(tt.ucResult, tt.ucError)
			}

			handler := NewVersionHandler(mockUc)

			req := httptest.NewRequest(tt.method, "/api/version?name="+tt.secretName, nil)

			if tt.withAuth {
				req = setUserIDInContext(req, testUserID)
			}

			rr := httptest.NewRecorder()
			handler.GetSecretVersion(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
