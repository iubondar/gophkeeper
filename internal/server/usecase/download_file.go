package usecase

import (
	"context"
	"gophkeeper/internal/models"
	"io"

	"github.com/google/uuid"
)

// DownloadFileStorage определяет интерфейс для скачивания файлов из файлового хранилища.
// Интерфейс используется для абстракции от конкретной реализации файлового хранилища
// и позволяет тестировать usecase с помощью моков.
type DownloadFileStorage interface {
	DownloadFile(ctx context.Context, fileKey string) (io.ReadCloser, error)
}

// DownloadFileRepository определяет интерфейс для работы с записями файлов при скачивании.
// Интерфейс используется для абстракции от конкретной реализации хранилища
// и позволяет тестировать usecase с помощью моков.
type DownloadFileRepository interface {
	GetRecordByLabel(ctx context.Context, label string, userID uuid.UUID) (*models.GetSecretOut, error)
}

// DownloadFileUsecase определяет интерфейс для скачивания файлов.
// Интерфейс содержит бизнес-логику получения файлов из файлового хранилища
// на основе метки и идентификатора пользователя.
type DownloadFileUsecase interface {
	DownloadFile(ctx context.Context, label string, userID uuid.UUID) (io.ReadCloser, error)
}

type downloadFileUsecase struct {
	repo    DownloadFileRepository
	storage DownloadFileStorage
}

// NewDownloadFileUsecase создает новый экземпляр DownloadFileUsecase.
// Принимает репозиторий для работы с записями файлов и хранилище файлов.
// Функция используется для внедрения зависимостей и создания usecase
// с конкретными реализациями хранилищ.
func NewDownloadFileUsecase(repo DownloadFileRepository, storage DownloadFileStorage) DownloadFileUsecase {
	return &downloadFileUsecase{
		repo:    repo,
		storage: storage,
	}
}

// DownloadFile скачивает файл из хранилища.
// Функция получает запись файла из базы данных по метке и пользователю,
// проверяет тип записи и скачивает файл из файлового хранилища.
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
