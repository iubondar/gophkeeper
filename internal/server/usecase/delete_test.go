package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockDeleteSecretRepository struct {
	deleteFunc func(ctx context.Context, label string, userID uuid.UUID) error
}

func (m *mockDeleteSecretRepository) DeleteRecordByLabel(ctx context.Context, label string, userID uuid.UUID) error {
	return m.deleteFunc(ctx, label, userID)
}

func TestDeleteSecretUsecase_DeleteSecret(t *testing.T) {
	userID := uuid.New()
	secretName := "test-secret"

	t.Run("success", func(t *testing.T) {
		mockRepo := &mockDeleteSecretRepository{
			deleteFunc: func(ctx context.Context, label string, uid uuid.UUID) error {
				assert.Equal(t, secretName, label)
				assert.Equal(t, userID, uid)
				return nil
			},
		}
		uc := NewDeleteSecretUsecase(mockRepo)
		err := uc.DeleteSecret(context.Background(), secretName, userID)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := &mockDeleteSecretRepository{
			deleteFunc: func(ctx context.Context, label string, uid uuid.UUID) error {
				return errors.New("not found")
			},
		}
		uc := NewDeleteSecretUsecase(mockRepo)
		err := uc.DeleteSecret(context.Background(), secretName, userID)
		assert.Error(t, err)
	})
}
