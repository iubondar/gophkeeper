package usecase

import (
	"context"
	"gophkeeper/internal/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// DeleteSecretRepository определяет интерфейс для удаления секретов из хранилища.
// Интерфейс используется для абстракции от конкретной реализации хранилища
// и позволяет тестировать usecase с помощью моков.
type DeleteSecretRepository interface {
	GetRecordByLabel(ctx context.Context, label string, userID uuid.UUID) (*models.GetSecretOut, error)
	DeleteRecordByLabel(ctx context.Context, label string, userID uuid.UUID) error
}

// DeleteSecretStorage определяет интерфейс для удаления файлов из файлового хранилища.
// Интерфейс используется для абстракции от конкретной реализации файлового хранилища
// и позволяет тестировать usecase с помощью моков.
type DeleteSecretStorage interface {
	DeleteFile(ctx context.Context, fileKey string) error
}

// DeleteSecretUsecase определяет интерфейс для удаления секретов.
// Интерфейс содержит бизнес-логику удаления секретов из хранилища,
// включая удаление связанных файлов для секретов типа "файл".
type DeleteSecretUsecase interface {
	DeleteSecret(ctx context.Context, secretName string, userID uuid.UUID) error
}

type deleteSecretUsecase struct {
	repo    DeleteSecretRepository
	storage DeleteSecretStorage
}

// NewDeleteSecretUsecase создает новый экземпляр DeleteSecretUsecase.
// Принимает репозиторий для удаления секретов и хранилище файлов.
// Функция используется для внедрения зависимостей и создания usecase
// с конкретными реализациями хранилищ.
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
