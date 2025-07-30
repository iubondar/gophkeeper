package usecase

import (
	"context"
	"gophkeeper/internal/models"
	"time"

	"github.com/google/uuid"
)

// UpdateSecretRepository определяет интерфейс для обновления секретов в хранилище.
// Интерфейс используется для абстракции от конкретной реализации хранилища
// и позволяет тестировать usecase с помощью моков.
type UpdateSecretRepository interface {
	UpdateRecordByLabel(ctx context.Context, label string, userID uuid.UUID, recordType, metadata string, encryptedData []byte, fileKey string, fileName string, expectedVersion int, updatedAt time.Time) (int, error)
}

// UpdateSecretUsecase определяет интерфейс для обновления секретов.
// Интерфейс содержит бизнес-логику обновления зашифрованных секретов
// с проверкой версий для предотвращения конфликтов.
type UpdateSecretUsecase interface {
	UpdateSecret(ctx context.Context, in models.UpdateSecretIn, userID uuid.UUID) (models.UpdateSecretOut, error)
}

type updateSecretUsecase struct {
	repo UpdateSecretRepository
}

// NewUpdateSecretUsecase создает новый экземпляр UpdateSecretUsecase.
// Принимает репозиторий для обновления секретов.
// Функция используется для внедрения зависимостей и создания usecase
// с конкретной реализацией хранилища секретов.
func NewUpdateSecretUsecase(repo UpdateSecretRepository) UpdateSecretUsecase {
	return &updateSecretUsecase{repo: repo}
}

func (uc *updateSecretUsecase) UpdateSecret(ctx context.Context, in models.UpdateSecretIn, userID uuid.UUID) (models.UpdateSecretOut, error) {
	updatedAt := time.Now()
	newVersion, err := uc.repo.UpdateRecordByLabel(ctx, in.Label, userID, in.Type, in.Metadata, in.EncryptedData, in.FileKey, "", in.Version, updatedAt)
	if err != nil {
		if err == models.ErrConflict {
			return models.UpdateSecretOut{}, models.ErrConflict
		}
		return models.UpdateSecretOut{}, err
	}
	return models.UpdateSecretOut{Version: newVersion}, nil
}
