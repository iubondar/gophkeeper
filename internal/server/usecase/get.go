package usecase

import (
	"context"
	"gophkeeper/internal/models"

	"github.com/google/uuid"
)

// GetSecretRepository определяет интерфейс для получения секретов из хранилища.
// Интерфейс используется для абстракции от конкретной реализации хранилища
// и позволяет тестировать usecase с помощью моков.
type GetSecretRepository interface {
	GetRecordByLabel(ctx context.Context, label string, userID uuid.UUID) (*models.GetSecretOut, error)
}

// GetSecretUsecase определяет интерфейс для получения секретов.
// Интерфейс содержит бизнес-логику получения зашифрованных секретов
// из хранилища по имени и идентификатору пользователя.
type GetSecretUsecase interface {
	GetSecret(ctx context.Context, secretName string, userID uuid.UUID) (*models.GetSecretOut, error)
}

type getSecretUsecase struct {
	repo GetSecretRepository
}

// NewGetSecretUsecase создает новый экземпляр GetSecretUsecase.
// Принимает репозиторий для получения секретов.
// Функция используется для внедрения зависимостей и создания usecase
// с конкретной реализацией хранилища секретов.
func NewGetSecretUsecase(repo GetSecretRepository) GetSecretUsecase {
	return &getSecretUsecase{repo: repo}
}

func (uc *getSecretUsecase) GetSecret(ctx context.Context, secretName string, userID uuid.UUID) (*models.GetSecretOut, error) {
	record, err := uc.repo.GetRecordByLabel(ctx, secretName, userID)
	if err != nil {
		return nil, err
	}
	return record, nil
}
