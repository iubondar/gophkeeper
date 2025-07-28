package cmd

import (
	"context"
	"errors"
	"testing"

	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/models"
)

// MockUpdateAPIClient представляет мок для UpdateAPIClient
type MockUpdateAPIClient struct {
	updateSecretFunc func(ctx context.Context, in models.UpdateSecretIn) error
}

func (m *MockUpdateAPIClient) UpdateSecret(ctx context.Context, in models.UpdateSecretIn) error {
	if m.updateSecretFunc != nil {
		return m.updateSecretFunc(ctx, in)
	}
	return nil
}

func TestNewUpdateCommand(t *testing.T) {
	mockAPI := &MockUpdateAPIClient{}
	crypto := crypto.NewCrypto()
	cmd := NewUpdateCommand(mockAPI, crypto)

	if cmd == nil {
		t.Fatal("NewUpdateCommand returned nil")
	}

	if cmd.apiClient != mockAPI {
		t.Error("apiClient not set correctly")
	}

	if cmd.crypto != crypto {
		t.Error("crypto not set correctly")
	}
}

func TestUpdateCommand_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("successful update text secret", func(t *testing.T) {
		textData := &models.TextSecretData{
			Name:     "test-secret",
			Text:     "updated text",
			Metadata: "updated metadata",
		}
		updateData := &UpdateData{
			SecretName: "test-secret",
			Version:    1,
			Type:       models.SecretTypeText,
			Data:       textData,
		}

		mockAPI := &MockUpdateAPIClient{
			updateSecretFunc: func(ctx context.Context, in models.UpdateSecretIn) error {
				return nil
			},
		}
		crypto := crypto.NewCrypto()
		// Устанавливаем encryption key для тестирования
		crypto.SetSalt("testsalt")
		err := crypto.GenerateAndStoreEncryptionKey("testpass")
		if err != nil {
			t.Fatalf("Failed to set encryption key: %v", err)
		}
		cmd := NewUpdateCommand(mockAPI, crypto)

		result, err := cmd.Execute(ctx, updateData)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("update API error", func(t *testing.T) {
		textData := &models.TextSecretData{
			Name:     "test-secret",
			Text:     "updated text",
			Metadata: "updated metadata",
		}
		updateData := &UpdateData{
			SecretName: "test-secret",
			Version:    1,
			Type:       models.SecretTypeText,
			Data:       textData,
		}

		expectedErr := errors.New("update failed")
		mockAPI := &MockUpdateAPIClient{
			updateSecretFunc: func(ctx context.Context, in models.UpdateSecretIn) error {
				return expectedErr
			},
		}
		crypto := crypto.NewCrypto()
		// Устанавливаем encryption key для тестирования
		crypto.SetSalt("testsalt")
		err := crypto.GenerateAndStoreEncryptionKey("testpass")
		if err != nil {
			t.Fatalf("Failed to set encryption key: %v", err)
		}
		cmd := NewUpdateCommand(mockAPI, crypto)

		result, err := cmd.Execute(ctx, updateData)

		if err == nil {
			t.Fatal("Expected error")
		}

		if err.Error() != "ошибка при обновлении секрета: "+expectedErr.Error() {
			t.Errorf("Expected error message containing 'ошибка при обновлении секрета', got '%s'", err.Error())
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("invalid args type", func(t *testing.T) {
		mockAPI := &MockUpdateAPIClient{}
		crypto := crypto.NewCrypto()
		cmd := NewUpdateCommand(mockAPI, crypto)

		result, err := cmd.Execute(ctx, "invalid args")

		if err == nil {
			t.Fatal("Expected error")
		}

		if err.Error() != "неверный тип аргументов для команды обновления" {
			t.Errorf("Expected error message 'неверный тип аргументов для команды обновления', got '%s'", err.Error())
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("unsupported data type", func(t *testing.T) {
		updateData := &UpdateData{
			SecretName: "test-secret",
			Version:    1,
			Type:       "unsupported",
			Data:       "unsupported data",
		}

		mockAPI := &MockUpdateAPIClient{}
		crypto := crypto.NewCrypto()
		cmd := NewUpdateCommand(mockAPI, crypto)

		result, err := cmd.Execute(ctx, updateData)

		if err == nil {
			t.Fatal("Expected error")
		}

		if err.Error() != "неподдерживаемый тип данных для обновления" {
			t.Errorf("Expected error message 'неподдерживаемый тип данных для обновления', got '%s'", err.Error())
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})
}

func TestUpdateCommand_GetName(t *testing.T) {
	mockAPI := &MockUpdateAPIClient{}
	crypto := crypto.NewCrypto()
	cmd := NewUpdateCommand(mockAPI, crypto)

	name := cmd.GetName()
	if name != "update" {
		t.Errorf("Expected name 'update', got '%s'", name)
	}
}

func TestUpdateCommand_GetDescription(t *testing.T) {
	mockAPI := &MockUpdateAPIClient{}
	crypto := crypto.NewCrypto()
	cmd := NewUpdateCommand(mockAPI, crypto)

	description := cmd.GetDescription()
	expected := "Обновить секрет на сервере"
	if description != expected {
		t.Errorf("Expected description '%s', got '%s'", expected, description)
	}
}
