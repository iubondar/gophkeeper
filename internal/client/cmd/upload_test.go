package cmd

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"gophkeeper/internal/models"
)

// MockUploadCryptoEncryptor представляет мок для UploadCryptoEncryptor
type MockUploadCryptoEncryptor struct {
	encryptStringFunc func(plaintext string) ([]byte, error)
	encryptStreamFunc func(writer io.Writer) (io.WriteCloser, error)
}

func (m *MockUploadCryptoEncryptor) EncryptString(plaintext string) ([]byte, error) {
	if m.encryptStringFunc != nil {
		return m.encryptStringFunc(plaintext)
	}
	return []byte("encrypted"), nil
}

func (m *MockUploadCryptoEncryptor) EncryptStream(writer io.Writer) (io.WriteCloser, error) {
	if m.encryptStreamFunc != nil {
		return m.encryptStreamFunc(writer)
	}
	// Возвращаем простой WriteCloser который просто копирует данные
	return &mockWriteCloser{writer: writer}, nil
}

// mockWriteCloser - простая реализация io.WriteCloser для тестов
type mockWriteCloser struct {
	writer io.Writer
}

func (m *mockWriteCloser) Write(p []byte) (n int, err error) {
	return m.writer.Write(p)
}

func (m *mockWriteCloser) Close() error {
	return nil
}

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
	mockCrypto := &MockUploadCryptoEncryptor{}
	cmd := NewUploadCommand(mockAPI, mockCrypto)

	if cmd == nil {
		t.Fatal("NewUploadCommand returned nil")
	}

	if cmd.apiClient != mockAPI {
		t.Error("apiClient not set correctly")
	}

	if cmd.crypto != mockCrypto {
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
		mockCrypto := &MockUploadCryptoEncryptor{
			encryptStringFunc: func(plaintext string) ([]byte, error) {
				return []byte("encrypted"), nil
			},
		}
		cmd := NewUploadCommand(mockAPI, mockCrypto)

		result, err := cmd.Execute(ctx, textData)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("successful upload login password secret", func(t *testing.T) {
		loginData := &models.LoginPasswordData{
			Name:     "test-secret",
			Login:    "test-login",
			Password: "test-password",
			URL:      "https://example.com",
			Metadata: "test metadata",
		}

		mockAPI := &MockUploadAPIClient{
			uploadSecretFunc: func(ctx context.Context, in models.UploadSecretIn) error {
				if in.Label != "test-secret" {
					t.Errorf("Expected label 'test-secret', got '%s'", in.Label)
				}
				if in.Type != models.SecretTypeLoginPassword {
					t.Errorf("Expected type '%s', got '%s'", models.SecretTypeLoginPassword, in.Type)
				}
				return nil
			},
		}
		mockCrypto := &MockUploadCryptoEncryptor{
			encryptStringFunc: func(plaintext string) ([]byte, error) {
				return []byte("encrypted"), nil
			},
		}
		cmd := NewUploadCommand(mockAPI, mockCrypto)

		result, err := cmd.Execute(ctx, loginData)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("successful upload card secret", func(t *testing.T) {
		cardData := &models.CardData{
			Name:     "test-secret",
			Number:   "1234567890123456",
			Holder:   "JOHN DOE",
			Expiry:   "12/25",
			CVV:      "123",
			Metadata: "test metadata",
		}

		mockAPI := &MockUploadAPIClient{
			uploadSecretFunc: func(ctx context.Context, in models.UploadSecretIn) error {
				if in.Label != "test-secret" {
					t.Errorf("Expected label 'test-secret', got '%s'", in.Label)
				}
				if in.Type != models.SecretTypeCard {
					t.Errorf("Expected type '%s', got '%s'", models.SecretTypeCard, in.Type)
				}
				return nil
			},
		}
		mockCrypto := &MockUploadCryptoEncryptor{
			encryptStringFunc: func(plaintext string) ([]byte, error) {
				return []byte("encrypted"), nil
			},
		}
		cmd := NewUploadCommand(mockAPI, mockCrypto)

		result, err := cmd.Execute(ctx, cardData)

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
		mockCrypto := &MockUploadCryptoEncryptor{
			encryptStringFunc: func(plaintext string) ([]byte, error) {
				return []byte("encrypted"), nil
			},
		}
		cmd := NewUploadCommand(mockAPI, mockCrypto)

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
		mockCrypto := &MockUploadCryptoEncryptor{}
		cmd := NewUploadCommand(mockAPI, mockCrypto)

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
	mockCrypto := &MockUploadCryptoEncryptor{}
	cmd := NewUploadCommand(mockAPI, mockCrypto)

	name := cmd.GetName()
	if name != "upload" {
		t.Errorf("Expected name 'upload', got '%s'", name)
	}
}

func TestUploadCommand_GetDescription(t *testing.T) {
	mockAPI := &MockUploadAPIClient{}
	mockCrypto := &MockUploadCryptoEncryptor{}
	cmd := NewUploadCommand(mockAPI, mockCrypto)

	description := cmd.GetDescription()
	expected := "Загрузить секрет на сервер"
	if description != expected {
		t.Errorf("Expected description '%s', got '%s'", expected, description)
	}
}

func TestUploadCommand_Execute_FileUpload(t *testing.T) {
	ctx := context.Background()

	// Создаем временный файл для теста
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Записываем тестовые данные в файл
	testContent := "test file content"
	_, err = tmpFile.WriteString(testContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	fileData := &models.FileData{
		Name:     "test-file",
		FilePath: tmpFile.Name(),
		Metadata: "test metadata",
	}

	mockAPI := &MockUploadAPIClient{
		uploadFileFunc: func(ctx context.Context, label, metadata string, file io.Reader, filename string) (*models.UploadSecretOut, error) {
			if label != "test-file" {
				t.Errorf("Expected label 'test-file', got '%s'", label)
			}
			if metadata != "test metadata" {
				t.Errorf("Expected metadata 'test metadata', got '%s'", metadata)
			}
			// Проверяем только расширение файла
			if !strings.HasSuffix(filename, ".txt") {
				t.Errorf("Expected filename to end with '.txt', got '%s'", filename)
			}
			return &models.UploadSecretOut{
				ID:      "test-id",
				Version: 1,
			}, nil
		},
	}

	mockCrypto := &MockUploadCryptoEncryptor{}
	cmd := NewUploadCommand(mockAPI, mockCrypto)

	result, err := cmd.Execute(ctx, fileData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	// Проверяем, что результат содержит UploadSecretOut
	uploadResult, ok := result.(*models.UploadSecretOut)
	if !ok {
		t.Fatal("expected result to be *models.UploadSecretOut")
	}

	if uploadResult.ID != "test-id" {
		t.Errorf("expected ID 'test-id', got '%s'", uploadResult.ID)
	}

	if uploadResult.Version != 1 {
		t.Errorf("expected version 1, got %d", uploadResult.Version)
	}
}

func TestUploadCommand_Execute_FileUploadError(t *testing.T) {
	ctx := context.Background()

	// Создаем временный файл для теста
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Записываем тестовые данные в файл
	testContent := "test file content"
	_, err = tmpFile.WriteString(testContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	fileData := &models.FileData{
		Name:     "test-file",
		FilePath: tmpFile.Name(),
		Metadata: "test metadata",
	}

	expectedErr := errors.New("upload failed")
	mockAPI := &MockUploadAPIClient{
		uploadFileFunc: func(ctx context.Context, label, metadata string, file io.Reader, filename string) (*models.UploadSecretOut, error) {
			return nil, expectedErr
		},
	}

	mockCrypto := &MockUploadCryptoEncryptor{}
	cmd := NewUploadCommand(mockAPI, mockCrypto)

	result, err := cmd.Execute(ctx, fileData)
	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "ошибка при загрузке зашифрованного файла: "+expectedErr.Error() {
		t.Errorf("Expected error message containing 'ошибка при загрузке зашифрованного файла', got '%s'", err.Error())
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}
