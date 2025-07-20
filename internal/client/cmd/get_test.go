package cmd

import (
	"context"
	"encoding/json"
	"testing"

	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGophKeeperClient struct {
	mock.Mock
}

func (m *MockGophKeeperClient) Register(ctx context.Context, in models.RegisterIn) error {
	args := m.Called(ctx, in)
	return args.Error(0)
}

func (m *MockGophKeeperClient) Login(ctx context.Context, in models.LoginIn) (string, error) {
	args := m.Called(ctx, in)
	return args.String(0), args.Error(1)
}

func (m *MockGophKeeperClient) Authenticate(ctx context.Context, in models.AuthenticateIn) error {
	args := m.Called(ctx, in)
	return args.Error(0)
}

func (m *MockGophKeeperClient) UploadSecret(ctx context.Context, in models.UploadSecretIn) error {
	args := m.Called(ctx, in)
	return args.Error(0)
}

func (m *MockGophKeeperClient) UpdateSecret(ctx context.Context, secret models.SecretData) error {
	args := m.Called(ctx, secret)
	return args.Error(0)
}

func (m *MockGophKeeperClient) GetSecret(ctx context.Context, secretName string) (*models.GetSecretOut, error) {
	args := m.Called(ctx, secretName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GetSecretOut), args.Error(1)
}

func (m *MockGophKeeperClient) DeleteSecret(ctx context.Context, secretName string) error {
	args := m.Called(ctx, secretName)
	return args.Error(0)
}

func (m *MockGophKeeperClient) HealthCheck(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestGetCommand_Execute(t *testing.T) {
	tests := []struct {
		name          string
		secretName    string
		secretType    string
		encryptedData []byte
		setupMock     func(*MockGophKeeperClient, *crypto.Crypto)
		expectedError bool
		errorContains string
	}{
		{
			name:       "successful get text secret",
			secretName: "test-text",
			secretType: models.SecretTypeText,
			setupMock: func(mockClient *MockGophKeeperClient, crypto *crypto.Crypto) {
				// Создаем тестовые данные
				textSecret := models.TextSecretData{
					Name: "test-text",
					Text: "This is a test text secret",
				}
				jsonData, _ := json.Marshal(textSecret)
				encryptedData, _ := crypto.EncryptString(string(jsonData))

				mockClient.On("GetSecret", mock.Anything, "test-text").Return(&models.GetSecretOut{
					Label:         "test-text",
					Type:          models.SecretTypeText,
					Metadata:      "test metadata",
					EncryptedData: encryptedData,
				}, nil)
			},
			expectedError: false,
		},
		{
			name:       "successful get file secret",
			secretName: "test-file",
			secretType: models.SecretTypeFile,
			setupMock: func(mockClient *MockGophKeeperClient, crypto *crypto.Crypto) {
				// Создаем тестовые данные для файла
				fileData := models.FileData{
					Name:     "test-file",
					FilePath: "/path/to/file.txt",
				}
				jsonData, _ := json.Marshal(fileData)
				encryptedData, _ := crypto.EncryptString(string(jsonData))

				mockClient.On("GetSecret", mock.Anything, "test-file").Return(&models.GetSecretOut{
					Label:         "test-file",
					Type:          models.SecretTypeFile,
					Metadata:      "test metadata",
					EncryptedData: encryptedData,
				}, nil)
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockGophKeeperClient{}
			crypto := crypto.NewCrypto()
			crypto.SetSalt("dGVzdC1zYWx0") // base64 encoded "test-salt"
			err := crypto.GenerateAndStoreEncryptionKey("test-password")
			assert.NoError(t, err)

			tt.setupMock(mockClient, crypto)

			command := NewGetCommand(mockClient, crypto)

			result, err := command.Execute(context.Background(), tt.secretName)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if getResult, ok := result.(*GetSecretResult); ok {
					assert.Equal(t, tt.secretType, getResult.Type)
					assert.NotNil(t, getResult.Data)
				}
			}

			mockClient.AssertExpectations(t)
		})
	}
}
