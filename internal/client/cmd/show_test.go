package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"gophkeeper/internal/models"
)

// MockShowCryptoDecryptor представляет мок для ShowCryptoDecryptor
type MockShowCryptoDecryptor struct {
	decryptStringFunc func(encryptedData []byte) (string, error)
}

func (m *MockShowCryptoDecryptor) DecryptString(encryptedData []byte) (string, error) {
	if m.decryptStringFunc != nil {
		return m.decryptStringFunc(encryptedData)
	}
	return "", nil
}

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
	mockCrypto := &MockShowCryptoDecryptor{}
	cmd := NewShowCommand(mockAPI, mockCrypto)

	if cmd == nil {
		t.Fatal("NewShowCommand returned nil")
	}

	if cmd.apiClient != mockAPI {
		t.Error("apiClient not set correctly")
	}

	if cmd.crypto != mockCrypto {
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
		mockCrypto := &MockShowCryptoDecryptor{}
		cmd := NewShowCommand(mockAPI, mockCrypto)

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
		mockCrypto := &MockShowCryptoDecryptor{}
		cmd := NewShowCommand(mockAPI, mockCrypto)

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
		mockCrypto := &MockShowCryptoDecryptor{}
		cmd := NewShowCommand(mockAPI, mockCrypto)

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
		mockCrypto := &MockShowCryptoDecryptor{}
		cmd := NewShowCommand(mockAPI, mockCrypto)

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

	t.Run("successful show text secret", func(t *testing.T) {
		secretName := "test-secret"
		textData := models.TextSecretData{Name: "n", Text: "t", Metadata: "m"}
		jsonData, _ := json.Marshal(textData)
		mockAPI := &MockShowAPIClient{
			getSecretFunc: func(ctx context.Context, secretName string) (*models.GetSecretOut, error) {
				return &models.GetSecretOut{
					Label:         secretName,
					Type:          models.SecretTypeText,
					Metadata:      "meta",
					EncryptedData: jsonData,
				}, nil
			},
			getSecretVersionFunc: func(ctx context.Context, secretName string) (int, error) {
				return 1, nil
			},
		}
		mockCrypto := &MockShowCryptoDecryptor{
			decryptStringFunc: func(encryptedData []byte) (string, error) {
				return string(encryptedData), nil
			},
		}
		cmd := NewShowCommand(mockAPI, mockCrypto)

		result, err := cmd.Execute(ctx, secretName)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("expected result, got nil")
		}
	})

	t.Run("successful show login password secret", func(t *testing.T) {
		secretName := "test-secret"
		loginData := models.LoginPasswordData{Name: "n", Login: "l", Password: "p", URL: "u", Metadata: "m"}
		jsonData, _ := json.Marshal(loginData)
		mockAPI := &MockShowAPIClient{
			getSecretFunc: func(ctx context.Context, secretName string) (*models.GetSecretOut, error) {
				return &models.GetSecretOut{
					Label:         secretName,
					Type:          models.SecretTypeLoginPassword,
					Metadata:      "meta",
					EncryptedData: jsonData,
				}, nil
			},
			getSecretVersionFunc: func(ctx context.Context, secretName string) (int, error) {
				return 1, nil
			},
		}
		mockCrypto := &MockShowCryptoDecryptor{
			decryptStringFunc: func(encryptedData []byte) (string, error) {
				return string(encryptedData), nil
			},
		}
		cmd := NewShowCommand(mockAPI, mockCrypto)

		result, err := cmd.Execute(ctx, secretName)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("expected result, got nil")
		}
	})

	t.Run("successful show card secret", func(t *testing.T) {
		secretName := "test-secret"
		cardData := models.CardData{Name: "n", Number: "1", Holder: "h", Expiry: "e", CVV: "c", Metadata: "m"}
		jsonData, _ := json.Marshal(cardData)
		mockAPI := &MockShowAPIClient{
			getSecretFunc: func(ctx context.Context, secretName string) (*models.GetSecretOut, error) {
				return &models.GetSecretOut{
					Label:         secretName,
					Type:          models.SecretTypeCard,
					Metadata:      "meta",
					EncryptedData: jsonData,
				}, nil
			},
			getSecretVersionFunc: func(ctx context.Context, secretName string) (int, error) {
				return 1, nil
			},
		}
		mockCrypto := &MockShowCryptoDecryptor{
			decryptStringFunc: func(encryptedData []byte) (string, error) {
				return string(encryptedData), nil
			},
		}
		cmd := NewShowCommand(mockAPI, mockCrypto)

		result, err := cmd.Execute(ctx, secretName)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("expected result, got nil")
		}
	})
}

func TestHandleShowTextSecret(t *testing.T) {
	valid := `{"name":"n","text":"t","metadata":"m"}`
	res, err := handleShowTextSecret(valid, "meta", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Type != models.SecretTypeText {
		t.Errorf("expected type %s, got %s", models.SecretTypeText, res.Type)
	}
	if res.Metadata != "meta" {
		t.Errorf("expected metadata 'meta', got '%s'", res.Metadata)
	}
	if res.Version != 1 {
		t.Errorf("expected version 1, got %d", res.Version)
	}

	_, err = handleShowTextSecret("not json", "meta", 1)
	if err == nil {
		t.Error("expected error for invalid json")
	}
}

func TestHandleShowLoginPasswordSecret(t *testing.T) {
	valid := `{"name":"n","login":"l","password":"p","url":"u","metadata":"m"}`
	res, err := handleShowLoginPasswordSecret(valid, "meta", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Type != models.SecretTypeLoginPassword {
		t.Errorf("expected type %s, got %s", models.SecretTypeLoginPassword, res.Type)
	}
	if res.Metadata != "meta" {
		t.Errorf("expected metadata 'meta', got '%s'", res.Metadata)
	}
	if res.Version != 1 {
		t.Errorf("expected version 1, got %d", res.Version)
	}

	_, err = handleShowLoginPasswordSecret("not json", "meta", 1)
	if err == nil {
		t.Error("expected error for invalid json")
	}
}

func TestHandleShowCardSecret(t *testing.T) {
	valid := `{"name":"n","number":"1","holder":"h","expiry":"e","cvv":"c","metadata":"m"}`
	res, err := handleShowCardSecret(valid, "meta", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Type != models.SecretTypeCard {
		t.Errorf("expected type %s, got %s", models.SecretTypeCard, res.Type)
	}
	if res.Metadata != "meta" {
		t.Errorf("expected metadata 'meta', got '%s'", res.Metadata)
	}
	if res.Version != 1 {
		t.Errorf("expected version 1, got %d", res.Version)
	}

	_, err = handleShowCardSecret("not json", "meta", 1)
	if err == nil {
		t.Error("expected error for invalid json")
	}
}

func TestShowCommand_GetName(t *testing.T) {
	mockAPI := &MockShowAPIClient{}
	mockCrypto := &MockShowCryptoDecryptor{}
	cmd := NewShowCommand(mockAPI, mockCrypto)

	name := cmd.GetName()
	if name != "show" {
		t.Errorf("Expected name 'show', got '%s'", name)
	}
}

func TestShowCommand_GetDescription(t *testing.T) {
	mockAPI := &MockShowAPIClient{}
	mockCrypto := &MockShowCryptoDecryptor{}
	cmd := NewShowCommand(mockAPI, mockCrypto)

	description := cmd.GetDescription()
	expected := "Показать текущие данные секрета"
	if description != expected {
		t.Errorf("Expected description '%s', got '%s'", expected, description)
	}
}
