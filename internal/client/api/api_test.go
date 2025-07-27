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
