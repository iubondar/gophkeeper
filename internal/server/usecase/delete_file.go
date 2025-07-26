package usecase

import (
	"context"
	"gophkeeper/internal/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// DeleteFileStorage интерфейс для удаления файлов из хранилища
type DeleteFileStorage interface {
	DeleteFile(ctx context.Context, fileKey string) error
}

// DeleteFileRepository интерфейс для работы с записями файлов при удалении
type DeleteFileRepository interface {
	GetRecordByLabel(ctx context.Context, label string, userID uuid.UUID) (*models.GetSecretOut, error)
	DeleteRecordByLabel(ctx context.Context, label string, userID uuid.UUID) error
}

// DeleteFileUsecase интерфейс для удаления файлов
type DeleteFileUsecase interface {
	DeleteFile(ctx context.Context, label string, userID uuid.UUID) error
}

type deleteFileUsecase struct {
	repo    DeleteFileRepository
	storage DeleteFileStorage
}

// NewDeleteFileUsecase создает новый usecase для удаления файлов
func NewDeleteFileUsecase(repo DeleteFileRepository, storage DeleteFileStorage) DeleteFileUsecase {
	return &deleteFileUsecase{
		repo:    repo,
		storage: storage,
	}
}

// DeleteFile удаляет файл из хранилища и запись из БД
func (uc *deleteFileUsecase) DeleteFile(ctx context.Context, label string, userID uuid.UUID) error {
	// Получаем запись из БД
	record, err := uc.repo.GetRecordByLabel(ctx, label, userID)
	if err != nil {
		return err
	}

	// Проверяем, что это файл
	if record.Type != models.SecretTypeFile {
		return models.ErrRecordNotFound
	}

	// Удаляем файл из хранилища, если fileKey не пустой
	if record.FileKey != "" {
		if err := uc.storage.DeleteFile(ctx, record.FileKey); err != nil {
			// Логируем ошибку удаления файла, но продолжаем удаление записи
			zap.L().Sugar().Debugln("Failed to delete file from storage", zap.Error(err))
		}
	}

	// Удаляем запись из БД
	return uc.repo.DeleteRecordByLabel(ctx, label, userID)
}
