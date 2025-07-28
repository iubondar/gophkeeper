package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"testing"

	"bytes"
	"gophkeeper/internal/models"
)

// MockCryptoDecryptor представляет мок для CryptoDecryptor
type MockCryptoDecryptor struct {
	decryptStringFunc func(encryptedData []byte) (string, error)
	decryptStreamFunc func(reader io.Reader) (io.ReadCloser, error)
}

func (m *MockCryptoDecryptor) DecryptString(encryptedData []byte) (string, error) {
	if m.decryptStringFunc != nil {
		return m.decryptStringFunc(encryptedData)
	}
	return "", nil
}

func (m *MockCryptoDecryptor) DecryptStream(reader io.Reader) (io.ReadCloser, error) {
	if m.decryptStreamFunc != nil {
		return m.decryptStreamFunc(reader)
	}
	return nil, nil
}

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
	mockCrypto := &MockCryptoDecryptor{}
	cmd := NewGetCommand(mockAPI, mockCrypto)

	if cmd == nil {
		t.Fatal("NewGetCommand returned nil")
	}

	if cmd.apiClient != mockAPI {
		t.Error("apiClient not set correctly")
	}

	if cmd.crypto != mockCrypto {
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
		mockCrypto := &MockCryptoDecryptor{}
		cmd := NewGetCommand(mockAPI, mockCrypto)

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
		mockCrypto := &MockCryptoDecryptor{}
		cmd := NewGetCommand(mockAPI, mockCrypto)

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
		mockCrypto := &MockCryptoDecryptor{}
		cmd := NewGetCommand(mockAPI, mockCrypto)

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

	t.Run("successful get text secret", func(t *testing.T) {
		secretName := "test-secret"
		textData := models.TextSecretData{Name: "n", Text: "t", Metadata: "m"}
		jsonData, _ := json.Marshal(textData)
		mockAPI := &MockGetAPIClient{
			getSecretFunc: func(ctx context.Context, secretName string) (*models.GetSecretOut, error) {
				return &models.GetSecretOut{
					Label:         secretName,
					Type:          models.SecretTypeText,
					Metadata:      "meta",
					EncryptedData: jsonData,
				}, nil
			},
		}
		mockCrypto := &MockCryptoDecryptor{
			decryptStringFunc: func(encryptedData []byte) (string, error) {
				return string(encryptedData), nil
			},
		}
		cmd := NewGetCommand(mockAPI, mockCrypto)

		result, err := cmd.Execute(ctx, secretName)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("expected result, got nil")
		}
	})

	t.Run("successful get login password secret", func(t *testing.T) {
		secretName := "test-secret"
		loginData := models.LoginPasswordData{Name: "n", Login: "l", Password: "p", URL: "u", Metadata: "m"}
		jsonData, _ := json.Marshal(loginData)
		mockAPI := &MockGetAPIClient{
			getSecretFunc: func(ctx context.Context, secretName string) (*models.GetSecretOut, error) {
				return &models.GetSecretOut{
					Label:         secretName,
					Type:          models.SecretTypeLoginPassword,
					Metadata:      "meta",
					EncryptedData: jsonData,
				}, nil
			},
		}
		mockCrypto := &MockCryptoDecryptor{
			decryptStringFunc: func(encryptedData []byte) (string, error) {
				return string(encryptedData), nil
			},
		}
		cmd := NewGetCommand(mockAPI, mockCrypto)

		result, err := cmd.Execute(ctx, secretName)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("expected result, got nil")
		}
	})

	t.Run("successful get card secret", func(t *testing.T) {
		secretName := "test-secret"
		cardData := models.CardData{Name: "n", Number: "1", Holder: "h", Expiry: "e", CVV: "c", Metadata: "m"}
		jsonData, _ := json.Marshal(cardData)
		mockAPI := &MockGetAPIClient{
			getSecretFunc: func(ctx context.Context, secretName string) (*models.GetSecretOut, error) {
				return &models.GetSecretOut{
					Label:         secretName,
					Type:          models.SecretTypeCard,
					Metadata:      "meta",
					EncryptedData: jsonData,
				}, nil
			},
		}
		mockCrypto := &MockCryptoDecryptor{
			decryptStringFunc: func(encryptedData []byte) (string, error) {
				return string(encryptedData), nil
			},
		}
		cmd := NewGetCommand(mockAPI, mockCrypto)

		result, err := cmd.Execute(ctx, secretName)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("expected result, got nil")
		}
	})

	t.Run("successful get file secret", func(t *testing.T) {
		secretName := "test-file"
		mockAPI := &MockGetAPIClient{
			getSecretFunc: func(ctx context.Context, secretName string) (*models.GetSecretOut, error) {
				return &models.GetSecretOut{
					Label:         secretName,
					Type:          models.SecretTypeFile,
					Metadata:      "meta",
					EncryptedData: []byte("file-data"),
				}, nil
			},
			downloadFileFunc: func(ctx context.Context, label string) (io.ReadCloser, error) {
				return io.NopCloser(bytes.NewReader([]byte("test file content"))), nil
			},
		}
		mockCrypto := &MockCryptoDecryptor{
			decryptStreamFunc: func(reader io.Reader) (io.ReadCloser, error) {
				return io.NopCloser(reader), nil
			},
		}
		cmd := NewGetCommand(mockAPI, mockCrypto)

		result, err := cmd.Execute(ctx, secretName)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("expected result, got nil")
		}
	})
}

