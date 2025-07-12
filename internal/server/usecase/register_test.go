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

func TestRegisterUsecase_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test cases
	tests := []struct {
		name           string
		in             models.RegisterIn
		repoOk         bool
		repoError      error
		expectedUserID uuid.UUID
		expectedOk     bool
		expectedError  error
	}{
		{
			name: "Successful registration",
			in: models.RegisterIn{
				Login:        "testuser",
				PasswordHash: "testpass",
				Salt:         "testsalt",
			},
			repoOk:         true,
			expectedOk:     true,
			expectedUserID: uuid.New(),
		},
		{
			name: "Empty login",
			in: models.RegisterIn{
				Login:        "",
				PasswordHash: "testpass",
				Salt:         "testsalt",
			},
			expectedOk: false,
		},
		{
			name: "Empty password",
			in: models.RegisterIn{
				Login:        "testuser",
				PasswordHash: "",
				Salt:         "testsalt",
			},
			expectedOk: false,
		},
		{
			name: "User already exists",
			in: models.RegisterIn{
				Login:        "testuser",
				PasswordHash: "testpass",
				Salt:         "testsalt",
			},
			repoOk:     false,
			expectedOk: false,
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
			userID, ok, err := uc.Register(context.Background(), tt.in)

			// Check results
			if tt.expectedUserID != uuid.Nil {
				assert.NotEqual(t, uuid.Nil, userID)
			}
			assert.Equal(t, tt.expectedOk, ok)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
