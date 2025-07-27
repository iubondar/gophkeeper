package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/models"
)

// ShowSecretResult представляет результат выполнения команды show
type ShowSecretResult struct {
	Type     string `json:"type"`
	Data     any    `json:"data"`
	Metadata string `json:"metadata"`
	Version  int    `json:"version"`
}

// ShowAPIClient интерфейс для получения секрета
type ShowAPIClient interface {
	GetSecret(ctx context.Context, secretName string) (*models.GetSecretOut, error)
	GetSecretVersion(ctx context.Context, secretName string) (int, error)
}

// ShowCommand представляет команду отображения секрета
type ShowCommand struct {
	apiClient ShowAPIClient
	crypto    *crypto.Crypto
}

// NewShowCommand создает новую команду отображения
func NewShowCommand(apiClient ShowAPIClient, crypto *crypto.Crypto) *ShowCommand {
	return &ShowCommand{
		apiClient: apiClient,
		crypto:    crypto,
	}
}

// Execute выполняет команду отображения
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
		fileData := &models.FileData{
			Name:     secret.Label,
			FilePath: "", // Путь будет установлен при скачивании
			Metadata: secret.Metadata,
		}
		return &ShowSecretResult{
			Type:     models.SecretTypeFile,
			Data:     fileData,
			Metadata: secret.Metadata,
			Version:  version,
		}, nil
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

	default:
		return nil, fmt.Errorf("неподдерживаемый тип секрета: %s", secret.Type)
	}
}

// GetName возвращает имя команды
func (c *ShowCommand) GetName() string {
	return "show"
}

// GetDescription возвращает описание команды
func (c *ShowCommand) GetDescription() string {
	return "Показать текущие данные секрета"
}

// Вынесенные приватные функции для обработки каждого типа секрета

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
