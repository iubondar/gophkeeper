package usecase

import (
	"context"
	"testing"

	"gophkeeper/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGetSecretRepository struct {
	mock.Mock
}

func (m *MockGetSecretRepository) GetRecordByLabel(ctx context.Context, label string, userID uuid.UUID) (*models.GetSecretOut, error) {
	args := m.Called(ctx, label, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GetSecretOut), args.Error(1)
}

func TestGetSecretUsecase_GetSecret(t *testing.T) {
	tests := []struct {
		name           string
		secretName     string
		userID         uuid.UUID
		setupMock      func(*MockGetSecretRepository)
		expectedResult *models.GetSecretOut
		expectedError  error
	}{
		{
			name:       "successful get secret",
			secretName: "test-secret",
			userID:     uuid.New(),
			setupMock: func(mockRepo *MockGetSecretRepository) {
				expectedOut := &models.GetSecretOut{
					ID:            "test-id",
					Label:         "test-secret",
					Type:          "text",
					Metadata:      "test metadata",
					EncryptedData: []byte("encrypted-data"),
					FileKey:       "",
					Version:       1,
				}
				mockRepo.On("GetRecordByLabel", mock.Anything, "test-secret", mock.Anything).Return(expectedOut, nil)
			},
			expectedResult: &models.GetSecretOut{
				ID:            "test-id",
				Label:         "test-secret",
				Type:          "text",
				Metadata:      "test metadata",
				EncryptedData: []byte("encrypted-data"),
				FileKey:       "",
				Version:       1,
			},
			expectedError: nil,
		},
		{
			name:       "record not found",
			secretName: "non-existent",
			userID:     uuid.New(),
			setupMock: func(mockRepo *MockGetSecretRepository) {
				mockRepo.On("GetRecordByLabel", mock.Anything, "non-existent", mock.Anything).Return(nil, models.ErrRecordNotFound)
			},
			expectedResult: nil,
			expectedError:  models.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockGetSecretRepository{}
			tt.setupMock(mockRepo)

			usecase := NewGetSecretUsecase(mockRepo)

			result, err := usecase.GetSecret(context.Background(), tt.secretName, tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
