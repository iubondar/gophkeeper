package usecase

import (
	"context"
	"gophkeeper/internal/models"
	"io"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UploadFileStorage определяет интерфейс для загрузки файлов в файловое хранилище.
// Интерфейс используется для абстракции от конкретной реализации файлового хранилища
// и позволяет тестировать usecase с помощью моков.
type UploadFileStorage interface {
	UploadFile(ctx context.Context, fileKey string, reader io.Reader, size int64) error
	DeleteFile(ctx context.Context, fileKey string) error
}

// UploadFileRepository определяет интерфейс для работы с записями файлов при загрузке.
// Интерфейс используется для абстракции от конкретной реализации хранилища
// и позволяет тестировать usecase с помощью моков.
type UploadFileRepository interface {
	InsertRecord(ctx context.Context, id, userID uuid.UUID, label, recordType, metadata string, encryptedData []byte, fileKey string, fileName string, version int, createdAt, updatedAt time.Time) error
}

// UploadFileUsecase определяет интерфейс для загрузки файлов.
// Интерфейс содержит бизнес-логику загрузки файлов в файловое хранилище
// и создания соответствующих записей в базе данных.
type UploadFileUsecase interface {
	UploadFile(ctx context.Context, label, metadata, fileName string, reader io.Reader, size int64, userID uuid.UUID) (models.UploadSecretOut, error)
}

type uploadFileUsecase struct {
	repo    UploadFileRepository
	storage UploadFileStorage
}

// NewUploadFileUsecase создает новый экземпляр UploadFileUsecase.
// Принимает репозиторий для работы с записями файлов и хранилище файлов.
// Функция используется для внедрения зависимостей и создания usecase
// с конкретными реализациями хранилищ.
func NewUploadFileUsecase(repo UploadFileRepository, storage UploadFileStorage) UploadFileUsecase {
	return &uploadFileUsecase{
		repo:    repo,
		storage: storage,
	}
}

// UploadFile загружает файл в хранилище и создает запись в БД.
// Функция генерирует уникальный ключ для файла, загружает его в файловое хранилище
// и создает соответствующую запись в базе данных с типом "файл".
func (uc *uploadFileUsecase) UploadFile(ctx context.Context, label, metadata, fileName string, reader io.Reader, size int64, userID uuid.UUID) (models.UploadSecretOut, error) {
	// Генерируем уникальный ключ для файла
	fileKey := uuid.New().String()

	// Загружаем файл в MinIO
	err := uc.storage.UploadFile(ctx, fileKey, reader, size)
	if err != nil {
		return models.UploadSecretOut{}, err
	}

	// Создаем запись в БД
	id := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now()
	version := 1

	err = uc.repo.InsertRecord(ctx, id, userID, label, models.SecretTypeFile, metadata, nil, fileKey, fileName, version, createdAt, updatedAt)
	if err != nil {
		// Если не удалось создать запись в БД, удаляем файл из хранилища
		if deleteErr := uc.storage.DeleteFile(ctx, fileKey); deleteErr != nil {
			// Логируем ошибку удаления, но возвращаем основную ошибку
			zap.L().Sugar().Debugln("Failed to delete file from storage", zap.Error(deleteErr))
		}

		if err == models.ErrConflict {
			return models.UploadSecretOut{}, models.ErrConflict
		}
		return models.UploadSecretOut{}, err
	}

	return models.UploadSecretOut{ID: id.String(), Version: version}, nil
}
