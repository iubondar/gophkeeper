package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/models"
	"gophkeeper/internal/os_utils"
	"io"
	"os"
	"path/filepath"
)

// UploadCryptoEncryptor интерфейс для шифрования данных
type UploadCryptoEncryptor interface {
	EncryptString(plaintext string) ([]byte, error)
	EncryptStream(writer io.Writer) (io.WriteCloser, error)
}

// UploadAPIClient интерфейс для загрузки секрета
type UploadAPIClient interface {
	UploadSecret(ctx context.Context, in models.UploadSecretIn) error
	UploadFile(ctx context.Context, label, metadata string, file io.Reader, filename string) (*models.UploadSecretOut, error)
}

// UploadCommand представляет команду загрузки секрета
type UploadCommand struct {
	apiClient UploadAPIClient
	crypto    UploadCryptoEncryptor
}

// NewUploadCommand создает новую команду загрузки
func NewUploadCommand(apiClient UploadAPIClient, crypto UploadCryptoEncryptor) *UploadCommand {
	return &UploadCommand{
		apiClient: apiClient,
		crypto:    crypto,
	}
}

// Execute выполняет команду загрузки
func (c *UploadCommand) Execute(ctx context.Context, args any) (any, error) {
	// Обрабатываем разные типы данных
	var secretData models.SecretData

	switch data := args.(type) {
	case *models.TextSecretData:
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сериализации текстовых данных: %w", err)
		}
		secretData = models.SecretData{
			Name:     data.Name,
			Type:     models.SecretTypeText,
			Data:     string(jsonData),
			Metadata: data.Metadata,
		}

	case *models.LoginPasswordData:
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сериализации данных логина/пароля: %w", err)
		}
		secretData = models.SecretData{
			Name:     data.Name,
			Type:     models.SecretTypeLoginPassword,
			Data:     string(jsonData),
			Metadata: data.Metadata,
		}

	case *models.CardData:
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сериализации данных карты: %w", err)
		}
		secretData = models.SecretData{
			Name:     data.Name,
			Type:     models.SecretTypeCard,
			Data:     string(jsonData),
			Metadata: data.Metadata,
		}

	case *models.FileData:
		return c.handleFileUpload(ctx, data)

	default:
		return nil, fmt.Errorf("неподдерживаемый тип данных для загрузки")
	}

	// Шифруем только конфиденциальные данные (поле Data)
	// Type и Metadata остаются открытыми для индексации и поиска
	encryptedData, err := c.crypto.EncryptString(secretData.Data)
	if err != nil {
		return nil, fmt.Errorf("ошибка при шифровании данных: %w", err)
	}

	in := models.UploadSecretIn{
		Label:         secretData.Name,
		Type:          secretData.Type,
		Metadata:      secretData.Metadata,
		EncryptedData: encryptedData,
	}

	// Выполняем загрузку через API клиент
	err = c.apiClient.UploadSecret(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("ошибка при загрузке секрета: %w", err)
	}

	return nil, nil
}

// handleFileUpload обрабатывает загрузку файла
func (c *UploadCommand) handleFileUpload(ctx context.Context, data *models.FileData) (any, error) {
	// Нормализуем путь к файлу (расширяем символ ~ и другие преобразования)
	normalizedPath, err := os_utils.NormalizeFilePath(data.FilePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка при обработке пути к файлу: %w", err)
	}

	file, err := os.Open(normalizedPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии файла: %w", err)
	}
	defer file.Close()

	// Получаем оригинальное имя файла из пути
	originalName := filepath.Base(data.FilePath)

	// Создаем буфер для зашифрованных данных
	var encryptedBuffer bytes.Buffer

	// Создаем потоковый шифратор
	encryptWriter, err := c.crypto.EncryptStream(&encryptedBuffer)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании шифратора: %w", err)
	}
	defer encryptWriter.Close()

	// Копируем и шифруем файл
	_, err = io.Copy(encryptWriter, file)
	if err != nil {
		return nil, fmt.Errorf("ошибка при шифровании файла: %w", err)
	}

	// Закрываем шифратор для завершения процесса
	if err := encryptWriter.Close(); err != nil {
		return nil, fmt.Errorf("ошибка при завершении шифрования: %w", err)
	}

	// Загружаем зашифрованный файл через специальный API
	result, err := c.apiClient.UploadFile(ctx, data.Name, data.Metadata, &encryptedBuffer, originalName)
	if err != nil {
		return nil, fmt.Errorf("ошибка при загрузке зашифрованного файла: %w", err)
	}

	return result, nil
}

// GetName возвращает имя команды
func (c *UploadCommand) GetName() string {
	return "upload"
}

// GetDescription возвращает описание команды
func (c *UploadCommand) GetDescription() string {
	return "Загрузить секрет на сервер"
}
