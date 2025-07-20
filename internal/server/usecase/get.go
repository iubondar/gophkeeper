package usecase

import (
	"context"
	"gophkeeper/internal/models"

	"github.com/google/uuid"
)

type GetSecretRepository interface {
	GetRecordByLabel(ctx context.Context, label string, userID uuid.UUID) (*models.GetSecretOut, error)
}

type GetSecretUsecase interface {
	GetSecret(ctx context.Context, secretName string, userID uuid.UUID) (*models.GetSecretOut, error)
}

type getSecretUsecase struct {
	repo GetSecretRepository
}

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
