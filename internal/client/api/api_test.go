package api

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"unicode/utf8"

	"gophkeeper/internal/models"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIClient_UploadFile(t *testing.T) {
	tests := []struct {
		name           string
		label          string
		metadata       string
		fileContent    string
		filename       string
		serverResponse string
		serverStatus   int
		expectedResult *models.UploadSecretOut
		expectedError  string
	}{
		{
			name:           "successful upload",
			label:          "test-file",
			metadata:       "test metadata",
			fileContent:    "test file content",
			filename:       "test.txt",
			serverResponse: `{"id":"file-123","version":1}`,
			serverStatus:   http.StatusOK,
			expectedResult: &models.UploadSecretOut{ID: "file-123", Version: 1},
		},
		{
			name:           "conflict error",
			label:          "existing-file",
			metadata:       "test metadata",
			fileContent:    "test file content",
			filename:       "test.txt",
			serverResponse: `{"message":"File with this label already exists"}`,
			serverStatus:   http.StatusConflict,
			expectedError:  "File with this label already exists",
		},
		{
			name:           "unauthorized",
			label:          "test-file",
			metadata:       "test metadata",
			fileContent:    "test file content",
			filename:       "test.txt",
			serverResponse: `{"message":"Unauthorized"}`,
			serverStatus:   http.StatusUnauthorized,
			expectedError:  "Unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/files", r.URL.Path)
				assert.Equal(t, http.MethodPost, r.Method)

				// Проверяем multipart форму
				err := r.ParseMultipartForm(32 << 20)
				require.NoError(t, err)

				assert.Equal(t, tt.label, r.FormValue("label"))
				assert.Equal(t, tt.metadata, r.FormValue("metadata"))

				file, header, err := r.FormFile("file")
				require.NoError(t, err)
				defer file.Close()

				assert.Equal(t, tt.filename, header.Filename)

				content, err := io.ReadAll(file)
				require.NoError(t, err)
				assert.Equal(t, tt.fileContent, string(content))

				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			client := NewAPIClient(server.URL)
			client.accessToken = "test-token"

			fileReader := strings.NewReader(tt.fileContent)
			result, err := client.UploadFile(context.Background(), tt.label, tt.metadata, fileReader, tt.filename)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestAPIClient_DownloadFile(t *testing.T) {
	tests := []struct {
		name           string
		label          string
		serverResponse string
		serverStatus   int
		expectedError  string
	}{
		{
			name:           "successful download",
			label:          "test-file",
			serverResponse: "test file content",
			serverStatus:   http.StatusOK,
		},
		{
			name:           "file not found",
			label:          "non-existent-file",
			serverResponse: `{"message":"File not found"}`,
			serverStatus:   http.StatusNotFound,
			expectedError:  "File not found",
		},
		{
			name:           "unauthorized",
			label:          "test-file",
			serverResponse: `{"message":"Unauthorized"}`,
			serverStatus:   http.StatusUnauthorized,
			expectedError:  "Unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := "/api/files/" + tt.label + "/download"
				assert.Equal(t, expectedPath, r.URL.Path)
				assert.Equal(t, http.MethodGet, r.Method)

				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			client := NewAPIClient(server.URL)
			client.accessToken = "test-token"

			reader, err := client.DownloadFile(context.Background(), tt.label)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, reader)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, reader)
				defer reader.Close()

				content, err := io.ReadAll(reader)
				require.NoError(t, err)
				assert.Equal(t, tt.serverResponse, string(content))
			}
		})
	}
}

