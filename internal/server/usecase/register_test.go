package usecase_test

import (
	"context"
	"testing"

	"gophkeeper/internal/models"
	"gophkeeper/internal/server/storage/mocks"
	"gophkeeper/internal/server/usecase"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRegisterUsecase_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test cases
	tests := []struct {
		name          string
		in            models.RegisterIn
		repoOk        bool
		repoError     error
		expectedError error
		expectSuccess bool
	}{
		{
			name: "Successful registration",
			in: models.RegisterIn{
				Login:        "testuser",
				PasswordHash: "testpass",
				Salt:         "testsalt",
			},
			repoOk:        true,
			expectSuccess: true,
		},
		{
			name: "Empty login",
			in: models.RegisterIn{
				Login:        "",
				PasswordHash: "testpass",
				Salt:         "testsalt",
			},
			expectedError: models.ErrLoginOrPasswordEmpty,
		},
		{
			name: "Empty password",
			in: models.RegisterIn{
				Login:        "testuser",
				PasswordHash: "",
				Salt:         "testsalt",
			},
			expectedError: models.ErrLoginOrPasswordEmpty,
		},
		{
			name: "User already exists",
			in: models.RegisterIn{
				Login:        "testuser",
				PasswordHash: "testpass",
				Salt:         "testsalt",
			},
			repoOk:        false,
			expectedError: models.ErrUserAlreadyExists,
		},
		{
			name: "Repository error",
			in: models.RegisterIn{
				Login:        "testuser",
				PasswordHash: "testpass",
				Salt:         "testsalt",
			},
			repoError:     assert.AnError,
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock repository
			mockRepo := mocks.NewMockUserRepository(ctrl)
			if tt.in.Login != "" && tt.in.PasswordHash != "" {
				mockRepo.EXPECT().
					Register(gomock.Any(), gomock.Any(), tt.in.Login, tt.in.PasswordHash, tt.in.Salt).
					Return(tt.repoOk, tt.repoError)
			}

			// Create usecase
			uc := usecase.NewRegisterUsecase(mockRepo)

			// Call usecase
			result, err := uc.Register(context.Background(), tt.in)

			// Check results
			if tt.expectSuccess {
				assert.NoError(t, err)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			}
		})
	}
}
