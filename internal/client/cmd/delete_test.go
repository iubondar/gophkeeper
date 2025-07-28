package cmd

import (
	"context"
	"errors"
	"testing"

	"gophkeeper/internal/client/crypto"
)

// MockDeleteAPIClient представляет мок для DeleteAPIClient
type MockDeleteAPIClient struct {
	deleteSecretFunc func(ctx context.Context, secretName string) error
}

func (m *MockDeleteAPIClient) DeleteSecret(ctx context.Context, secretName string) error {
	if m.deleteSecretFunc != nil {
		return m.deleteSecretFunc(ctx, secretName)
	}
	return nil
}

func TestNewDeleteCommand(t *testing.T) {
	mockAPI := &MockDeleteAPIClient{}
	crypto := crypto.NewCrypto()
	cmd := NewDeleteCommand(mockAPI, crypto)

	if cmd == nil {
		t.Fatal("NewDeleteCommand returned nil")
	}

	if cmd.apiClient != mockAPI {
		t.Error("apiClient not set correctly")
	}

	if cmd.crypto != crypto {
		t.Error("crypto not set correctly")
	}
}

func TestDeleteCommand_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("successful delete", func(t *testing.T) {
		secretName := "test-secret"
		mockAPI := &MockDeleteAPIClient{
			deleteSecretFunc: func(ctx context.Context, secretName string) error {
				return nil
			},
		}
		crypto := crypto.NewCrypto()
		cmd := NewDeleteCommand(mockAPI, crypto)

		result, err := cmd.Execute(ctx, secretName)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("delete with API error", func(t *testing.T) {
		secretName := "test-secret"
		expectedErr := errors.New("secret not found")
		mockAPI := &MockDeleteAPIClient{
			deleteSecretFunc: func(ctx context.Context, secretName string) error {
				return expectedErr
			},
		}
		crypto := crypto.NewCrypto()
		cmd := NewDeleteCommand(mockAPI, crypto)

		result, err := cmd.Execute(ctx, secretName)

		if err == nil {
			t.Fatal("Expected error")
		}

		if err.Error() != "ошибка при удалении секрета: "+expectedErr.Error() {
			t.Errorf("Expected error message containing 'ошибка при удалении секрета', got '%s'", err.Error())
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("invalid args type", func(t *testing.T) {
		mockAPI := &MockDeleteAPIClient{}
		crypto := crypto.NewCrypto()
		cmd := NewDeleteCommand(mockAPI, crypto)

		result, err := cmd.Execute(ctx, 123) // int instead of string

		if err == nil {
			t.Fatal("Expected error")
		}

		if err.Error() != "неверный тип аргументов для команды удаления" {
			t.Errorf("Expected error message 'неверный тип аргументов для команды удаления', got '%s'", err.Error())
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})
}

func TestDeleteCommand_GetName(t *testing.T) {
	mockAPI := &MockDeleteAPIClient{}
	crypto := crypto.NewCrypto()
	cmd := NewDeleteCommand(mockAPI, crypto)

	name := cmd.GetName()
	if name != "delete" {
		t.Errorf("Expected name 'delete', got '%s'", name)
	}
}

func TestDeleteCommand_GetDescription(t *testing.T) {
	mockAPI := &MockDeleteAPIClient{}
	crypto := crypto.NewCrypto()
	cmd := NewDeleteCommand(mockAPI, crypto)

	description := cmd.GetDescription()
	expected := "Удалить секрет с сервера"
	if description != expected {
		t.Errorf("Expected description '%s', got '%s'", expected, description)
	}
}
