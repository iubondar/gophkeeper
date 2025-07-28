package cmd

import (
	"context"
	"errors"
	"io"
	"testing"

	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/models"
)

// MockUploadAPIClient представляет мок для UploadAPIClient
type MockUploadAPIClient struct {
	uploadSecretFunc func(ctx context.Context, in models.UploadSecretIn) error
	uploadFileFunc   func(ctx context.Context, label, metadata string, file io.Reader, filename string) (*models.UploadSecretOut, error)
}

func (m *MockUploadAPIClient) UploadSecret(ctx context.Context, in models.UploadSecretIn) error {
	if m.uploadSecretFunc != nil {
		return m.uploadSecretFunc(ctx, in)
	}
	return nil
}

func (m *MockUploadAPIClient) UploadFile(ctx context.Context, label, metadata string, file io.Reader, filename string) (*models.UploadSecretOut, error) {
	if m.uploadFileFunc != nil {
		return m.uploadFileFunc(ctx, label, metadata, file, filename)
	}
	return nil, nil
}

func TestNewUploadCommand(t *testing.T) {
	mockAPI := &MockUploadAPIClient{}
	crypto := crypto.NewCrypto()
	cmd := NewUploadCommand(mockAPI, crypto)

	if cmd == nil {
		t.Fatal("NewUploadCommand returned nil")
	}

	if cmd.apiClient != mockAPI {
		t.Error("apiClient not set correctly")
	}

	if cmd.crypto != crypto {
		t.Error("crypto not set correctly")
	}
}

func TestUploadCommand_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("successful upload text secret", func(t *testing.T) {
		textData := &models.TextSecretData{
			Name:     "test-secret",
			Text:     "test text",
			Metadata: "test metadata",
		}

		mockAPI := &MockUploadAPIClient{
			uploadSecretFunc: func(ctx context.Context, in models.UploadSecretIn) error {
				if in.Label != "test-secret" {
					t.Errorf("Expected label 'test-secret', got '%s'", in.Label)
				}
				if in.Type != models.SecretTypeText {
					t.Errorf("Expected type '%s', got '%s'", models.SecretTypeText, in.Type)
				}
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
		cmd := NewUploadCommand(mockAPI, crypto)

		result, err := cmd.Execute(ctx, textData)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("upload API error", func(t *testing.T) {
		textData := &models.TextSecretData{
			Name:     "test-secret",
			Text:     "test text",
			Metadata: "test metadata",
		}

		expectedErr := errors.New("upload failed")
		mockAPI := &MockUploadAPIClient{
			uploadSecretFunc: func(ctx context.Context, in models.UploadSecretIn) error {
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
		cmd := NewUploadCommand(mockAPI, crypto)

		result, err := cmd.Execute(ctx, textData)

		if err == nil {
			t.Fatal("Expected error")
		}

		if err.Error() != "ошибка при загрузке секрета: "+expectedErr.Error() {
			t.Errorf("Expected error message containing 'ошибка при загрузке секрета', got '%s'", err.Error())
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("unsupported data type", func(t *testing.T) {
		mockAPI := &MockUploadAPIClient{}
		crypto := crypto.NewCrypto()
		cmd := NewUploadCommand(mockAPI, crypto)

		result, err := cmd.Execute(ctx, "unsupported data")

		if err == nil {
			t.Fatal("Expected error")
		}

		if err.Error() != "неподдерживаемый тип данных для загрузки" {
			t.Errorf("Expected error message 'неподдерживаемый тип данных для загрузки', got '%s'", err.Error())
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})
}

func TestUploadCommand_GetName(t *testing.T) {
	mockAPI := &MockUploadAPIClient{}
	crypto := crypto.NewCrypto()
	cmd := NewUploadCommand(mockAPI, crypto)

	name := cmd.GetName()
	if name != "upload" {
		t.Errorf("Expected name 'upload', got '%s'", name)
	}
}

func TestUploadCommand_GetDescription(t *testing.T) {
	mockAPI := &MockUploadAPIClient{}
	crypto := crypto.NewCrypto()
	cmd := NewUploadCommand(mockAPI, crypto)

	description := cmd.GetDescription()
	expected := "Загрузить секрет на сервер"
	if description != expected {
		t.Errorf("Expected description '%s', got '%s'", expected, description)
	}
}