func TestAPIClient_UploadFile_NoAuthToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"file-123","version":1}`))
	}))
	defer server.Close()

	client := NewAPIClient(server.URL)
	// Не устанавливаем accessToken

	fileReader := strings.NewReader("test content")
	result, err := client.UploadFile(context.Background(), "test-file", "metadata", fileReader, "test.txt")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access token is required")
	assert.Nil(t, result)
}

func TestAPIClient_DownloadFile_NoAuthToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test content"))
	}))
	defer server.Close()

	client := NewAPIClient(server.URL)
	// Не устанавливаем accessToken

	reader, err := client.DownloadFile(context.Background(), "test-file")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access token is required")
	assert.Nil(t, reader)
}

func TestAPIClient_DownloadFile_BinaryData(t *testing.T) {
	// Создаем бинарные данные (например, PNG заголовок)
	binaryData := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/files/test-image/download", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write(binaryData)
	}))
	defer server.Close()

	client := NewAPIClient(server.URL)
	client.accessToken = "test-token"

	reader, err := client.DownloadFile(context.Background(), "test-image")
	assert.NoError(t, err)
	require.NotNil(t, reader)
	defer reader.Close()

	content, err := io.ReadAll(reader)
	require.NoError(t, err)

	// Проверяем, что бинарные данные не были повреждены
	assert.Equal(t, binaryData, content)

	// Проверяем, что данные действительно бинарные (не UTF-8)
	assert.False(t, utf8.Valid(content))
}

func TestNewAPIClient(t *testing.T) {
	tests := []struct {
		name        string
		serverURL   string
		expectedURL string
	}{
		{
			name:        "with https protocol",
			serverURL:   "https://example.com",
			expectedURL: "https://example.com",
		},
		{
			name:        "with http protocol",
			serverURL:   "http://example.com",
			expectedURL: "http://example.com",
		},
		{
			name:        "without protocol",
			serverURL:   "example.com",
			expectedURL: "https://example.com",
		},
		{
			name:        "localhost without protocol",
			serverURL:   "localhost:8080",
			expectedURL: "https://localhost:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewAPIClient(tt.serverURL)
			assert.NotNil(t, client)
			assert.NotNil(t, client.httpc)
			assert.Equal(t, "", client.accessToken)
			assert.Equal(t, "", client.refreshToken)
		})
	}
}

func TestAPIClient_Register(t *testing.T) {
	tests := []struct {
		name           string
		input          models.RegisterIn
		serverResponse string
		serverStatus   int
		expectedError  string
	}{
		{
			name: "successful registration",
			input: models.RegisterIn{
				Login:        "testuser",
				PasswordHash: "hash123",
				Salt:         "salt123",
			},
			serverResponse: `{"access_token":"access123","refresh_token":"refresh123"}`,
			serverStatus:   http.StatusOK,
		},
		{
			name: "user already exists",
			input: models.RegisterIn{
				Login:        "existinguser",
				PasswordHash: "hash123",
				Salt:         "salt123",
			},
			serverResponse: `{"message":"User already exists"}`,
			serverStatus:   http.StatusConflict,
			expectedError:  "User already exists",
		},
		{
			name: "invalid request",
			input: models.RegisterIn{
				Login:        "",
				PasswordHash: "hash123",
				Salt:         "salt123",
			},
			serverResponse: `{"message":"Invalid request"}`,
			serverStatus:   http.StatusBadRequest,
			expectedError:  "Invalid request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/register", r.URL.Path)
				assert.Equal(t, http.MethodPost, r.Method)

				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			client := NewAPIClient(server.URL)
			err := client.Register(context.Background(), tt.input)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "access123", client.accessToken)
				assert.Equal(t, "refresh123", client.refreshToken)
			}
		})
	}
}

func TestAPIClient_Login(t *testing.T) {
	tests := []struct {
		name           string
		input          models.LoginIn
		serverResponse string
		serverStatus   int
		expectedSalt   string
		expectedError  string
	}{
		{
			name: "successful login",
			input: models.LoginIn{
				Login: "testuser",
			},
			serverResponse: `{"salt":"salt123"}`,
			serverStatus:   http.StatusOK,
			expectedSalt:   "salt123",
		},
		{
			name: "user not found",
			input: models.LoginIn{
				Login: "nonexistent",
			},
			serverResponse: `{"message":"User not found"}`,
			serverStatus:   http.StatusNotFound,
			expectedError:  "User not found",
		},
		{
			name: "invalid request",
			input: models.LoginIn{
				Login: "",
			},
			serverResponse: `{"message":"Invalid request"}`,
			serverStatus:   http.StatusBadRequest,
			expectedError:  "Invalid request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/login", r.URL.Path)
				assert.Equal(t, http.MethodPost, r.Method)

				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			client := NewAPIClient(server.URL)
			salt, err := client.Login(context.Background(), tt.input)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Equal(t, "", salt)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedSalt, salt)
			}
		})
	}
}

func TestAPIClient_Authenticate(t *testing.T) {
	tests := []struct {
		name           string
		input          models.AuthenticateIn
		serverResponse string
		serverStatus   int
		expectedError  string
	}{
		{
			name: "successful authentication",
			input: models.AuthenticateIn{
				Login:        "testuser",
				PasswordHash: "hash123",
			},
			serverResponse: `{"access_token":"access123","refresh_token":"refresh123"}`,
			serverStatus:   http.StatusOK,
		},
		{
			name: "invalid credentials",
			input: models.AuthenticateIn{
				Login:        "testuser",
				PasswordHash: "wronghash",
			},
			serverResponse: `{"message":"Invalid credentials"}`,
			serverStatus:   http.StatusUnauthorized,
			expectedError:  "Invalid credentials",
		},
		{
			name: "user not found",
			input: models.AuthenticateIn{
				Login:        "nonexistent",
				PasswordHash: "hash123",
			},
			serverResponse: `{"message":"User not found"}`,
			serverStatus:   http.StatusNotFound,
			expectedError:  "User not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/authenticate", r.URL.Path)
				assert.Equal(t, http.MethodPost, r.Method)

				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			client := NewAPIClient(server.URL)
			err := client.Authenticate(context.Background(), tt.input)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "access123", client.accessToken)
				assert.Equal(t, "refresh123", client.refreshToken)
			}
		})
	}
}

func TestAPIClient_UploadSecret(t *testing.T) {
	tests := []struct {
		name           string
		input          models.UploadSecretIn
		serverResponse string
		serverStatus   int
		expectedError  string
	}{
		{
			name: "successful upload",
			input: models.UploadSecretIn{
				Label:         "test-secret",
				Type:          "text",
				Metadata:      "test metadata",
				EncryptedData: []byte("encrypted data"),
			},
			serverResponse: `{"message":"Secret uploaded successfully"}`,
			serverStatus:   http.StatusOK,
		},
		{
			name: "unauthorized",
			input: models.UploadSecretIn{
				Label:         "test-secret",
				Type:          "text",
				Metadata:      "test metadata",
				EncryptedData: []byte("encrypted data"),
			},
			serverResponse: `{"message":"Unauthorized"}`,
			serverStatus:   http.StatusUnauthorized,
			expectedError:  "Unauthorized",
		},
		{
			name: "conflict",
			input: models.UploadSecretIn{
				Label:         "existing-secret",
				Type:          "text",
				Metadata:      "test metadata",
				EncryptedData: []byte("encrypted data"),
			},
			serverResponse: `{"message":"Secret already exists"}`,
			serverStatus:   http.StatusConflict,
			expectedError:  "Secret already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/upload", r.URL.Path)
				assert.Equal(t, http.MethodPost, r.Method)

				// Проверяем cookie
				cookies := r.Cookies()
				found := false
				for _, cookie := range cookies {
					if cookie.Name == "Authorization" && cookie.Value == "test-token" {
						found = true
						break
					}
				}
				assert.True(t, found, "Authorization cookie not found")

				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			client := NewAPIClient(server.URL)
			client.accessToken = "test-token"
			err := client.UploadSecret(context.Background(), tt.input)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAPIClient_GetSecretVersion(t *testing.T) {
	tests := []struct {
		name            string
		secretName      string
		serverResponse  string
		serverStatus    int
		expectedVersion int
		expectedError   string
	}{
		{
			name:            "successful get version",
			secretName:      "test-secret",
			serverResponse:  `{"version":5}`,
			serverStatus:    http.StatusOK,
			expectedVersion: 5,
		},
		{
			name:           "secret not found",
			secretName:     "nonexistent",
			serverResponse: `{"message":"Secret not found"}`,
			serverStatus:   http.StatusNotFound,
			expectedError:  "Secret not found",
		},
		{
			name:           "unauthorized",
			secretName:     "test-secret",
			serverResponse: `{"message":"Unauthorized"}`,
			serverStatus:   http.StatusUnauthorized,
			expectedError:  "Unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/version", r.URL.Path)
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, tt.secretName, r.URL.Query().Get("name"))

				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			client := NewAPIClient(server.URL)
			client.accessToken = "test-token"
			version, err := client.GetSecretVersion(context.Background(), tt.secretName)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Equal(t, 0, version)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedVersion, version)
			}
		})
	}
}

func TestAPIClient_UpdateSecret(t *testing.T) {
	tests := []struct {
		name           string
		input          models.UpdateSecretIn
		serverResponse string
		serverStatus   int
		expectedError  string
	}{
		{
			name: "successful update",
			input: models.UpdateSecretIn{
				UploadSecretIn: models.UploadSecretIn{
					Label:         "test-secret",
					Type:          "text",
					Metadata:      "updated metadata",
					EncryptedData: []byte("updated encrypted data"),
				},
				Version: 5,
			},
			serverResponse: `{"message":"Secret updated successfully"}`,
			serverStatus:   http.StatusOK,
		},
		{
			name: "version conflict",
			input: models.UpdateSecretIn{
				UploadSecretIn: models.UploadSecretIn{
					Label:         "test-secret",
					Type:          "text",
					Metadata:      "updated metadata",
					EncryptedData: []byte("updated encrypted data"),
				},
				Version: 3,
			},
			serverResponse: `{"message":"Version conflict"}`,
			serverStatus:   http.StatusConflict,
			expectedError:  "Version conflict",
		},
		{
			name: "secret not found",
			input: models.UpdateSecretIn{
				UploadSecretIn: models.UploadSecretIn{
					Label:         "nonexistent",
					Type:          "text",
					Metadata:      "updated metadata",
					EncryptedData: []byte("updated encrypted data"),
				},
				Version: 1,
			},
			serverResponse: `{"message":"Secret not found"}`,
			serverStatus:   http.StatusNotFound,
			expectedError:  "Secret not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/update", r.URL.Path)
				assert.Equal(t, http.MethodPut, r.Method)

				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			client := NewAPIClient(server.URL)
			client.accessToken = "test-token"
			err := client.UpdateSecret(context.Background(), tt.input)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAPIClient_GetSecret(t *testing.T) {
	tests := []struct {
		name           string
		secretName     string
		serverResponse string
		serverStatus   int
		expectedResult *models.GetSecretOut
		expectedError  string
	}{
		{
			name:       "successful get secret",
			secretName: "test-secret",
			serverResponse: `{
				"id": "secret-123",
				"label": "test-secret",
				"type": "text",
				"metadata": "test metadata",
				"encrypted_data": "ZW5jcnlwdGVkIGRhdGE=",
				"version": 5
			}`,
			serverStatus: http.StatusOK,
			expectedResult: &models.GetSecretOut{
				ID:            "secret-123",
				Label:         "test-secret",
				Type:          "text",
				Metadata:      "test metadata",
				EncryptedData: []byte("encrypted data"),
				Version:       5,
			},
		},
		{
			name:           "secret not found",
			secretName:     "nonexistent",
			serverResponse: `{"message":"Secret not found"}`,
			serverStatus:   http.StatusNotFound,
			expectedError:  "Secret not found",
		},
		{
			name:           "unauthorized",
			secretName:     "test-secret",
			serverResponse: `{"message":"Unauthorized"}`,
			serverStatus:   http.StatusUnauthorized,
			expectedError:  "Unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/get", r.URL.Path)
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, tt.secretName, r.URL.Query().Get("name"))

				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			client := NewAPIClient(server.URL)
			client.accessToken = "test-token"
			result, err := client.GetSecret(context.Background(), tt.secretName)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestAPIClient_DeleteSecret(t *testing.T) {
	tests := []struct {
		name           string
		secretName     string
		serverResponse string
		serverStatus   int
		expectedError  string
	}{
		{
			name:           "successful delete",
			secretName:     "test-secret",
			serverResponse: `{"message":"Secret deleted successfully"}`,
			serverStatus:   http.StatusOK,
		},
		{
			name:           "secret not found",
			secretName:     "nonexistent",
			serverResponse: `{"message":"Secret not found"}`,
			serverStatus:   http.StatusNotFound,
			expectedError:  "Secret not found",
		},
		{
			name:           "unauthorized",
			secretName:     "test-secret",
			serverResponse: `{"message":"Unauthorized"}`,
			serverStatus:   http.StatusUnauthorized,
			expectedError:  "Unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/delete", r.URL.Path)
				assert.Equal(t, http.MethodDelete, r.Method)
				assert.Equal(t, tt.secretName, r.URL.Query().Get("name"))

				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			client := NewAPIClient(server.URL)
			client.accessToken = "test-token"
			err := client.DeleteSecret(context.Background(), tt.secretName)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAPIClient_HealthCheck(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse string
		serverStatus   int
		expectedError  string
	}{
		{
			name:           "successful health check",
			serverResponse: `{"status":"ok"}`,
			serverStatus:   http.StatusOK,
		},
		{
			name:           "server error",
			serverResponse: `{"message":"Internal server error"}`,
			serverStatus:   http.StatusInternalServerError,
			expectedError:  "Internal server error",
		},
		{
			name:           "service unavailable",
			serverResponse: `{"message":"Service unavailable"}`,
			serverStatus:   http.StatusServiceUnavailable,
			expectedError:  "Service unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/health", r.URL.Path)
				assert.Equal(t, http.MethodGet, r.Method)

				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			client := NewAPIClient(server.URL)
			err := client.HealthCheck(context.Background())

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAPIClient_UploadSecret_NoAuthToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Success"}`))
	}))
	defer server.Close()

	client := NewAPIClient(server.URL)
	// Не устанавливаем accessToken

	input := models.UploadSecretIn{
		Label:         "test-secret",
		Type:          "text",
		Metadata:      "test metadata",
		EncryptedData: []byte("encrypted data"),
	}

	err := client.UploadSecret(context.Background(), input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access token is required")
}

