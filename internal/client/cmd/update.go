package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/models"
)

// UpdateData представляет данные для обновления секрета
type UpdateData struct {
	SecretName string `json:"secret_name"`
	Version    int    `json:"version"`
	Type       string `json:"type"`
	Data       any    `json:"data"`
}

// UpdateAPIClient интерфейс для обновления секрета
type UpdateAPIClient interface {
	UpdateSecret(ctx context.Context, in models.UpdateSecretIn) error
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
	// Получаем данные для обновления
	updateData, ok := args.(*UpdateData)
	if !ok {
		return nil, fmt.Errorf("неверный тип аргументов для команды обновления")
	}

	// Подготавливаем данные для отправки
	var secretData models.SecretData
	switch data := updateData.Data.(type) {
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

	// Шифруем только конфиденциальные данные (поле Data)
	// Type и Metadata остаются открытыми для индексации и поиска
	encryptedData, err := c.crypto.EncryptString(secretData.Data)
	if err != nil {
		return nil, fmt.Errorf("ошибка при шифровании данных: %w", err)
	}

	// Выполняем обновление через API клиент
	in := models.UpdateSecretIn{
		UploadSecretIn: models.UploadSecretIn{
			Label:         secretData.Name,
			Type:          secretData.Type,
			Metadata:      secretData.Metadata,
			EncryptedData: encryptedData,
		},
		Version: updateData.Version,
	}
	err = c.apiClient.UpdateSecret(ctx, in)
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
