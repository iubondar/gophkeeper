package cmd

import (
	"context"
	"errors"
	"testing"

	"gophkeeper/internal/models"
)

// MockUpdateCryptoEncryptor представляет мок для UpdateCryptoEncryptor
type MockUpdateCryptoEncryptor struct {
	encryptStringFunc func(plaintext string) ([]byte, error)
}

func (m *MockUpdateCryptoEncryptor) EncryptString(plaintext string) ([]byte, error) {
	if m.encryptStringFunc != nil {
		return m.encryptStringFunc(plaintext)
	}
	return []byte("encrypted"), nil
}

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
	mockCrypto := &MockUpdateCryptoEncryptor{}
	cmd := NewUpdateCommand(mockAPI, mockCrypto)

	if cmd == nil {
		t.Fatal("NewUpdateCommand returned nil")
	}

	if cmd.apiClient != mockAPI {
		t.Error("apiClient not set correctly")
	}

	if cmd.crypto != mockCrypto {
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
		mockCrypto := &MockUpdateCryptoEncryptor{
			encryptStringFunc: func(plaintext string) ([]byte, error) {
				return []byte("encrypted"), nil
			},
		}
		cmd := NewUpdateCommand(mockAPI, mockCrypto)

		result, err := cmd.Execute(ctx, updateData)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("successful update login password secret", func(t *testing.T) {
		loginData := &models.LoginPasswordData{
			Name:     "test-secret",
			Login:    "updated-login",
			Password: "updated-password",
			URL:      "https://example.com",
			Metadata: "updated metadata",
		}
		updateData := &UpdateData{
			SecretName: "test-secret",
			Version:    1,
			Type:       models.SecretTypeLoginPassword,
			Data:       loginData,
		}

		mockAPI := &MockUpdateAPIClient{
			updateSecretFunc: func(ctx context.Context, in models.UpdateSecretIn) error {
				return nil
			},
		}
		mockCrypto := &MockUpdateCryptoEncryptor{
			encryptStringFunc: func(plaintext string) ([]byte, error) {
				return []byte("encrypted"), nil
			},
		}
		cmd := NewUpdateCommand(mockAPI, mockCrypto)

		result, err := cmd.Execute(ctx, updateData)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("successful update card secret", func(t *testing.T) {
		cardData := &models.CardData{
			Name:     "test-secret",
			Number:   "1234567890123456",
			Holder:   "JOHN DOE",
			Expiry:   "12/25",
			CVV:      "123",
			Metadata: "updated metadata",
		}
		updateData := &UpdateData{
			SecretName: "test-secret",
			Version:    1,
			Type:       models.SecretTypeCard,
			Data:       cardData,
		}

		mockAPI := &MockUpdateAPIClient{
			updateSecretFunc: func(ctx context.Context, in models.UpdateSecretIn) error {
				return nil
			},
		}
		mockCrypto := &MockUpdateCryptoEncryptor{
			encryptStringFunc: func(plaintext string) ([]byte, error) {
				return []byte("encrypted"), nil
			},
		}
		cmd := NewUpdateCommand(mockAPI, mockCrypto)

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
		mockCrypto := &MockUpdateCryptoEncryptor{
			encryptStringFunc: func(plaintext string) ([]byte, error) {
				return []byte("encrypted"), nil
			},
		}
		cmd := NewUpdateCommand(mockAPI, mockCrypto)

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
		mockCrypto := &MockUpdateCryptoEncryptor{}
		cmd := NewUpdateCommand(mockAPI, mockCrypto)

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
		mockCrypto := &MockUpdateCryptoEncryptor{}
		cmd := NewUpdateCommand(mockAPI, mockCrypto)

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
	mockCrypto := &MockUpdateCryptoEncryptor{}
	cmd := NewUpdateCommand(mockAPI, mockCrypto)

	name := cmd.GetName()
	if name != "update" {
		t.Errorf("Expected name 'update', got '%s'", name)
	}
}

func TestUpdateCommand_GetDescription(t *testing.T) {
	mockAPI := &MockUpdateAPIClient{}
	mockCrypto := &MockUpdateCryptoEncryptor{}
	cmd := NewUpdateCommand(mockAPI, mockCrypto)

	description := cmd.GetDescription()
	expected := "Обновить секрет на сервере"
	if description != expected {
		t.Errorf("Expected description '%s', got '%s'", expected, description)
	}
}