func TestAPIClient_GetSecretVersion_NoAuthToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"version":1}`))
	}))
	defer server.Close()

	client := NewAPIClient(server.URL)
	// Не устанавливаем accessToken

	version, err := client.GetSecretVersion(context.Background(), "test-secret")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access token is required")
	assert.Equal(t, 0, version)
}

func TestAPIClient_UpdateSecret_NoAuthToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Success"}`))
	}))
	defer server.Close()

	client := NewAPIClient(server.URL)
	// Не устанавливаем accessToken

	input := models.UpdateSecretIn{
		UploadSecretIn: models.UploadSecretIn{
			Label:         "test-secret",
			Type:          "text",
			Metadata:      "test metadata",
			EncryptedData: []byte("encrypted data"),
		},
		Version: 1,
	}

	err := client.UpdateSecret(context.Background(), input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access token is required")
}

func TestAPIClient_GetSecret_NoAuthToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"secret-123","label":"test-secret","type":"text","metadata":"test","encrypted_data":"","version":1}`))
	}))
	defer server.Close()

	client := NewAPIClient(server.URL)
	// Не устанавливаем accessToken

	result, err := client.GetSecret(context.Background(), "test-secret")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access token is required")
	assert.Nil(t, result)
}

func TestAPIClient_DeleteSecret_NoAuthToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Success"}`))
	}))
	defer server.Close()

	client := NewAPIClient(server.URL)
	// Не устанавливаем accessToken

	err := client.DeleteSecret(context.Background(), "test-secret")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access token is required")
}

