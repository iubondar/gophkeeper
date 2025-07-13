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

func TestAuthenticateUsecase_Authenticate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validUserID := uuid.New()

	tests := []struct {
		name          string
		login         string
		passwordHash  string
		repoUserID    uuid.UUID
		repoError     error
		expectedError error
		expectSuccess bool
	}{
		{
			name:          "Successful authentication",
			login:         "testuser",
			passwordHash:  "validhash",
			repoUserID:    validUserID,
			expectSuccess: true,
		},
		{
			name:          "User not found",
			login:         "nouser",
			passwordHash:  "invalidhash",
			repoUserID:    uuid.Nil,
			expectedError: models.ErrUserNotFound,
		},
		{
			name:          "Repository error",
			login:         "testuser",
			passwordHash:  "validhash",
			repoError:     assert.AnError,
			expectedError: assert.AnError,
		},
		{
			name:          "Empty login",
			login:         "",
			passwordHash:  "validhash",
			expectedError: models.ErrLoginOrPasswordEmpty,
		},
		{
			name:          "Empty password hash",
			login:         "testuser",
			passwordHash:  "",
			expectedError: models.ErrLoginOrPasswordEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockAuthenticateUserRepository(ctrl)
			if tt.login != "" && tt.passwordHash != "" {
				mockRepo.EXPECT().
					GetUserByLoginAndPassword(gomock.Any(), tt.login, tt.passwordHash).
					Return(tt.repoUserID, tt.repoError)
			}

			uc := usecase.NewAuthenticateUsecase(mockRepo)
			result, err := uc.Authenticate(context.Background(), tt.login, tt.passwordHash)

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
				assert.Equal(t, models.AuthenticateOut{}, result)
			} else {
				assert.NoError(t, err)
				if tt.expectSuccess {
					assert.NotEmpty(t, result.AccessToken)
					assert.NotEmpty(t, result.RefreshToken)
				} else {
					assert.Empty(t, result.AccessToken)
					assert.Empty(t, result.RefreshToken)
				}
			}
		})
	}
}
