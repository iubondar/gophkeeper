package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/models"
)

// UploadAPIClient интерфейс для загрузки секрета
type UploadAPIClient interface {
	UploadSecret(ctx context.Context, in models.UploadSecretIn) error
}

// UploadCommand представляет команду загрузки секрета
type UploadCommand struct {
	apiClient UploadAPIClient
	crypto    *crypto.Crypto
}

// NewUploadCommand создает новую команду загрузки
func NewUploadCommand(apiClient UploadAPIClient, crypto *crypto.Crypto) *UploadCommand {
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

// GetName возвращает имя команды
func (c *UploadCommand) GetName() string {
	return "upload"
}

// GetDescription возвращает описание команды
func (c *UploadCommand) GetDescription() string {
	return "Загрузить секрет на сервер"
}
