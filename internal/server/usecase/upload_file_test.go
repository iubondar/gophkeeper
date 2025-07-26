package usecase

import (
	"context"
	"strings"
	"testing"

	"gophkeeper/internal/models"
	"gophkeeper/internal/server/storage/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUploadFileUsecase_UploadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUploadFileRepository(ctrl)
	mockStorage := mocks.NewMockUploadFileStorage(ctrl)

	uc := NewUploadFileUsecase(mockRepo, mockStorage)

	ctx := context.Background()
	userID := uuid.New()
	label := "test-file"
	metadata := "test metadata"
	content := "test content"
	reader := strings.NewReader(content)
	size := int64(len(content))

	t.Run("successful upload", func(t *testing.T) {
		// Ожидаем вызов UploadFile в storage
		mockStorage.EXPECT().
			UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), size).
			Return(nil)

		// Ожидаем вызов InsertRecord в repo
		mockRepo.EXPECT().
			InsertRecord(
				gomock.Any(),
				gomock.Any(),
				userID,
				label,
				models.SecretTypeFile,
				metadata,
				nil,
				gomock.Any(),
				1,
				gomock.Any(),
				gomock.Any(),
			).
			Return(nil)

		result, err := uc.UploadFile(ctx, label, metadata, reader, size, userID)

		require.NoError(t, err)
		assert.NotEmpty(t, result.ID)
		assert.Equal(t, 1, result.Version)
	})

	t.Run("storage error", func(t *testing.T) {
		// Ожидаем вызов UploadFile в storage с ошибкой
		mockStorage.EXPECT().
			UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), size).
			Return(assert.AnError)

		result, err := uc.UploadFile(ctx, label, metadata, reader, size, userID)

		assert.Error(t, err)
		assert.Equal(t, models.UploadSecretOut{}, result)
	})

	t.Run("database conflict error", func(t *testing.T) {
		// Ожидаем вызов UploadFile в storage
		mockStorage.EXPECT().
			UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), size).
			Return(nil)

		// Ожидаем вызов DeleteFile в storage (rollback)
		mockStorage.EXPECT().
			DeleteFile(gomock.Any(), gomock.Any()).
			Return(nil)

		// Ожидаем вызов InsertRecord в repo с конфликтом
		mockRepo.EXPECT().
			InsertRecord(
				gomock.Any(),
				gomock.Any(),
				userID,
				label,
				models.SecretTypeFile,
				metadata,
				nil,
				gomock.Any(),
				1,
				gomock.Any(),
				gomock.Any(),
			).
			Return(models.ErrConflict)

		result, err := uc.UploadFile(ctx, label, metadata, reader, size, userID)

		assert.Error(t, err)
		assert.Equal(t, models.ErrConflict, err)
		assert.Equal(t, models.UploadSecretOut{}, result)
	})
}