func TestAPIClient_handleAuthenticateResponse_InvalidJSON(t *testing.T) {
	client := &APIClient{}

	err := client.handleAuthenticateResponse([]byte("invalid json"))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal authenticate response")
}

func TestAPIClient_handleErrorResponse_NonJSONError(t *testing.T) {
	client := &APIClient{}

	// Создаем mock response с невалидным JSON
	response := &resty.Response{}
	// Симулируем ответ с ошибкой, но невалидным JSON
	response.RawResponse = &http.Response{
		StatusCode: 400,
		Body:       io.NopCloser(strings.NewReader("invalid json error")),
	}
	// Устанавливаем тело ответа
	response.SetBody([]byte("invalid json error"))

	err := client.handleErrorResponse(response)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid json error")
}

func TestAPIClient_handleErrorResponse_SuccessStatus(t *testing.T) {
	client := &APIClient{}

	// Создаем mock response с успешным статусом
	response := &resty.Response{}
	response.RawResponse = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("success")),
	}

	err := client.handleErrorResponse(response)

	assert.NoError(t, err)
}

func TestAPIClient_Register_NetworkError(t *testing.T) {
	client := NewAPIClient("http://invalid-server:9999")

	input := models.RegisterIn{
		Login:        "testuser",
		PasswordHash: "hash123",
		Salt:         "salt123",
	}

	err := client.Register(context.Background(), input)

	assert.Error(t, err)
}

