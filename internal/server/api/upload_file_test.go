package api

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"gophkeeper/internal/models"
	"gophkeeper/internal/server/storage/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUploadFileHandler_UploadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUploadFileUsecase(ctrl)
	handler := NewUploadFileHandler(mockUsecase)

	userID := uuid.New()
	label := "test-file"
	metadata := "test metadata"
	content := "test content"

	t.Run("successful upload", func(t *testing.T) {
		// Создаем multipart форму
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Добавляем поля формы
		writer.WriteField("label", label)
		writer.WriteField("metadata", metadata)

		// Добавляем файл
		part, err := writer.CreateFormFile("file", "test.txt")
		require.NoError(t, err)
		part.Write([]byte(content))
		writer.Close()

		// Создаем запрос
		req := httptest.NewRequest(http.MethodPost, "/api/files", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Добавляем userID в контекст
		req = setUserIDInContext(req, userID)

		w := httptest.NewRecorder()

		// Ожидаем вызов usecase
		expectedResult := models.UploadSecretOut{
			ID:      uuid.New().String(),
			Version: 1,
		}
		mockUsecase.EXPECT().
			UploadFile(gomock.Any(), label, metadata, "test.txt", gomock.Any(), int64(len(content)), userID).
			Return(expectedResult, nil)

		// Выполняем запрос
		handler.UploadFile(w, req)

		// Проверяем результат
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response models.UploadSecretOut
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, expectedResult.ID, response.ID)
		assert.Equal(t, expectedResult.Version, response.Version)
	})

	t.Run("wrong HTTP method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/files", nil)
		w := httptest.NewRecorder()

		handler.UploadFile(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("missing user ID", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("label", label)
		writer.WriteField("metadata", metadata)
		part, err := writer.CreateFormFile("file", "test.txt")
		require.NoError(t, err)
		part.Write([]byte(content))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/files", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		// Не добавляем userID в контекст

		w := httptest.NewRecorder()

		handler.UploadFile(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("missing label", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("metadata", metadata)
		part, err := writer.CreateFormFile("file", "test.txt")
		require.NoError(t, err)
		part.Write([]byte(content))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/files", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req = setUserIDInContext(req, userID)

		w := httptest.NewRecorder()

		handler.UploadFile(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing file", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("label", label)
		writer.WriteField("metadata", metadata)
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/files", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req = setUserIDInContext(req, userID)

		w := httptest.NewRecorder()

		handler.UploadFile(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase conflict error", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("label", label)
		writer.WriteField("metadata", metadata)
		part, err := writer.CreateFormFile("file", "test.txt")
		require.NoError(t, err)
		part.Write([]byte(content))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/files", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req = setUserIDInContext(req, userID)

		w := httptest.NewRecorder()

		mockUsecase.EXPECT().
			UploadFile(gomock.Any(), label, metadata, "test.txt", gomock.Any(), int64(len(content)), userID).
			Return(models.UploadSecretOut{}, models.ErrConflict)

		handler.UploadFile(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("label", label)
		writer.WriteField("metadata", metadata)
		part, err := writer.CreateFormFile("file", "test.txt")
		require.NoError(t, err)
		part.Write([]byte(content))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/files", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req = setUserIDInContext(req, userID)

		w := httptest.NewRecorder()

		mockUsecase.EXPECT().
			UploadFile(gomock.Any(), label, metadata, "test.txt", gomock.Any(), int64(len(content)), userID).
			Return(models.UploadSecretOut{}, assert.AnError)

		handler.UploadFile(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
