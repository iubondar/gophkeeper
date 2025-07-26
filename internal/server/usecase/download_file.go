package usecase

import (
	"context"
	"gophkeeper/internal/models"
	"io"

	"github.com/google/uuid"
)

// DownloadFileStorage интерфейс для скачивания файлов из хранилища
type DownloadFileStorage interface {
	DownloadFile(ctx context.Context, fileKey string) (io.ReadCloser, error)
}

// DownloadFileRepository интерфейс для работы с записями файлов при скачивании
type DownloadFileRepository interface {
	GetRecordByLabel(ctx context.Context, label string, userID uuid.UUID) (*models.GetSecretOut, error)
}

// DownloadFileUsecase интерфейс для скачивания файлов
type DownloadFileUsecase interface {
	DownloadFile(ctx context.Context, label string, userID uuid.UUID) (io.ReadCloser, error)
}

type downloadFileUsecase struct {
	repo    DownloadFileRepository
	storage DownloadFileStorage
}

// NewDownloadFileUsecase создает новый usecase для скачивания файлов
func NewDownloadFileUsecase(repo DownloadFileRepository, storage DownloadFileStorage) DownloadFileUsecase {
	return &downloadFileUsecase{
		repo:    repo,
		storage: storage,
	}
}

// DownloadFile скачивает файл из хранилища
func (uc *downloadFileUsecase) DownloadFile(ctx context.Context, label string, userID uuid.UUID) (io.ReadCloser, error) {
	// Получаем запись из БД
	record, err := uc.repo.GetRecordByLabel(ctx, label, userID)
	if err != nil {
		return nil, err
	}

	// Проверяем, что это файл
	if record.Type != models.SecretTypeFile {
		return nil, models.ErrRecordNotFound
	}

	// Проверяем, что fileKey не пустой
	if record.FileKey == "" {
		return nil, models.ErrRecordNotFound
	}

	// Скачиваем файл из хранилища
	reader, err := uc.storage.DownloadFile(ctx, record.FileKey)
	if err != nil {
		return nil, err
	}

	return reader, nil
}
