package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/models"
)

// ShowCryptoDecryptor интерфейс для расшифровки данных при отображении секретов.
type ShowCryptoDecryptor interface {
	DecryptString(encryptedData []byte) (string, error)
}

// ShowSecretResult представляет результат выполнения команды show.
// Содержит расшифрованные данные секрета, его тип, метаданные и версию.
type ShowSecretResult struct {
	Type     string `json:"type"`     // Тип секрета
	Data     any    `json:"data"`     // Расшифрованные данные секрета
	Metadata string `json:"metadata"` // Метаданные секрета
	Version  int    `json:"version"`  // Версия секрета
}

// ShowAPIClient интерфейс для получения секретов с сервера.
type ShowAPIClient interface {
	GetSecret(ctx context.Context, secretName string) (*models.GetSecretOut, error)
	GetSecretVersion(ctx context.Context, secretName string) (int, error)
}

// ShowCommand представляет команду отображения секретов.
// Получает секрет с сервера, расшифровывает его и возвращает структурированные данные.
type ShowCommand struct {
	apiClient ShowAPIClient
	crypto    ShowCryptoDecryptor
}

// NewShowCommand создает новую команду отображения секретов.
//
// Параметры:
//   - apiClient: API клиент для получения секретов с сервера
//   - crypto: криптографический модуль для расшифровки данных
//
// Возвращает:
//   - *ShowCommand: новый экземпляр команды отображения
func NewShowCommand(apiClient ShowAPIClient, crypto ShowCryptoDecryptor) *ShowCommand {
	return &ShowCommand{
		apiClient: apiClient,
		crypto:    crypto,
	}
}

// Execute выполняет команду отображения секрета.
// Получает секрет с сервера, расшифровывает его и возвращает структурированные данные
// в зависимости от типа секрета (текст, логин/пароль, карта, файл).
//
// Параметры:
//   - ctx: контекст выполнения
//   - args: имя секрета (string)
//
// Возвращает:
//   - any: ShowSecretResult с расшифрованными данными секрета
//   - error: ошибка в случае неудачи
func (c *ShowCommand) Execute(ctx context.Context, args any) (any, error) {
	// Получаем название секрета
	secretName, ok := args.(string)
	if !ok {
		return nil, fmt.Errorf("неверный тип аргументов для команды отображения")
	}

	// Получаем информацию о секрете с сервера
	secret, err := c.apiClient.GetSecret(ctx, secretName)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении информации о секрете: %w", err)
	}

	// Получаем версию секрета
	version, err := c.apiClient.GetSecretVersion(ctx, secretName)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении версии секрета: %w", err)
	}

	// Для файлов создаем структуру с информацией о файле
	// Файлы не шифруются в EncryptedData, а хранятся в файловом хранилище
	if secret.Type == models.SecretTypeFile {
		// Используем оригинальное имя файла из БД, или fallback на label
		originalName := secret.FileName
		if originalName == "" {
			originalName = secret.Label
		}

		fileData := &models.FileData{
			Name:     secret.Label,
			FilePath: originalName, // Используем оригинальное имя как FilePath для отображения
			Metadata: secret.Metadata,
		}
		return &ShowSecretResult{
			Type:     models.SecretTypeFile,
			Data:     fileData,
			Metadata: secret.Metadata,
			Version:  version,
		}, nil
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
		return handleShowTextSecret(decryptedData, secret.Metadata, version)

	case models.SecretTypeLoginPassword:
		return handleShowLoginPasswordSecret(decryptedData, secret.Metadata, version)

	case models.SecretTypeCard:
		return handleShowCardSecret(decryptedData, secret.Metadata, version)
	}

	return nil, fmt.Errorf("неподдерживаемый тип секрета: %s", secret.Type)
}

// GetName возвращает имя команды.
//
// Возвращает:
//   - string: "show"
func (c *ShowCommand) GetName() string {
	return "show"
}

// GetDescription возвращает описание команды.
//
// Возвращает:
//   - string: описание команды отображения
func (c *ShowCommand) GetDescription() string {
	return "Показать текущие данные секрета"
}

// Вынесенные приватные функции для обработки каждого типа секрета

// handleShowTextSecret обрабатывает отображение текстового секрета.
//
// Параметры:
//   - decryptedData: расшифрованные данные секрета
//   - metadata: метаданные секрета
//   - version: версия секрета
//
// Возвращает:
//   - *ShowSecretResult: результат с текстовыми данными
//   - error: ошибка в случае неудачи
func handleShowTextSecret(decryptedData string, metadata string, version int) (*ShowSecretResult, error) {
	var textSecret models.TextSecretData
	if err := json.Unmarshal([]byte(decryptedData), &textSecret); err != nil {
		return nil, fmt.Errorf("ошибка при разборе текстового секрета: %w", err)
	}
	return &ShowSecretResult{
		Type:     models.SecretTypeText,
		Data:     &textSecret,
		Metadata: metadata,
		Version:  version,
	}, nil
}

// handleShowLoginPasswordSecret обрабатывает отображение секрета логин/пароль.
//
// Параметры:
//   - decryptedData: расшифрованные данные секрета
//   - metadata: метаданные секрета
//   - version: версия секрета
//
// Возвращает:
//   - *ShowSecretResult: результат с данными логина/пароля
//   - error: ошибка в случае неудачи
func handleShowLoginPasswordSecret(decryptedData string, metadata string, version int) (*ShowSecretResult, error) {
	var loginPassword models.LoginPasswordData
	if err := json.Unmarshal([]byte(decryptedData), &loginPassword); err != nil {
		return nil, fmt.Errorf("ошибка при разборе секрета логин/пароль: %w", err)
	}
	return &ShowSecretResult{
		Type:     models.SecretTypeLoginPassword,
		Data:     &loginPassword,
		Metadata: metadata,
		Version:  version,
	}, nil
}

// handleShowCardSecret обрабатывает отображение данных банковской карты.
//
// Параметры:
//   - decryptedData: расшифрованные данные секрета
//   - metadata: метаданные секрета
//   - version: версия секрета
//
// Возвращает:
//   - *ShowSecretResult: результат с данными карты
//   - error: ошибка в случае неудачи
func handleShowCardSecret(decryptedData string, metadata string, version int) (*ShowSecretResult, error) {
	var cardData models.CardData
	if err := json.Unmarshal([]byte(decryptedData), &cardData); err != nil {
		return nil, fmt.Errorf("ошибка при разборе данных карты: %w", err)
	}
	return &ShowSecretResult{
		Type:     models.SecretTypeCard,
		Data:     &cardData,
		Metadata: metadata,
		Version:  version,
	}, nil
}
