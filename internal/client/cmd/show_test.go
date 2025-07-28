package cmd

import (
	"context"
	"errors"
	"testing"

	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/models"
)

// MockShowAPIClient представляет мок для ShowAPIClient
type MockShowAPIClient struct {
	getSecretFunc        func(ctx context.Context, secretName string) (*models.GetSecretOut, error)
	getSecretVersionFunc func(ctx context.Context, secretName string) (int, error)
}

func (m *MockShowAPIClient) GetSecret(ctx context.Context, secretName string) (*models.GetSecretOut, error) {
	if m.getSecretFunc != nil {
		return m.getSecretFunc(ctx, secretName)
	}
	return nil, nil
}

func (m *MockShowAPIClient) GetSecretVersion(ctx context.Context, secretName string) (int, error) {
	if m.getSecretVersionFunc != nil {
		return m.getSecretVersionFunc(ctx, secretName)
	}
	return 0, nil
}

func TestNewShowCommand(t *testing.T) {
	mockAPI := &MockShowAPIClient{}
	crypto := crypto.NewCrypto()
	cmd := NewShowCommand(mockAPI, crypto)

	if cmd == nil {
		t.Fatal("NewShowCommand returned nil")
	}

	if cmd.apiClient != mockAPI {
		t.Error("apiClient not set correctly")
	}

	if cmd.crypto != crypto {
		t.Error("crypto not set correctly")
	}
}

func TestShowCommand_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("get secret API error", func(t *testing.T) {
		secretName := "test-secret"
		expectedErr := errors.New("secret not found")
		mockAPI := &MockShowAPIClient{
			getSecretFunc: func(ctx context.Context, secretName string) (*models.GetSecretOut, error) {
				return nil, expectedErr
			},
		}
		crypto := crypto.NewCrypto()
		cmd := NewShowCommand(mockAPI, crypto)

		result, err := cmd.Execute(ctx, secretName)

		if err == nil {
			t.Fatal("Expected error")
		}

		if err.Error() != "ошибка при получении информации о секрете: "+expectedErr.Error() {
			t.Errorf("Expected error message containing 'ошибка при получении информации о секрете', got '%s'", err.Error())
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("get version API error", func(t *testing.T) {
		secretName := "test-secret"
		expectedErr := errors.New("version not found")
		mockAPI := &MockShowAPIClient{
			getSecretFunc: func(ctx context.Context, secretName string) (*models.GetSecretOut, error) {
				return &models.GetSecretOut{
					Label:         secretName,
					Type:          models.SecretTypeText,
					Metadata:      "test metadata",
					EncryptedData: []byte("encrypted-data"),
				}, nil
			},
			getSecretVersionFunc: func(ctx context.Context, secretName string) (int, error) {
				return 0, expectedErr
			},
		}
		crypto := crypto.NewCrypto()
		cmd := NewShowCommand(mockAPI, crypto)

		result, err := cmd.Execute(ctx, secretName)

		if err == nil {
			t.Fatal("Expected error")
		}

		if err.Error() != "ошибка при получении версии секрета: "+expectedErr.Error() {
			t.Errorf("Expected error message containing 'ошибка при получении версии секрета', got '%s'", err.Error())
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("invalid args type", func(t *testing.T) {
		mockAPI := &MockShowAPIClient{}
		crypto := crypto.NewCrypto()
		cmd := NewShowCommand(mockAPI, crypto)

		result, err := cmd.Execute(ctx, 123) // int instead of string

		if err == nil {
			t.Fatal("Expected error")
		}

		if err.Error() != "неверный тип аргументов для команды отображения" {
			t.Errorf("Expected error message 'неверный тип аргументов для команды отображения', got '%s'", err.Error())
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("unsupported secret type", func(t *testing.T) {
		secretName := "test-secret"
		mockAPI := &MockShowAPIClient{
			getSecretFunc: func(ctx context.Context, secretName string) (*models.GetSecretOut, error) {
				return &models.GetSecretOut{
					Label:         secretName,
					Type:          "unsupported",
					Metadata:      "test metadata",
					EncryptedData: []byte("encrypted-data"),
				}, nil
			},
			getSecretVersionFunc: func(ctx context.Context, secretName string) (int, error) {
				return 1, nil
			},
		}
		crypto := crypto.NewCrypto()
		// Устанавливаем encryption key для тестирования
		crypto.SetSalt("testsalt")
		err := crypto.GenerateAndStoreEncryptionKey("testpass")
		if err != nil {
			t.Fatalf("Failed to set encryption key: %v", err)
		}
		cmd := NewShowCommand(mockAPI, crypto)

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

func TestShowCommand_GetName(t *testing.T) {
	mockAPI := &MockShowAPIClient{}
	crypto := crypto.NewCrypto()
	cmd := NewShowCommand(mockAPI, crypto)

	name := cmd.GetName()
	if name != "show" {
		t.Errorf("Expected name 'show', got '%s'", name)
	}
}

func TestShowCommand_GetDescription(t *testing.T) {
	mockAPI := &MockShowAPIClient{}
	crypto := crypto.NewCrypto()
	cmd := NewShowCommand(mockAPI, crypto)

	description := cmd.GetDescription()
	expected := "Показать текущие данные секрета"
	if description != expected {
		t.Errorf("Expected description '%s', got '%s'", expected, description)
	}
}
