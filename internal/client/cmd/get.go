package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/models"
)

// GetSecretResult представляет результат выполнения команды get
type GetSecretResult struct {
	Type     string `json:"type"`
	Data     any    `json:"data"`
	Metadata string `json:"metadata"`
}

type GetAPIClient interface {
	GetSecret(ctx context.Context, secretName string) (*models.GetSecretOut, error)
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

	// Расшифровываем данные секрета только для поддерживаемых типов
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

	case models.SecretTypeFile:
		return handleFileSecret(decryptedData, secret.Metadata)
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

func handleFileSecret(decryptedData string, metadata string) (*GetSecretResult, error) {
	var fileData models.FileData
	if err := json.Unmarshal([]byte(decryptedData), &fileData); err != nil {
		return nil, fmt.Errorf("ошибка при разборе данных файла: %w", err)
	}
	return &GetSecretResult{
		Type:     models.SecretTypeFile,
		Data:     &fileData,
		Metadata: metadata,
	}, nil
}
