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

func TestLoginHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test cases
	tests := []struct {
		name           string
		method         string
		body           []byte
		ucSalt         string
		ucError        error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Successful login request",
			method:         http.MethodPost,
			body:           mustMarshal(t, models.LoginIn{Login: "testuser"}),
			ucSalt:         "dGVzdC1zYWx0",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"salt":"dGVzdC1zYWx0"}` + "\n",
		},
		{
			name:           "User not found",
			method:         http.MethodPost,
			body:           mustMarshal(t, models.LoginIn{Login: "nonexistent"}),
			ucSalt:         "",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Usecase error",
			method:         http.MethodPost,
			body:           mustMarshal(t, models.LoginIn{Login: "testuser"}),
			ucError:        assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to get user salt\n",
		},
		{
			name:           "Wrong HTTP method",
			method:         http.MethodGet,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Only POST requests are allowed!\n",
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
			// Setup mock usecase
			mockUc := mocks.NewMockLoginUsecase(ctrl)
			if tt.method == http.MethodPost && tt.expectedStatus != http.StatusBadRequest {
				mockUc.EXPECT().
					GetSalt(gomock.Any(), gomock.Any()).
					Return(tt.ucSalt, tt.ucError)
			}

			// Create handler
			handler := NewLoginHandler(mockUc)

			// Create request
			var req *http.Request
			if tt.method == http.MethodPost {
				req = httptest.NewRequest(tt.method, "/api/login", bytes.NewBuffer(tt.body))
			} else {
				req = httptest.NewRequest(tt.method, "/api/login", nil)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.Login(rr, req)

			resp := rr.Result()
			defer resp.Body.Close()

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Check response body if expected
			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			}

			// Check content type for successful response
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
			}
		})
	}
}
