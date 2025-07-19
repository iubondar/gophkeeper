package usecase_test

import (
	"context"
	"testing"

	"gophkeeper/internal/models"
	"gophkeeper/internal/server/storage/mocks"
	"gophkeeper/internal/server/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUploadSecretUsecase_UploadSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		in             models.UploadSecretIn
		repoError      error
		expectErr      bool
		expectConflict bool
	}{
		{
			name: "Success",
			in: models.UploadSecretIn{
				UserID:        uuid.New().String(),
				Label:         "label",
				Type:          "note",
				Metadata:      "meta",
				EncryptedData: "data",
				FileKey:       "key",
			},
		},
		{
			name: "Invalid userID",
			in: models.UploadSecretIn{
				UserID: "not-a-uuid",
			},
			expectErr: true,
		},
		{
			name: "Repo error",
			in: models.UploadSecretIn{
				UserID: uuid.New().String(),
				Label:  "label",
				Type:   "note",
			},
			repoError: assert.AnError,
			expectErr: true,
		},
		{
			name: "Conflict error",
			in: models.UploadSecretIn{
				UserID: uuid.New().String(),
				Label:  "label",
				Type:   "note",
			},
			repoError:      models.ErrConflict,
			expectErr:      true,
			expectConflict: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockRecordRepository(ctrl)
			if tt.in.UserID != "not-a-uuid" {
				mockRepo.EXPECT().
					InsertRecord(gomock.Any(), gomock.Any(), gomock.Any(), tt.in.Label, tt.in.Type, tt.in.Metadata, gomock.Any(), tt.in.FileKey, 1, gomock.Any(), gomock.Any()).
					Return(tt.repoError)
			}

			uc := usecase.NewUploadSecretUsecase(mockRepo)
			result, err := uc.UploadSecret(context.Background(), tt.in)
			if tt.expectErr {
				assert.Error(t, err)
				if tt.expectConflict {
					assert.Equal(t, models.ErrConflict, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 1, result.Version)
				assert.NotEmpty(t, result.ID)
			}
		})
	}
}
