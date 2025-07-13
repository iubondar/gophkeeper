package usecase_test

import (
	"context"
	"testing"

	"gophkeeper/internal/server/storage/mocks"
	"gophkeeper/internal/server/usecase"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLoginUsecase_GetSalt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		login         string
		repoSalt      string
		repoError     error
		expectedSalt  string
		expectedError error
	}{
		{
			name:         "User found",
			login:        "testuser",
			repoSalt:     "dGVzdC1zYWx0",
			expectedSalt: "dGVzdC1zYWx0",
		},
		{
			name:         "User not found",
			login:        "nouser",
			repoSalt:     "",
			expectedSalt: "",
		},
		{
			name:          "Repository error",
			login:         "testuser",
			repoError:     assert.AnError,
			expectedError: assert.AnError,
		},
		{
			name:         "Empty login",
			login:        "",
			expectedSalt: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockLoginUserRepository(ctrl)
			if tt.login != "" {
				mockRepo.EXPECT().
					GetUserSalt(gomock.Any(), tt.login).
					Return(tt.repoSalt, tt.repoError)
			}

			uc := usecase.NewLoginUsecase(mockRepo)
			salt, err := uc.GetSalt(context.Background(), tt.login)
			assert.Equal(t, tt.expectedSalt, salt)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
