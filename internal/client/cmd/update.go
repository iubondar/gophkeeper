package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/models"
)

type UpdateAPIClient interface {
	UpdateSecret(ctx context.Context, secret models.SecretData) error
}

// UpdateCommand представляет команду обновления секрета
type UpdateCommand struct {
	apiClient UpdateAPIClient
	crypto    *crypto.Crypto
}

// NewUpdateCommand создает новую команду обновления
func NewUpdateCommand(apiClient UpdateAPIClient, crypto *crypto.Crypto) *UpdateCommand {
	return &UpdateCommand{
		apiClient: apiClient,
		crypto:    crypto,
	}
}

// Execute выполняет команду обновления
func (c *UpdateCommand) Execute(ctx context.Context, args any) (any, error) {
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
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сериализации данных файла: %w", err)
		}
		secretData = models.SecretData{
			Name:     data.Name,
			Type:     models.SecretTypeFile,
			Data:     string(jsonData),
			Metadata: data.Metadata,
		}

	default:
		return nil, fmt.Errorf("неподдерживаемый тип данных для обновления")
	}

	// Выполняем обновление через API клиент
	err := c.apiClient.UpdateSecret(ctx, secretData)
	if err != nil {
		return nil, fmt.Errorf("ошибка при обновлении секрета: %w", err)
	}

	return nil, nil
}

// GetName возвращает имя команды
func (c *UpdateCommand) GetName() string {
	return "update"
}

// GetDescription возвращает описание команды
func (c *UpdateCommand) GetDescription() string {
	return "Обновить секрет на сервере"
}
