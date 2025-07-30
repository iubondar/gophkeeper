package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/models"
	"gophkeeper/internal/os_utils"
	"io"
	"os"
	"path/filepath"
)

// CryptoDecryptor интерфейс для расшифровки данных при получении секретов.
type CryptoDecryptor interface {
	DecryptString(encryptedData []byte) (string, error)
	DecryptStream(reader io.Reader) (io.ReadCloser, error)
}

// GetSecretResult представляет результат выполнения команды get.
// Содержит расшифрованные данные секрета, его тип и метаданные.
type GetSecretResult struct {
	Type     string `json:"type"`     // Тип секрета
	Data     any    `json:"data"`     // Расшифрованные данные секрета
	Metadata string `json:"metadata"` // Метаданные секрета
}

// GetAPIClient интерфейс для получения секретов и файлов с сервера.
type GetAPIClient interface {
	GetSecret(ctx context.Context, secretName string) (*models.GetSecretOut, error)
	DownloadFile(ctx context.Context, label string) (io.ReadCloser, error)
}

// GetCommand представляет команду получения секретов с сервера.
// Получает секрет с сервера, расшифровывает его и возвращает структурированные данные.
// Для файлов также скачивает их на локальный диск.
type GetCommand struct {
	apiClient GetAPIClient
	crypto    CryptoDecryptor
}

// NewGetCommand создает новую команду получения секретов.
//
// Параметры:
//   - apiClient: API клиент для получения секретов и файлов с сервера
//   - crypto: криптографический модуль для расшифровки данных
//
// Возвращает:
//   - *GetCommand: новый экземпляр команды получения
func NewGetCommand(apiClient GetAPIClient, crypto CryptoDecryptor) *GetCommand {
	return &GetCommand{
		apiClient: apiClient,
		crypto:    crypto,
	}
}

// Execute выполняет команду получения секрета.
// Получает секрет с сервера, расшифровывает его и возвращает структурированные данные
// в зависимости от типа секрета (текст, логин/пароль, карта, файл).
// Для файлов также скачивает их в папку загрузок.
//
// Параметры:
//   - ctx: контекст выполнения
//   - args: имя секрета (string)
//
// Возвращает:
//   - any: GetSecretResult с расшифрованными данными секрета
//   - error: ошибка в случае неудачи
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
	if secret.Type == models.SecretTypeFile {
		return handleFileDownload(c.apiClient, c.crypto, secret.Label, secret.Metadata, secret.FileName)
	}

	// Проверяем поддерживаемый тип до расшифровки
	switch secret.Type {
	case models.SecretTypeText, models.SecretTypeLoginPassword, models.SecretTypeCard:
		// ok
	default:
		return nil, fmt.Errorf("неподдерживаемый тип секрета: %s", secret.Type)
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
	}

	return nil, fmt.Errorf("неподдерживаемый тип секрета: %s", secret.Type)
}

// GetName возвращает имя команды.
//
// Возвращает:
//   - string: "get"
func (c *GetCommand) GetName() string {
	return "get"
}

// GetDescription возвращает описание команды.
//
// Возвращает:
//   - string: описание команды получения
func (c *GetCommand) GetDescription() string {
	return "Получить секрет с сервера"
}

// Вынесенные приватные функции для обработки каждого типа секрета

// handleTextSecret обрабатывает получение текстового секрета.
//
// Параметры:
//   - decryptedData: расшифрованные данные секрета
//   - metadata: метаданные секрета
//
// Возвращает:
//   - *GetSecretResult: результат с текстовыми данными
//   - error: ошибка в случае неудачи
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

// handleLoginPasswordSecret обрабатывает получение секрета логин/пароль.
//
// Параметры:
//   - decryptedData: расшифрованные данные секрета
//   - metadata: метаданные секрета
//
// Возвращает:
//   - *GetSecretResult: результат с данными логина/пароля
//   - error: ошибка в случае неудачи
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

// handleCardSecret обрабатывает получение данных банковской карты.
//
// Параметры:
//   - decryptedData: расшифрованные данные секрета
//   - metadata: метаданные секрета
//
// Возвращает:
//   - *GetSecretResult: результат с данными карты
//   - error: ошибка в случае неудачи
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

// handleFileDownload обрабатывает скачивание файла с сервера.
// Скачивает зашифрованный файл, расшифровывает его и сохраняет в папку загрузок.
//
// Параметры:
//   - apiClient: API клиент для скачивания файлов
//   - crypto: криптографический модуль для расшифровки
//   - label: метка файла
//   - metadata: метаданные файла
//   - fileName: оригинальное имя файла
//
// Возвращает:
//   - *GetSecretResult: результат с информацией о скачанном файле
//   - error: ошибка в случае неудачи
func handleFileDownload(apiClient GetAPIClient, crypto CryptoDecryptor, label, metadata, fileName string) (*GetSecretResult, error) {
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