func TestAPIClient_Login_NetworkError(t *testing.T) {
	client := NewAPIClient("http://invalid-server:9999")

	input := models.LoginIn{
		Login: "testuser",
	}

	salt, err := client.Login(context.Background(), input)

	assert.Error(t, err)
	assert.Equal(t, "", salt)
}

func TestAPIClient_Authenticate_NetworkError(t *testing.T) {
	client := NewAPIClient("http://invalid-server:9999")

	input := models.AuthenticateIn{
		Login:        "testuser",
		PasswordHash: "hash123",
	}

	err := client.Authenticate(context.Background(), input)

	assert.Error(t, err)
}

func TestAPIClient_UploadSecret_NetworkError(t *testing.T) {
	client := NewAPIClient("http://invalid-server:9999")
	client.accessToken = "test-token"

	input := models.UploadSecretIn{
		Label:         "test-secret",
		Type:          "text",
		Metadata:      "test metadata",
		EncryptedData: []byte("encrypted data"),
	}

	err := client.UploadSecret(context.Background(), input)

	assert.Error(t, err)
}

func TestAPIClient_GetSecretVersion_NetworkError(t *testing.T) {
	client := NewAPIClient("http://invalid-server:9999")
	client.accessToken = "test-token"

	version, err := client.GetSecretVersion(context.Background(), "test-secret")

	assert.Error(t, err)
	assert.Equal(t, 0, version)
}

