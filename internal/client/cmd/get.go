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
	Type     string      `json:"type"`
	Data     interface{} `json:"data"`
	Metadata string      `json:"metadata"`
}

// GetCommand представляет команду получения секрета
type GetCommand struct {
	apiClient GophKeeperClient
	crypto    *crypto.Crypto
}

// NewGetCommand создает новую команду получения
func NewGetCommand(apiClient GophKeeperClient, crypto *crypto.Crypto) *GetCommand {
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

	// Проверяем поддерживаемые типы секретов
	switch secret.Type {
	case "text", "login_password", "card":
		// Расшифровываем данные секрета только для поддерживаемых типов
		decryptedData, err := c.crypto.DecryptString(secret.EncryptedData)
		if err != nil {
			return nil, fmt.Errorf("ошибка при расшифровке секрета: %w", err)
		}

		// Обрабатываем секрет в зависимости от его типа
		switch secret.Type {
		case "text":
			var textSecret models.TextSecretData
			if err := json.Unmarshal([]byte(decryptedData), &textSecret); err != nil {
				return nil, fmt.Errorf("ошибка при разборе текстового секрета: %w", err)
			}
			return &GetSecretResult{
				Type:     "text",
				Data:     &textSecret,
				Metadata: secret.Metadata,
			}, nil

		case "login_password":
			var loginPassword models.LoginPasswordData
			if err := json.Unmarshal([]byte(decryptedData), &loginPassword); err != nil {
				return nil, fmt.Errorf("ошибка при разборе секрета логин/пароль: %w", err)
			}
			return &GetSecretResult{
				Type:     "login_password",
				Data:     &loginPassword,
				Metadata: secret.Metadata,
			}, nil

		case "card":
			var cardData models.CardData
			if err := json.Unmarshal([]byte(decryptedData), &cardData); err != nil {
				return nil, fmt.Errorf("ошибка при разборе данных карты: %w", err)
			}
			return &GetSecretResult{
				Type:     "card",
				Data:     &cardData,
				Metadata: secret.Metadata,
			}, nil
		}

	default:
		return nil, fmt.Errorf("неподдерживаемый тип секрета: %s", secret.Type)
	}

	return nil, nil
}

// GetName возвращает имя команды
func (c *GetCommand) GetName() string {
	return "get"
}

// GetDescription возвращает описание команды
func (c *GetCommand) GetDescription() string {
	return "Получить секрет с сервера"
}
