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
		userID         uuid.UUID
		repoError      error
		expectErr      bool
		expectConflict bool
	}{
		{
			name: "Success",
			in: models.UploadSecretIn{
				Label:         "label",
				Type:          "note",
				Metadata:      "meta",
				EncryptedData: []byte("data"),
				FileKey:       "key",
			},
			userID: uuid.New(),
		},
		{
			name: "Repo error",
			in: models.UploadSecretIn{
				Label: "label",
				Type:  "note",
			},
			userID:    uuid.New(),
			repoError: assert.AnError,
			expectErr: true,
		},
		{
			name: "Conflict error",
			in: models.UploadSecretIn{
				Label: "label",
				Type:  "note",
			},
			userID:         uuid.New(),
			repoError:      models.ErrConflict,
			expectErr:      true,
			expectConflict: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockRecordRepository(ctrl)
			mockRepo.EXPECT().
				InsertRecord(gomock.Any(), gomock.Any(), tt.userID, tt.in.Label, tt.in.Type, tt.in.Metadata, gomock.Any(), tt.in.FileKey, 1, gomock.Any(), gomock.Any()).
				Return(tt.repoError)

			uc := usecase.NewUploadSecretUsecase(mockRepo)
			result, err := uc.UploadSecret(context.Background(), tt.in, tt.userID)
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
