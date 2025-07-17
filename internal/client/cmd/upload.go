package cmd

import (
	"context"
	"fmt"
	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/models"
)

// UploadCommand представляет команду загрузки секрета
type UploadCommand struct {
	apiClient GophKeeperClient
	crypto    *crypto.Crypto
}

// NewUploadCommand создает новую команду загрузки
func NewUploadCommand(apiClient GophKeeperClient, crypto *crypto.Crypto) *UploadCommand {
	return &UploadCommand{
		apiClient: apiClient,
		crypto:    crypto,
	}
}

// Execute выполняет команду загрузки
func (c *UploadCommand) Execute(ctx context.Context, args any) error {
	// Обрабатываем разные типы данных
	var secretData models.SecretData

	switch data := args.(type) {
	case *models.TextSecretData:
		secretData = models.SecretData{
			Name:     data.Name,
			Type:     "text",
			Data:     data.Text,
			Metadata: data.Metadata,
		}

	case *models.LoginPasswordData:
		secretData = models.SecretData{
			Name:     data.Name,
			Type:     "login_password",
			Data:     fmt.Sprintf("login:%s;password:%s;url:%s", data.Login, data.Password, data.URL),
			Metadata: data.Metadata,
		}

	case *models.CardData:
		secretData = models.SecretData{
			Name:     data.Name,
			Type:     "card",
			Data:     fmt.Sprintf("number:%s;holder:%s;expiry:%s;cvv:%s", data.Number, data.Holder, data.Expiry, data.CVV),
			Metadata: data.Metadata,
		}

	case *models.FileData:
		secretData = models.SecretData{
			Name:     data.Name,
			Type:     "file",
			Data:     data.FilePath,
			Metadata: data.Metadata,
		}

	default:
		return fmt.Errorf("неподдерживаемый тип данных для загрузки")
	}

	// Выполняем загрузку через API клиент
	err := c.apiClient.UploadSecret(ctx, secretData)
	if err != nil {
		return fmt.Errorf("ошибка при загрузке секрета: %w", err)
	}

	return nil
}

// GetName возвращает имя команды
func (c *UploadCommand) GetName() string {
	return "upload"
}

// GetDescription возвращает описание команды
func (c *UploadCommand) GetDescription() string {
	return "Загрузить секрет на сервер"
}