func TestAPIClient_UpdateSecret_NetworkError(t *testing.T) {
	client := NewAPIClient("http://invalid-server:9999")
	client.accessToken = "test-token"

	input := models.UpdateSecretIn{
		UploadSecretIn: models.UploadSecretIn{
			Label:         "test-secret",
			Type:          "text",
			Metadata:      "test metadata",
			EncryptedData: []byte("encrypted data"),
		},
		Version: 1,
	}

	err := client.UpdateSecret(context.Background(), input)

	assert.Error(t, err)
}

func TestAPIClient_GetSecret_NetworkError(t *testing.T) {
	client := NewAPIClient("http://invalid-server:9999")
	client.accessToken = "test-token"

	result, err := client.GetSecret(context.Background(), "test-secret")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestAPIClient_DeleteSecret_NetworkError(t *testing.T) {
	client := NewAPIClient("http://invalid-server:9999")
	client.accessToken = "test-token"

	err := client.DeleteSecret(context.Background(), "test-secret")

	assert.Error(t, err)
}

func TestAPIClient_UploadFile_NetworkError(t *testing.T) {
	client := NewAPIClient("http://invalid-server:9999")
	client.accessToken = "test-token"

	fileReader := strings.NewReader("test content")
	result, err := client.UploadFile(context.Background(), "test-file", "metadata", fileReader, "test.txt")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestAPIClient_DownloadFile_NetworkError(t *testing.T) {
	client := NewAPIClient("http://invalid-server:9999")
	client.accessToken = "test-token"

	reader, err := client.DownloadFile(context.Background(), "test-file")

	assert.Error(t, err)
	assert.Nil(t, reader)
}

func TestAPIClient_HealthCheck_NetworkError(t *testing.T) {
	client := NewAPIClient("http://invalid-server:9999")

	err := client.HealthCheck(context.Background())

	assert.Error(t, err)
}

func TestAPIClient_GetSecret_InvalidJSONResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewAPIClient(server.URL)
	client.accessToken = "test-token"

	result, err := client.GetSecret(context.Background(), "test-secret")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal get secret response")
	assert.Nil(t, result)
}

func TestAPIClient_GetSecretVersion_InvalidJSONResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewAPIClient(server.URL)
	client.accessToken = "test-token"

	version, err := client.GetSecretVersion(context.Background(), "test-secret")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal version response")
	assert.Equal(t, 0, version)
}

