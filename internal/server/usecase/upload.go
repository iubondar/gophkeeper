package usecase

import (
	"context"
	"gophkeeper/internal/models"
	"time"

	"github.com/google/uuid"
)

// RecordRepository определяет интерфейс для работы с записями секретов в хранилище.
// Интерфейс используется для абстракции от конкретной реализации хранилища
// и позволяет тестировать usecase с помощью моков.
type RecordRepository interface {
	InsertRecord(ctx context.Context, id, userID uuid.UUID, label, recordType, metadata string, encryptedData []byte, fileKey string, fileName string, version int, createdAt, updatedAt time.Time) error
}

// UploadSecretUsecase определяет интерфейс для загрузки секретов.
// Интерфейс содержит бизнес-логику сохранения зашифрованных секретов
// в хранилище с генерацией уникального ID и версии.
type UploadSecretUsecase interface {
	UploadSecret(ctx context.Context, in models.UploadSecretIn, userID uuid.UUID) (models.UploadSecretOut, error)
}

type uploadSecretUsecase struct {
	repo RecordRepository
}

// NewUploadSecretUsecase создает новый экземпляр UploadSecretUsecase.
// Принимает репозиторий для работы с записями секретов.
// Функция используется для внедрения зависимостей и создания usecase
// с конкретной реализацией хранилища секретов.
func NewUploadSecretUsecase(repo RecordRepository) UploadSecretUsecase {
	return &uploadSecretUsecase{repo: repo}
}

func (uc *uploadSecretUsecase) UploadSecret(ctx context.Context, in models.UploadSecretIn, userID uuid.UUID) (models.UploadSecretOut, error) {
	id := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now()
	version := 1 // всегда 1 для новой записи

	err := uc.repo.InsertRecord(ctx, id, userID, in.Label, in.Type, in.Metadata, in.EncryptedData, in.FileKey, "", version, createdAt, updatedAt)
	if err != nil {
		if err == models.ErrConflict {
			return models.UploadSecretOut{}, models.ErrConflict
		}
		return models.UploadSecretOut{}, err
	}
	return models.UploadSecretOut{ID: id.String(), Version: version}, nil
}
