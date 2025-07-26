package usecase

import (
	"context"
	"io"
	"strings"
	"testing"

	"gophkeeper/internal/models"
	"gophkeeper/internal/server/storage/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDownloadFileUsecase_DownloadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDownloadFileRepository(ctrl)
	mockStorage := mocks.NewMockDownloadFileStorage(ctrl)

	uc := NewDownloadFileUsecase(mockRepo, mockStorage)

	ctx := context.Background()
	userID := uuid.New()
	label := "test-file"
	fileKey := "test-file-key"
	content := "test content"

	t.Run("successful download", func(t *testing.T) {
		// Ожидаем вызов GetRecordByLabel в repo
		mockRepo.EXPECT().
			GetRecordByLabel(ctx, label, userID).
			Return(&models.GetSecretOut{
				Type:    models.SecretTypeFile,
				FileKey: fileKey,
			}, nil)

		// Ожидаем вызов DownloadFile в storage
		mockStorage.EXPECT().
			DownloadFile(ctx, fileKey).
			Return(io.NopCloser(strings.NewReader(content)), nil)

		reader, err := uc.DownloadFile(ctx, label, userID)

		require.NoError(t, err)
		require.NotNil(t, reader)
		defer reader.Close()

		// Читаем содержимое
		data, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Equal(t, content, string(data))
	})

	t.Run("record not found", func(t *testing.T) {
		// Ожидаем вызов GetRecordByLabel в repo с ошибкой
		mockRepo.EXPECT().
			GetRecordByLabel(ctx, label, userID).
			Return(nil, models.ErrRecordNotFound)

		reader, err := uc.DownloadFile(ctx, label, userID)

		assert.Error(t, err)
		assert.Equal(t, models.ErrRecordNotFound, err)
		assert.Nil(t, reader)
	})

	t.Run("wrong record type", func(t *testing.T) {
		// Ожидаем вызов GetRecordByLabel в repo с неправильным типом
		mockRepo.EXPECT().
			GetRecordByLabel(ctx, label, userID).
			Return(&models.GetSecretOut{
				Type: models.SecretTypeText,
			}, nil)

		reader, err := uc.DownloadFile(ctx, label, userID)

		assert.Error(t, err)
		assert.Equal(t, models.ErrRecordNotFound, err)
		assert.Nil(t, reader)
	})

	t.Run("empty file key", func(t *testing.T) {
		// Ожидаем вызов GetRecordByLabel в repo с пустым fileKey
		mockRepo.EXPECT().
			GetRecordByLabel(ctx, label, userID).
			Return(&models.GetSecretOut{
				Type:    models.SecretTypeFile,
				FileKey: "",
			}, nil)

		reader, err := uc.DownloadFile(ctx, label, userID)

		assert.Error(t, err)
		assert.Equal(t, models.ErrRecordNotFound, err)
		assert.Nil(t, reader)
	})

	t.Run("storage error", func(t *testing.T) {
		// Ожидаем вызов GetRecordByLabel в repo
		mockRepo.EXPECT().
			GetRecordByLabel(ctx, label, userID).
			Return(&models.GetSecretOut{
				Type:    models.SecretTypeFile,
				FileKey: fileKey,
			}, nil)

		// Ожидаем вызов DownloadFile в storage с ошибкой
		mockStorage.EXPECT().
			DownloadFile(ctx, fileKey).
			Return(nil, assert.AnError)

		reader, err := uc.DownloadFile(ctx, label, userID)

		assert.Error(t, err)
		assert.Nil(t, reader)
	})
}