func TestAPIClient_Login_InvalidJSONResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewAPIClient(server.URL)

	salt, err := client.Login(context.Background(), models.LoginIn{Login: "testuser"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal login response")
	assert.Equal(t, "", salt)
}

func TestAPIClient_UploadFile_InvalidJSONResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewAPIClient(server.URL)
	client.accessToken = "test-token"

	fileReader := strings.NewReader("test content")
	result, err := client.UploadFile(context.Background(), "test-file", "metadata", fileReader, "test.txt")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal upload file response")
	assert.Nil(t, result)
}

func TestAPIClient_Refresh(t *testing.T) {
	tests := []struct {
		name           string
		refreshToken   string
		serverResponse string
		serverStatus   int
		expectedError  string
	}{
		{
			name:           "successful refresh",
			refreshToken:   "valid-refresh-token",
			serverResponse: `{"access_token":"new-access-token","refresh_token":"new-refresh-token"}`,
			serverStatus:   http.StatusOK,
		},
		{
			name:           "invalid refresh token",
			refreshToken:   "invalid-refresh-token",
			serverResponse: `{"message":"invalid refresh token"}`,
			serverStatus:   http.StatusBadRequest,
			expectedError:  "invalid refresh token",
		},
		{
			name:           "expired refresh token",
			refreshToken:   "expired-refresh-token",
			serverResponse: `{"message":"refresh token expired"}`,
			serverStatus:   http.StatusUnauthorized,
			expectedError:  "refresh token expired",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/refresh", r.URL.Path)
				assert.Equal(t, http.MethodPost, r.Method)

				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			client := NewAPIClient(server.URL)
			client.refreshToken = tt.refreshToken

			err := client.Refresh(context.Background())

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "new-access-token", client.accessToken)
				assert.Equal(t, "new-refresh-token", client.refreshToken)
			}
		})
	}
}

func TestAPIClient_Refresh_NoRefreshToken(t *testing.T) {
	client := NewAPIClient("http://localhost:8080")

	err := client.Refresh(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "refresh token is required")
}

func TestAPIClient_Refresh_NetworkError(t *testing.T) {
	client := NewAPIClient("http://invalid-server:9999")
	client.refreshToken = "test-token"

	err := client.Refresh(context.Background())

	assert.Error(t, err)
}

