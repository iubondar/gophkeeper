package cmd

import (
	"context"
	"fmt"
	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/models"
)

// UpdateCommand представляет команду обновления секрета
type UpdateCommand struct {
	apiClient GophKeeperClient
	crypto    *crypto.Crypto
}

// NewUpdateCommand создает новую команду обновления
func NewUpdateCommand(apiClient GophKeeperClient, crypto *crypto.Crypto) *UpdateCommand {
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
		secretData = models.SecretData{
			Name:     data.Name,
			Type:     models.SecretTypeText,
			Data:     data.Text,
			Metadata: data.Metadata,
		}

	case *models.LoginPasswordData:
		secretData = models.SecretData{
			Name:     data.Name,
			Type:     models.SecretTypeLoginPassword,
			Data:     fmt.Sprintf("login:%s;password:%s;url:%s", data.Login, data.Password, data.URL),
			Metadata: data.Metadata,
		}

	case *models.CardData:
		secretData = models.SecretData{
			Name:     data.Name,
			Type:     models.SecretTypeCard,
			Data:     fmt.Sprintf("number:%s;holder:%s;expiry:%s;cvv:%s", data.Number, data.Holder, data.Expiry, data.CVV),
			Metadata: data.Metadata,
		}

	case *models.FileData:
		secretData = models.SecretData{
			Name:     data.Name,
			Type:     models.SecretTypeFile,
			Data:     data.FilePath,
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
