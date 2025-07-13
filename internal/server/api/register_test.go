package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gophkeeper/internal/models"
	"gophkeeper/internal/server/storage/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRegisterHandler_Register(t *testing.T) {
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
			name:   "Successful registration",
			method: http.MethodPost,
			body:   mustMarshal(t, models.RegisterIn{Login: "testuser", PasswordHash: "testpass", Salt: "testsalt"}),
			ucResult: models.AuthenticateOut{
				AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				RefreshToken: "refresh-token-123",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "User already exists",
			method:         http.MethodPost,
			body:           mustMarshal(t, models.RegisterIn{Login: "testuser", PasswordHash: "testpass", Salt: "testsalt"}),
			ucError:        models.ErrUserAlreadyExists,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "Login or password empty",
			method:         http.MethodPost,
			body:           mustMarshal(t, models.RegisterIn{Login: "testuser", PasswordHash: "testpass", Salt: "testsalt"}),
			ucError:        models.ErrLoginOrPasswordEmpty,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Usecase error",
			method:         http.MethodPost,
			body:           mustMarshal(t, models.RegisterIn{Login: "testuser", PasswordHash: "testpass", Salt: "testsalt"}),
			ucError:        assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to register user\n",
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
			mockUc := mocks.NewMockRegisterUsecase(ctrl)
			if tt.method == http.MethodPost && tt.name != "Invalid JSON" {
				mockUc.EXPECT().
					Register(gomock.Any(), gomock.Any()).
					Return(tt.ucResult, tt.ucError)
			}

			// Create handler
			handler := NewRegisterHandler(mockUc)

			// Create request
			var req *http.Request
			if tt.method == http.MethodPost {
				req = httptest.NewRequest(tt.method, "/api/user/register", bytes.NewBuffer(tt.body))
			} else {
				req = httptest.NewRequest(tt.method, "/api/user/register", nil)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.Register(rr, req)

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

func mustMarshal(t *testing.T, v any) []byte {
	t.Helper()
	data, err := json.Marshal(v)
	assert.NoError(t, err)
	return data
}
