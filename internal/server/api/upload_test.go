package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"gophkeeper/internal/models"
	"gophkeeper/internal/server/storage/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUploadHandler_Upload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		method         string
		body           []byte
		ucResult       models.UploadSecretOut
		ucError        error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Successful upload",
			method:         http.MethodPost,
			body:           mustMarshal(t, models.UploadSecretIn{UserID: "b3b3b3b3-b3b3-b3b3-b3b3-b3b3b3b3b3b3", Label: "test", Type: "note", Metadata: "meta", EncryptedData: "data", FileKey: "key", Version: 1}),
			ucResult:       models.UploadSecretOut{ID: "id-123", Version: 1},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Usecase error",
			method:         http.MethodPost,
			body:           mustMarshal(t, models.UploadSecretIn{UserID: "b3b3b3b3-b3b3-b3b3-b3b3-b3b3b3b3b3b3", Label: "test", Type: "note"}),
			ucError:        assert.AnError,
			expectedStatus: http.StatusInternalServerError,
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
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUc := mocks.NewMockUploadSecretUsecase(ctrl)
			if tt.method == http.MethodPost && tt.name != "Invalid JSON" {
				mockUc.EXPECT().
					UploadSecret(gomock.Any(), gomock.Any()).
					Return(tt.ucResult, tt.ucError)
			}

			handler := NewUploadHandler(mockUc)

			var req *http.Request
			if tt.method == http.MethodPost {
				req = httptest.NewRequest(tt.method, "/api/upload", bytes.NewBuffer(tt.body))
			} else {
				req = httptest.NewRequest(tt.method, "/api/upload", nil)
			}

			rr := httptest.NewRecorder()
			handler.Upload(rr, req)

			resp := rr.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
