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

func TestAuthenticateHandler_Authenticate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test cases
	tests := []struct {
		name           string
		method         string
		body           []byte
		ucResult       models.AuthenticateOut
		ucError        error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Successful authentication",
			method: http.MethodPost,
			body:   mustMarshal(t, models.AuthenticateIn{Login: "testuser", PasswordHash: "validhash"}),
			ucResult: models.AuthenticateOut{
				AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				RefreshToken: "refresh-token-123",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...","refresh_token":"refresh-token-123"}` + "\n",
		},
		{
			name:           "Invalid credentials",
			method:         http.MethodPost,
			body:           mustMarshal(t, models.AuthenticateIn{Login: "testuser", PasswordHash: "invalidhash"}),
			ucError:        models.ErrUserNotFound,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Login or password empty",
			method:         http.MethodPost,
			body:           mustMarshal(t, models.AuthenticateIn{Login: "testuser", PasswordHash: "validhash"}),
			ucError:        models.ErrLoginOrPasswordEmpty,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Usecase error",
			method:         http.MethodPost,
			body:           mustMarshal(t, models.AuthenticateIn{Login: "testuser", PasswordHash: "validhash"}),
			ucError:        assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"Failed to authenticate user","code":500}` + "\n",
		},
		{
			name:           "Wrong HTTP method",
			method:         http.MethodGet,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   `{"message":"Only POST requests are allowed!","code":405}` + "\n",
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
			mockUc := mocks.NewMockAuthenticateUsecase(ctrl)
			if tt.method == http.MethodPost && tt.name != "Invalid JSON" {
				mockUc.EXPECT().
					Authenticate(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(tt.ucResult, tt.ucError)
			}

			// Create handler
			handler := NewAuthenticateHandler(mockUc)

			// Create request
			var req *http.Request
			if tt.method == http.MethodPost {
				req = httptest.NewRequest(tt.method, "/api/authenticate", bytes.NewBuffer(tt.body))
			} else {
				req = httptest.NewRequest(tt.method, "/api/authenticate", nil)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.Authenticate(rr, req)

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
