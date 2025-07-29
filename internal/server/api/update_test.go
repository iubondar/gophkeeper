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

func TestUpdateHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserID := uuid.New()

	tests := []struct {
		name           string
		method         string
		body           []byte
		withAuth       bool
		ucResult       models.UpdateSecretOut
		ucError        error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Successful update",
			method:         http.MethodPut,
			body:           mustMarshal(t, models.UpdateSecretIn{UploadSecretIn: models.UploadSecretIn{Label: "test", Type: "note", Metadata: "meta", EncryptedData: []byte("data"), FileKey: "key"}, Version: 1}),
			withAuth:       true,
			ucResult:       models.UpdateSecretOut{Version: 2},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Version conflict",
			method:         http.MethodPut,
			body:           mustMarshal(t, models.UpdateSecretIn{UploadSecretIn: models.UploadSecretIn{Label: "test", Type: "note"}, Version: 1}),
			withAuth:       true,
			ucError:        models.ErrConflict,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "Unauthorized - no auth context",
			method:         http.MethodPut,
			body:           mustMarshal(t, models.UpdateSecretIn{UploadSecretIn: models.UploadSecretIn{Label: "test", Type: "note"}, Version: 1}),
			withAuth:       false,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Wrong HTTP method",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Invalid JSON",
			method:         http.MethodPut,
			body:           []byte("invalid json"),
			withAuth:       true,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUc := mocks.NewMockUpdateSecretUsecase(ctrl)
			if tt.method == http.MethodPut && tt.withAuth && tt.name != "Invalid JSON" && tt.name != "Unauthorized - no auth context" {
				mockUc.EXPECT().
					UpdateSecret(gomock.Any(), gomock.Any(), testUserID).
					Return(tt.ucResult, tt.ucError)
			}

			handler := NewUpdateHandler(mockUc)

			var req *http.Request
			if tt.body != nil {
				req = httptest.NewRequest(tt.method, "/api/update", bytes.NewBuffer(tt.body))
			} else {
				req = httptest.NewRequest(tt.method, "/api/update", nil)
			}

			if tt.withAuth {
				req = setUserIDInContext(req, testUserID)
			}

			rr := httptest.NewRecorder()
			handler.Update(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