func TestHandleTextSecret(t *testing.T) {
	valid := `{"name":"n","text":"t","metadata":"m"}`
	res, err := handleTextSecret(valid, "meta")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Type != models.SecretTypeText {
		t.Errorf("expected type %s, got %s", models.SecretTypeText, res.Type)
	}
	if res.Metadata != "meta" {
		t.Errorf("expected metadata 'meta', got '%s'", res.Metadata)
	}

	_, err = handleTextSecret("not json", "meta")
	if err == nil {
		t.Error("expected error for invalid json")
	}
}

func TestHandleLoginPasswordSecret(t *testing.T) {
	valid := `{"name":"n","login":"l","password":"p","url":"u","metadata":"m"}`
	res, err := handleLoginPasswordSecret(valid, "meta")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Type != models.SecretTypeLoginPassword {
		t.Errorf("expected type %s, got %s", models.SecretTypeLoginPassword, res.Type)
	}
	if res.Metadata != "meta" {
		t.Errorf("expected metadata 'meta', got '%s'", res.Metadata)
	}

	_, err = handleLoginPasswordSecret("not json", "meta")
	if err == nil {
		t.Error("expected error for invalid json")
	}
}

func TestHandleCardSecret(t *testing.T) {
	valid := `{"name":"n","number":"1","holder":"h","expiry":"e","cvv":"c","metadata":"m"}`
	res, err := handleCardSecret(valid, "meta")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Type != models.SecretTypeCard {
		t.Errorf("expected type %s, got %s", models.SecretTypeCard, res.Type)
	}
	if res.Metadata != "meta" {
		t.Errorf("expected metadata 'meta', got '%s'", res.Metadata)
	}

	_, err = handleCardSecret("not json", "meta")
	if err == nil {
		t.Error("expected error for invalid json")
	}
}

func TestHandleFileDownload(t *testing.T) {
	label := "test-file"
	metadata := "test metadata"
	fileName := "test.txt"

	// Создаем мок для API клиента
	mockAPI := &MockGetAPIClient{
		downloadFileFunc: func(ctx context.Context, label string) (io.ReadCloser, error) {
			if label != "test-file" {
				t.Errorf("Expected label 'test-file', got '%s'", label)
			}
			// Возвращаем простой reader с тестовыми данными
			return io.NopCloser(bytes.NewReader([]byte("test file content"))), nil
		},
	}

	// Создаем мок для crypto
	mockCrypto := &MockCryptoDecryptor{
		decryptStreamFunc: func(reader io.Reader) (io.ReadCloser, error) {
			// Просто возвращаем тот же reader для теста
			return io.NopCloser(reader), nil
		},
	}

	result, err := handleFileDownload(mockAPI, mockCrypto, label, metadata, fileName)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	if result.Type != models.SecretTypeFile {
		t.Errorf("expected type %s, got %s", models.SecretTypeFile, result.Type)
	}

	if result.Metadata != metadata {
		t.Errorf("expected metadata '%s', got '%s'", metadata, result.Metadata)
	}

	// Проверяем, что данные содержат FileData
	fileData, ok := result.Data.(*models.FileData)
	if !ok {
		t.Fatal("expected result.Data to be *models.FileData")
	}

	if fileData.Name != label {
		t.Errorf("expected name '%s', got '%s'", label, fileData.Name)
	}

	if fileData.Metadata != metadata {
		t.Errorf("expected metadata '%s', got '%s'", metadata, fileData.Metadata)
	}
}

func TestHandleFileDownload_APIError(t *testing.T) {
	label := "test-file"
	metadata := "test metadata"
	fileName := "test.txt"

	expectedErr := errors.New("download failed")
	mockAPI := &MockGetAPIClient{
		downloadFileFunc: func(ctx context.Context, label string) (io.ReadCloser, error) {
			return nil, expectedErr
		},
	}

	mockCrypto := &MockCryptoDecryptor{}

	result, err := handleFileDownload(mockAPI, mockCrypto, label, metadata, fileName)
	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "ошибка при скачивании файла: "+expectedErr.Error() {
		t.Errorf("Expected error message containing 'ошибка при скачивании файла', got '%s'", err.Error())
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}

func TestHandleFileDownload_DecryptError(t *testing.T) {
	label := "test-file"
	metadata := "test metadata"
	fileName := "test.txt"

	mockAPI := &MockGetAPIClient{
		downloadFileFunc: func(ctx context.Context, label string) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader([]byte("test file content"))), nil
		},
	}

	expectedErr := errors.New("decrypt failed")
	mockCrypto := &MockCryptoDecryptor{
		decryptStreamFunc: func(reader io.Reader) (io.ReadCloser, error) {
			return nil, expectedErr
		},
	}

	result, err := handleFileDownload(mockAPI, mockCrypto, label, metadata, fileName)
	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "ошибка при создании дешифратора: "+expectedErr.Error() {
		t.Errorf("Expected error message containing 'ошибка при создании дешифратора', got '%s'", err.Error())
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}

func TestGetCommand_GetName(t *testing.T) {
	mockAPI := &MockGetAPIClient{}
	mockCrypto := &MockCryptoDecryptor{}
	cmd := NewGetCommand(mockAPI, mockCrypto)

	name := cmd.GetName()
	if name != "get" {
		t.Errorf("Expected name 'get', got '%s'", name)
	}
}

func TestGetCommand_GetDescription(t *testing.T) {
	mockAPI := &MockGetAPIClient{}
	mockCrypto := &MockCryptoDecryptor{}
	cmd := NewGetCommand(mockAPI, mockCrypto)

	description := cmd.GetDescription()
	expected := "Получить секрет с сервера"
	if description != expected {
		t.Errorf("Expected description '%s', got '%s'", expected, description)
	}
}
