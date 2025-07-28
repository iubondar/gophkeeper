package cmd

import (
	"context"
	"errors"
	"io"
	"testing"

	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/models"
)

// MockGetAPIClient представляет мок для GetAPIClient
type MockGetAPIClient struct {
	getSecretFunc    func(ctx context.Context, secretName string) (*models.GetSecretOut, error)
	downloadFileFunc func(ctx context.Context, label string) (io.ReadCloser, error)
}

func (m *MockGetAPIClient) GetSecret(ctx context.Context, secretName string) (*models.GetSecretOut, error) {
	if m.getSecretFunc != nil {
		return m.getSecretFunc(ctx, secretName)
	}
	return nil, nil
}

func (m *MockGetAPIClient) DownloadFile(ctx context.Context, label string) (io.ReadCloser, error) {
	if m.downloadFileFunc != nil {
		return m.downloadFileFunc(ctx, label)
	}
	return nil, nil
}

func TestNewGetCommand(t *testing.T) {
	mockAPI := &MockGetAPIClient{}
	crypto := crypto.NewCrypto()
	cmd := NewGetCommand(mockAPI, crypto)

	if cmd == nil {
		t.Fatal("NewGetCommand returned nil")
	}

	if cmd.apiClient != mockAPI {
		t.Error("apiClient not set correctly")
	}

	if cmd.crypto != crypto {
		t.Error("crypto not set correctly")
	}
}

func TestGetCommand_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("get secret API error", func(t *testing.T) {
		secretName := "test-secret"
		expectedErr := errors.New("secret not found")
		mockAPI := &MockGetAPIClient{
			getSecretFunc: func(ctx context.Context, secretName string) (*models.GetSecretOut, error) {
				return nil, expectedErr
			},
		}
		crypto := crypto.NewCrypto()
		cmd := NewGetCommand(mockAPI, crypto)

		result, err := cmd.Execute(ctx, secretName)

		if err == nil {
			t.Fatal("Expected error")
		}

		if err.Error() != "ошибка при получении секрета: "+expectedErr.Error() {
			t.Errorf("Expected error message containing 'ошибка при получении секрета', got '%s'", err.Error())
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("invalid args type", func(t *testing.T) {
		mockAPI := &MockGetAPIClient{}
		crypto := crypto.NewCrypto()
		cmd := NewGetCommand(mockAPI, crypto)

		result, err := cmd.Execute(ctx, 123) // int instead of string

		if err == nil {
			t.Fatal("Expected error")
		}

		if err.Error() != "неверный тип аргументов для команды получения" {
			t.Errorf("Expected error message 'неверный тип аргументов для команды получения', got '%s'", err.Error())
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("unsupported secret type", func(t *testing.T) {
		secretName := "test-secret"
		mockAPI := &MockGetAPIClient{
			getSecretFunc: func(ctx context.Context, secretName string) (*models.GetSecretOut, error) {
				return &models.GetSecretOut{
					Label:         secretName,
					Type:          "unsupported",
					Metadata:      "test metadata",
					EncryptedData: []byte("encrypted-data"),
				}, nil
			},
		}
		crypto := crypto.NewCrypto()
		// Устанавливаем encryption key для тестирования
		crypto.SetSalt("testsalt")
		err := crypto.GenerateAndStoreEncryptionKey("testpass")
		if err != nil {
			t.Fatalf("Failed to set encryption key: %v", err)
		}
		cmd := NewGetCommand(mockAPI, crypto)

		result, err := cmd.Execute(ctx, secretName)

		if err == nil {
			t.Fatal("Expected error")
		}

		if err.Error() != "неподдерживаемый тип секрета: unsupported" {
			t.Errorf("Expected error message 'неподдерживаемый тип секрета: unsupported', got '%s'", err.Error())
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})
}

func TestGetCommand_GetName(t *testing.T) {
	mockAPI := &MockGetAPIClient{}
	crypto := crypto.NewCrypto()
	cmd := NewGetCommand(mockAPI, crypto)

	name := cmd.GetName()
	if name != "get" {
		t.Errorf("Expected name 'get', got '%s'", name)
	}
}

func TestGetCommand_GetDescription(t *testing.T) {
	mockAPI := &MockGetAPIClient{}
	crypto := crypto.NewCrypto()
	cmd := NewGetCommand(mockAPI, crypto)

	description := cmd.GetDescription()
	expected := "Получить секрет с сервера"
	if description != expected {
		t.Errorf("Expected description '%s', got '%s'", expected, description)
	}
}
