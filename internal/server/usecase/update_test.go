package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"gophkeeper/internal/models"
	"gophkeeper/internal/server/storage/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestUpdateSecretUsecase_UpdateSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUpdateSecretRepository(ctrl)
	uc := NewUpdateSecretUsecase(mockRepo)

	ctx := context.Background()
	userID := uuid.New()
	in := models.UpdateSecretIn{
		UploadSecretIn: models.UploadSecretIn{
			Label:         "label",
			Type:          "note",
			Metadata:      "meta",
			EncryptedData: []byte("data"),
			FileKey:       "key",
		},
		Version: 1,
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().UpdateRecordByLabel(gomock.Any(), in.Label, userID, in.Type, in.Metadata, in.EncryptedData, in.FileKey, in.Version, gomock.AssignableToTypeOf(time.Now())).Return(2, nil)
		out, err := uc.UpdateSecret(ctx, in, userID)
		require.NoError(t, err)
		require.Equal(t, 2, out.Version)
	})

	t.Run("conflict", func(t *testing.T) {
		mockRepo.EXPECT().UpdateRecordByLabel(gomock.Any(), in.Label, userID, in.Type, in.Metadata, in.EncryptedData, in.FileKey, in.Version, gomock.AssignableToTypeOf(time.Now())).Return(0, models.ErrConflict)
		_, err := uc.UpdateSecret(ctx, in, userID)
		require.ErrorIs(t, err, models.ErrConflict)
	})

	t.Run("other error", func(t *testing.T) {
		errSome := errors.New("some error")
		mockRepo.EXPECT().UpdateRecordByLabel(gomock.Any(), in.Label, userID, in.Type, in.Metadata, in.EncryptedData, in.FileKey, in.Version, gomock.AssignableToTypeOf(time.Now())).Return(0, errSome)
		_, err := uc.UpdateSecret(ctx, in, userID)
		require.ErrorIs(t, err, errSome)
	})
}
