package usecase

import (
	"context"
	"errors"
	"gophkeeper/internal/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockDeleteSecretRepository struct {
	getFunc    func(ctx context.Context, label string, userID uuid.UUID) (*models.GetSecretOut, error)
	deleteFunc func(ctx context.Context, label string, userID uuid.UUID) error
}

func (m *mockDeleteSecretRepository) GetRecordByLabel(ctx context.Context, label string, userID uuid.UUID) (*models.GetSecretOut, error) {
	return m.getFunc(ctx, label, userID)
}

func (m *mockDeleteSecretRepository) DeleteRecordByLabel(ctx context.Context, label string, userID uuid.UUID) error {
	return m.deleteFunc(ctx, label, userID)
}

type mockDeleteSecretStorage struct {
	deleteFunc func(ctx context.Context, fileKey string) error
}

func (m *mockDeleteSecretStorage) DeleteFile(ctx context.Context, fileKey string) error {
	return m.deleteFunc(ctx, fileKey)
}

func TestDeleteSecretUsecase_DeleteSecret(t *testing.T) {
	userID := uuid.New()
	secretName := "test-secret"

	t.Run("success - text secret", func(t *testing.T) {
		mockRepo := &mockDeleteSecretRepository{
			getFunc: func(ctx context.Context, label string, uid uuid.UUID) (*models.GetSecretOut, error) {
				assert.Equal(t, secretName, label)
				assert.Equal(t, userID, uid)
				return &models.GetSecretOut{
					Type: models.SecretTypeText,
				}, nil
			},
			deleteFunc: func(ctx context.Context, label string, uid uuid.UUID) error {
				assert.Equal(t, secretName, label)
				assert.Equal(t, userID, uid)
				return nil
			},
		}
		mockStorage := &mockDeleteSecretStorage{}
		uc := NewDeleteSecretUsecase(mockRepo, mockStorage)
		err := uc.DeleteSecret(context.Background(), secretName, userID)
		assert.NoError(t, err)
	})

	t.Run("success - file secret", func(t *testing.T) {
		mockRepo := &mockDeleteSecretRepository{
			getFunc: func(ctx context.Context, label string, uid uuid.UUID) (*models.GetSecretOut, error) {
				assert.Equal(t, secretName, label)
				assert.Equal(t, userID, uid)
				return &models.GetSecretOut{
					Type:    models.SecretTypeFile,
					FileKey: "test-file-key",
				}, nil
			},
			deleteFunc: func(ctx context.Context, label string, uid uuid.UUID) error {
				assert.Equal(t, secretName, label)
				assert.Equal(t, userID, uid)
				return nil
			},
		}
		mockStorage := &mockDeleteSecretStorage{
			deleteFunc: func(ctx context.Context, fileKey string) error {
				assert.Equal(t, "test-file-key", fileKey)
				return nil
			},
		}
		uc := NewDeleteSecretUsecase(mockRepo, mockStorage)
		err := uc.DeleteSecret(context.Background(), secretName, userID)
		assert.NoError(t, err)
	})

	t.Run("get record error", func(t *testing.T) {
		mockRepo := &mockDeleteSecretRepository{
			getFunc: func(ctx context.Context, label string, uid uuid.UUID) (*models.GetSecretOut, error) {
				return nil, errors.New("get error")
			},
		}
		mockStorage := &mockDeleteSecretStorage{}
		uc := NewDeleteSecretUsecase(mockRepo, mockStorage)
		err := uc.DeleteSecret(context.Background(), secretName, userID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "get error")
	})

	t.Run("delete record error", func(t *testing.T) {
		mockRepo := &mockDeleteSecretRepository{
			getFunc: func(ctx context.Context, label string, uid uuid.UUID) (*models.GetSecretOut, error) {
				return &models.GetSecretOut{
					Type: models.SecretTypeText,
				}, nil
			},
			deleteFunc: func(ctx context.Context, label string, uid uuid.UUID) error {
				return errors.New("delete error")
			},
		}
		mockStorage := &mockDeleteSecretStorage{}
		uc := NewDeleteSecretUsecase(mockRepo, mockStorage)
		err := uc.DeleteSecret(context.Background(), secretName, userID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete error")
	})
}
