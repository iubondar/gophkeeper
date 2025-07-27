package usecase

import (
	"context"
	"gophkeeper/internal/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DeleteSecretRepository interface {
	GetRecordByLabel(ctx context.Context, label string, userID uuid.UUID) (*models.GetSecretOut, error)
	DeleteRecordByLabel(ctx context.Context, label string, userID uuid.UUID) error
}

type DeleteSecretStorage interface {
	DeleteFile(ctx context.Context, fileKey string) error
}

type DeleteSecretUsecase interface {
	DeleteSecret(ctx context.Context, secretName string, userID uuid.UUID) error
}

type deleteSecretUsecase struct {
	repo    DeleteSecretRepository
	storage DeleteSecretStorage
}

func NewDeleteSecretUsecase(repo DeleteSecretRepository, storage DeleteSecretStorage) DeleteSecretUsecase {
	return &deleteSecretUsecase{
		repo:    repo,
		storage: storage,
	}
}

func (uc *deleteSecretUsecase) DeleteSecret(ctx context.Context, secretName string, userID uuid.UUID) error {
	// Сначала получаем информацию о секрете, чтобы определить его тип
	record, err := uc.repo.GetRecordByLabel(ctx, secretName, userID)
	if err != nil {
		return err
	}

	if record.Type == models.SecretTypeFile {
		if record.FileKey != "" {
			if err := uc.storage.DeleteFile(ctx, record.FileKey); err != nil {
				zap.L().Sugar().Debugln("Failed to delete file from storage", zap.Error(err))
			}
		} else {
			zap.L().Sugar().Debugln("File key is empty")
		}
	}

	return uc.repo.DeleteRecordByLabel(ctx, secretName, userID)
}
