package cmd

import (
	"context"
	"fmt"
	"gophkeeper/internal/client/crypto"
)

// DeleteCommand представляет команду удаления секрета
type DeleteCommand struct {
	apiClient GophKeeperClient
	crypto    *crypto.Crypto
}

// NewDeleteCommand создает новую команду удаления
func NewDeleteCommand(apiClient GophKeeperClient, crypto *crypto.Crypto) *DeleteCommand {
	return &DeleteCommand{
		apiClient: apiClient,
		crypto:    crypto,
	}
}

// Execute выполняет команду удаления
func (c *DeleteCommand) Execute(ctx context.Context, args any) error {
	// Получаем название секрета
	secretName, ok := args.(string)
	if !ok {
		return fmt.Errorf("неверный тип аргументов для команды удаления")
	}

	// Выполняем удаление через API клиент
	err := c.apiClient.DeleteSecret(ctx, secretName)
	if err != nil {
		return fmt.Errorf("ошибка при удалении секрета: %w", err)
	}

	return nil
}

// GetName возвращает имя команды
func (c *DeleteCommand) GetName() string {
	return "delete"
}

// GetDescription возвращает описание команды
func (c *DeleteCommand) GetDescription() string {
	return "Удалить секрет с сервера"
}
