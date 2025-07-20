package cmd

import (
	"context"
	"fmt"
	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/models"
)

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
func (c *GetCommand) Execute(ctx context.Context, args any) error {
	// Получаем название секрета
	secretName, ok := args.(string)
	if !ok {
		return fmt.Errorf("неверный тип аргументов для команды получения")
	}

	// Выполняем получение через API клиент
	secret, err := c.apiClient.GetSecret(ctx, secretName)
	if err != nil {
		return fmt.Errorf("ошибка при получении секрета: %w", err)
	}

	// Расшифровываем только конфиденциальные данные (поле Data)
	decryptedData, err := c.crypto.DecryptString(secret.EncryptedData)
	if err != nil {
		return fmt.Errorf("ошибка при расшифровке секрета: %w", err)
	}

	// Собираем полную структуру секрета из открытых и расшифрованных данных
	secretData := models.SecretData{
		Name:     secret.Label,
		Type:     secret.Type,
		Data:     decryptedData,
		Metadata: secret.Metadata,
	}

	// TODO: Здесь можно добавить логику для обработки расшифрованного секрета
	// Например, сохранить в переменную или передать дальше
	_ = secretData // Пока просто игнорируем, чтобы избежать ошибки компиляции

	return nil
}

// GetName возвращает имя команды
func (c *GetCommand) GetName() string {
	return "get"
}

// GetDescription возвращает описание команды
func (c *GetCommand) GetDescription() string {
	return "Получить секрет с сервера"
}