func TestAPIClient_UploadSecret_TokenRefresh(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++

		if r.URL.Path == "/api/refresh" {
			// Первый вызов refresh - успешный
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"access_token":"new-access-token","refresh_token":"new-refresh-token"}`))
			return
		}

		if r.URL.Path == "/api/upload" {
			if callCount == 1 {
				// Первый вызов upload - 401 с expired token
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"message":"access token expired"}`))
			} else {
				// Второй вызов upload - успешный
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"message":"success"}`))
			}
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewAPIClient(server.URL)
	client.accessToken = "old-access-token"
	client.refreshToken = "old-refresh-token"

	err := client.UploadSecret(context.Background(), models.UploadSecretIn{
		Label:         "test-secret",
		Type:          "text",
		Metadata:      "test metadata",
		EncryptedData: []byte("test-data"),
	})

	assert.NoError(t, err)
	assert.Equal(t, "new-access-token", client.accessToken)
	assert.Equal(t, "new-refresh-token", client.refreshToken)
	assert.Equal(t, 3, callCount) // 1 upload + 1 refresh + 1 upload
}

func TestAPIClient_UploadSecret_TokenRefreshFailed(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++

		if r.URL.Path == "/api/refresh" {
			// Refresh неуспешен
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"message":"refresh token expired"}`))
			return
		}

		if r.URL.Path == "/api/upload" {
			// Первый вызов upload - 401 с expired token
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"message":"access token expired"}`))
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewAPIClient(server.URL)
	client.accessToken = "old-access-token"
	client.refreshToken = "old-refresh-token"

	err := client.UploadSecret(context.Background(), models.UploadSecretIn{
		Label:         "test-secret",
		Type:          "text",
		Metadata:      "test metadata",
		EncryptedData: []byte("test-data"),
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to refresh tokens")
	assert.Equal(t, 2, callCount) // 1 upload + 1 refresh
}

func TestAPIClient_GetSecret_TokenRefresh(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++

		if r.URL.Path == "/api/refresh" {
			// Первый вызов refresh - успешный
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"access_token":"new-access-token","refresh_token":"new-refresh-token"}`))
			return
		}

		if r.URL.Path == "/api/get" {
			if callCount == 1 {
				// Первый вызов get - 401 с expired token
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"message":"access token expired"}`))
			} else {
				// Второй вызов get - успешный
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"id":"secret-123","label":"test-secret","type":"text","metadata":"test metadata","encrypted_data":"dGVzdC1kYXRh","version":1}`))
			}
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewAPIClient(server.URL)
	client.accessToken = "old-access-token"
	client.refreshToken = "old-refresh-token"

	result, err := client.GetSecret(context.Background(), "test-secret")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test-secret", result.Label)
	assert.Equal(t, []byte("test-data"), result.EncryptedData)
	assert.Equal(t, 1, result.Version)
	assert.Equal(t, "new-access-token", client.accessToken)
	assert.Equal(t, "new-refresh-token", client.refreshToken)
	assert.Equal(t, 3, callCount) // 1 get + 1 refresh + 1 get
}

func TestAPIClient_UploadFile_TokenRefresh(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++

		if r.URL.Path == "/api/refresh" {
			// Первый вызов refresh - успешный
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"access_token":"new-access-token","refresh_token":"new-refresh-token"}`))
			return
		}

		if r.URL.Path == "/api/files" {
			if callCount == 1 {
				// Первый вызов upload - 401 с expired token
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"message":"access token expired"}`))
			} else {
				// Второй вызов upload - успешный
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"id":"file-123","version":1}`))
			}
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewAPIClient(server.URL)
	client.accessToken = "old-access-token"
	client.refreshToken = "old-refresh-token"

	fileReader := strings.NewReader("test file content")
	result, err := client.UploadFile(context.Background(), "test-file", "test metadata", fileReader, "test.txt")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "file-123", result.ID)
	assert.Equal(t, 1, result.Version)
	assert.Equal(t, "new-access-token", client.accessToken)
	assert.Equal(t, "new-refresh-token", client.refreshToken)
	assert.Equal(t, 3, callCount) // 1 upload + 1 refresh + 1 upload
}

func TestAPIClient_DownloadFile_TokenRefresh(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++

		if r.URL.Path == "/api/refresh" {
			// Первый вызов refresh - успешный
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"access_token":"new-access-token","refresh_token":"new-refresh-token"}`))
			return
		}

		if r.URL.Path == "/api/files/test-file/download" {
			if callCount == 1 {
				// Первый вызов download - 401 с expired token
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"message":"access token expired"}`))
			} else {
				// Второй вызов download - успешный
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("file content"))
			}
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewAPIClient(server.URL)
	client.accessToken = "old-access-token"
	client.refreshToken = "old-refresh-token"

	reader, err := client.DownloadFile(context.Background(), "test-file")

	assert.NoError(t, err)
	assert.NotNil(t, reader)

	content, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, "file content", string(content))

	assert.Equal(t, "new-access-token", client.accessToken)
	assert.Equal(t, "new-refresh-token", client.refreshToken)
	assert.Equal(t, 3, callCount) // 1 download + 1 refresh + 1 download
}
