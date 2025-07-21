package usecase

import (
	"context"

	"github.com/google/uuid"
)

type DeleteSecretRepository interface {
	DeleteRecordByLabel(ctx context.Context, label string, userID uuid.UUID) error
}

type DeleteSecretUsecase interface {
	DeleteSecret(ctx context.Context, secretName string, userID uuid.UUID) error
}

type deleteSecretUsecase struct {
	repo DeleteSecretRepository
}

func NewDeleteSecretUsecase(repo DeleteSecretRepository) DeleteSecretUsecase {
	return &deleteSecretUsecase{repo: repo}
}

func (uc *deleteSecretUsecase) DeleteSecret(ctx context.Context, secretName string, userID uuid.UUID) error {
	return uc.repo.DeleteRecordByLabel(ctx, secretName, userID)
}
