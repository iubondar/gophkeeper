package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/models"
	"gophkeeper/internal/os_utils"
	"io"
	"os"
	"path/filepath"
)

// GetSecretResult представляет результат выполнения команды get
type GetSecretResult struct {
	Type     string `json:"type"`
	Data     any    `json:"data"`
	Metadata string `json:"metadata"`
}

type GetAPIClient interface {
	GetSecret(ctx context.Context, secretName string) (*models.GetSecretOut, error)
	DownloadFile(ctx context.Context, label string) (io.ReadCloser, error)
}

// GetCommand представляет команду получения секрета
type GetCommand struct {
	apiClient GetAPIClient
	crypto    *crypto.Crypto
}

// NewGetCommand создает новую команду получения
func NewGetCommand(apiClient GetAPIClient, crypto *crypto.Crypto) *GetCommand {
	return &GetCommand{
		apiClient: apiClient,
		crypto:    crypto,
	}
}

// Execute выполняет команду получения
func (c *GetCommand) Execute(ctx context.Context, args any) (any, error) {
	// Получаем название секрета
	secretName, ok := args.(string)
	if !ok {
		return nil, fmt.Errorf("неверный тип аргументов для команды получения")
	}

	secret, err := c.apiClient.GetSecret(ctx, secretName)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении секрета: %w", err)
	}

	// Для файлов используем специальную логику скачивания
	// Файлы не шифруются в EncryptedData, а хранятся в файловом хранилище
	if secret.Type == models.SecretTypeFile {
		return handleFileDownload(c.apiClient, c.crypto, secret.Label, secret.Metadata, secret.FileName)
	}

	// Расшифровываем данные секрета для всех остальных типов
	decryptedData, err := c.crypto.DecryptString(secret.EncryptedData)
	if err != nil {
		return nil, fmt.Errorf("ошибка при расшифровке секрета: %w", err)
	}

	// Обрабатываем секрет в зависимости от его типа
	switch secret.Type {
	case models.SecretTypeText:
		return handleTextSecret(decryptedData, secret.Metadata)

	case models.SecretTypeLoginPassword:
		return handleLoginPasswordSecret(decryptedData, secret.Metadata)

	case models.SecretTypeCard:
		return handleCardSecret(decryptedData, secret.Metadata)

	default:
		return nil, fmt.Errorf("неподдерживаемый тип секрета: %s", secret.Type)
	}
}

// GetName возвращает имя команды
func (c *GetCommand) GetName() string {
	return "get"
}

// GetDescription возвращает описание команды
func (c *GetCommand) GetDescription() string {
	return "Получить секрет с сервера"
}

// Вынесенные приватные функции для обработки каждого типа секрета

func handleTextSecret(decryptedData string, metadata string) (*GetSecretResult, error) {
	var textSecret models.TextSecretData
	if err := json.Unmarshal([]byte(decryptedData), &textSecret); err != nil {
		return nil, fmt.Errorf("ошибка при разборе текстового секрета: %w", err)
	}
	return &GetSecretResult{
		Type:     models.SecretTypeText,
		Data:     &textSecret,
		Metadata: metadata,
	}, nil
}

func handleLoginPasswordSecret(decryptedData string, metadata string) (*GetSecretResult, error) {
	var loginPassword models.LoginPasswordData
	if err := json.Unmarshal([]byte(decryptedData), &loginPassword); err != nil {
		return nil, fmt.Errorf("ошибка при разборе секрета логин/пароль: %w", err)
	}
	return &GetSecretResult{
		Type:     models.SecretTypeLoginPassword,
		Data:     &loginPassword,
		Metadata: metadata,
	}, nil
}

func handleCardSecret(decryptedData string, metadata string) (*GetSecretResult, error) {
	var cardData models.CardData
	if err := json.Unmarshal([]byte(decryptedData), &cardData); err != nil {
		return nil, fmt.Errorf("ошибка при разборе данных карты: %w", err)
	}
	return &GetSecretResult{
		Type:     models.SecretTypeCard,
		Data:     &cardData,
		Metadata: metadata,
	}, nil
}

func handleFileDownload(apiClient GetAPIClient, crypto *crypto.Crypto, label, metadata, fileName string) (*GetSecretResult, error) {
	// Скачиваем файл с сервера
	reader, err := apiClient.DownloadFile(context.Background(), label)
	if err != nil {
		return nil, fmt.Errorf("ошибка при скачивании файла: %w", err)
	}
	defer reader.Close()

	// Проверяем, является ли файл зашифрованным (по метаданным или другим признакам)
	// Пока что будем считать все файлы зашифрованными для безопасности
	isEncrypted := true

	var decryptedReader io.ReadCloser
	if isEncrypted {
		// Создаем потоковый дешифратор
		decryptedReader, err = crypto.DecryptStream(reader)
		if err != nil {
			return nil, fmt.Errorf("ошибка при создании дешифратора: %w", err)
		}
		defer decryptedReader.Close()
	} else {
		decryptedReader = reader
	}

	// Получаем путь к папке загрузок
	downloadsDir, err := os_utils.GetDownloadsDir()
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении пути к папке загрузок: %w", err)
	}

	// Используем оригинальное имя файла из БД, или fallback на label
	originalName := fileName
	if originalName == "" {
		originalName = label
	}

	// Создаем путь для сохранения файла с оригинальным именем
	filePath := filepath.Join(downloadsDir, originalName)

	// Генерируем уникальное имя файла
	filePath, err = os_utils.GenerateUniqueFilePath(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка при генерации уникального имени файла: %w", err)
	}

	// Создаем файл для записи
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании файла: %w", err)
	}
	defer file.Close()

	// Копируем данные из reader в файл (расшифрованные, если файл был зашифрован)
	_, err = io.Copy(file, decryptedReader)
	if err != nil {
		return nil, fmt.Errorf("ошибка при записи файла: %w", err)
	}

	// Создаем результат с информацией о скачанном файле
	fileData := &models.FileData{
		Name:     label,
		FilePath: filePath,
		Metadata: metadata,
	}

	return &GetSecretResult{
		Type:     models.SecretTypeFile,
		Data:     fileData,
		Metadata: metadata,
	}, nil
}
