package usecase

import (
	"context"
	"testing"

	"gophkeeper/internal/models"
	"gophkeeper/internal/server/storage/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDeleteFileUsecase_DeleteFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDeleteFileRepository(ctrl)
	mockStorage := mocks.NewMockDeleteFileStorage(ctrl)

	uc := NewDeleteFileUsecase(mockRepo, mockStorage)

	ctx := context.Background()
	userID := uuid.New()
	label := "test-file"
	fileKey := "test-file-key"

	t.Run("successful delete", func(t *testing.T) {
		// Ожидаем вызов GetRecordByLabel в repo
		mockRepo.EXPECT().
			GetRecordByLabel(ctx, label, userID).
			Return(&models.GetSecretOut{
				Type:    models.SecretTypeFile,
				FileKey: fileKey,
			}, nil)

		// Ожидаем вызов DeleteFile в storage
		mockStorage.EXPECT().
			DeleteFile(ctx, fileKey).
			Return(nil)

		// Ожидаем вызов DeleteRecordByLabel в repo
		mockRepo.EXPECT().
			DeleteRecordByLabel(ctx, label, userID).
			Return(nil)

		err := uc.DeleteFile(ctx, label, userID)

		assert.NoError(t, err)
	})

	t.Run("record not found", func(t *testing.T) {
		// Ожидаем вызов GetRecordByLabel в repo с ошибкой
		mockRepo.EXPECT().
			GetRecordByLabel(ctx, label, userID).
			Return(nil, models.ErrRecordNotFound)

		err := uc.DeleteFile(ctx, label, userID)

		assert.Error(t, err)
		assert.Equal(t, models.ErrRecordNotFound, err)
	})

	t.Run("wrong record type", func(t *testing.T) {
		// Ожидаем вызов GetRecordByLabel в repo с неправильным типом
		mockRepo.EXPECT().
			GetRecordByLabel(ctx, label, userID).
			Return(&models.GetSecretOut{
				Type: models.SecretTypeText,
			}, nil)

		err := uc.DeleteFile(ctx, label, userID)

		assert.Error(t, err)
		assert.Equal(t, models.ErrRecordNotFound, err)
	})

	t.Run("empty file key", func(t *testing.T) {
		// Ожидаем вызов GetRecordByLabel в repo с пустым fileKey
		mockRepo.EXPECT().
			GetRecordByLabel(ctx, label, userID).
			Return(&models.GetSecretOut{
				Type:    models.SecretTypeFile,
				FileKey: "",
			}, nil)

		// Ожидаем вызов DeleteRecordByLabel в repo (без удаления файла)
		mockRepo.EXPECT().
			DeleteRecordByLabel(ctx, label, userID).
			Return(nil)

		err := uc.DeleteFile(ctx, label, userID)

		assert.NoError(t, err)
	})

	t.Run("storage delete error", func(t *testing.T) {
		// Ожидаем вызов GetRecordByLabel в repo
		mockRepo.EXPECT().
			GetRecordByLabel(ctx, label, userID).
			Return(&models.GetSecretOut{
				Type:    models.SecretTypeFile,
				FileKey: fileKey,
			}, nil)

		// Ожидаем вызов DeleteFile в storage с ошибкой
		mockStorage.EXPECT().
			DeleteFile(ctx, fileKey).
			Return(assert.AnError)

		// Ожидаем вызов DeleteRecordByLabel в repo (продолжаем несмотря на ошибку storage)
		mockRepo.EXPECT().
			DeleteRecordByLabel(ctx, label, userID).
			Return(nil)

		err := uc.DeleteFile(ctx, label, userID)

		assert.NoError(t, err)
	})

	t.Run("repository delete error", func(t *testing.T) {
		// Ожидаем вызов GetRecordByLabel в repo
		mockRepo.EXPECT().
			GetRecordByLabel(ctx, label, userID).
			Return(&models.GetSecretOut{
				Type:    models.SecretTypeFile,
				FileKey: fileKey,
			}, nil)

		// Ожидаем вызов DeleteFile в storage
		mockStorage.EXPECT().
			DeleteFile(ctx, fileKey).
			Return(nil)

		// Ожидаем вызов DeleteRecordByLabel в repo с ошибкой
		mockRepo.EXPECT().
			DeleteRecordByLabel(ctx, label, userID).
			Return(assert.AnError)

		err := uc.DeleteFile(ctx, label, userID)

		assert.Error(t, err)
	})
}
