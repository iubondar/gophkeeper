package usecase

import (
	"context"
	"gophkeeper/internal/models"
	"time"

	"github.com/google/uuid"
)

type RecordRepository interface {
	InsertRecord(ctx context.Context, id, userID uuid.UUID, label, recordType, metadata string, encryptedData []byte, fileKey string, version int, createdAt, updatedAt time.Time) error
}

type UploadSecretUsecase interface {
	UploadSecret(ctx context.Context, in models.UploadSecretIn) (models.UploadSecretOut, error)
}

type uploadSecretUsecase struct {
	repo RecordRepository
}

func NewUploadSecretUsecase(repo RecordRepository) UploadSecretUsecase {
	return &uploadSecretUsecase{repo: repo}
}

func (uc *uploadSecretUsecase) UploadSecret(ctx context.Context, in models.UploadSecretIn) (models.UploadSecretOut, error) {
	id := uuid.New()
	userID, err := uuid.Parse(in.UserID)
	if err != nil {
		return models.UploadSecretOut{}, err
	}
	createdAt := time.Now()
	updatedAt := time.Now()
	version := 1 // всегда 1 для новой записи
	var encryptedData []byte
	if in.EncryptedData != "" {
		encryptedData = []byte(in.EncryptedData)
	}
	err = uc.repo.InsertRecord(ctx, id, userID, in.Label, in.Type, in.Metadata, encryptedData, in.FileKey, version, createdAt, updatedAt)
	if err != nil {
		if err == models.ErrConflict {
			return models.UploadSecretOut{}, models.ErrConflict
		}
		return models.UploadSecretOut{}, err
	}
	return models.UploadSecretOut{ID: id.String(), Version: version}, nil
}
