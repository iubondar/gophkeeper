package cmd

import (
	"context"
	"fmt"
	"gophkeeper/internal/client/crypto"
)

type DeleteAPIClient interface {
	DeleteSecret(ctx context.Context, secretName string) error
}

// DeleteCommand представляет команду удаления секрета
type DeleteCommand struct {
	apiClient DeleteAPIClient
	crypto    *crypto.Crypto
}

// NewDeleteCommand создает новую команду удаления
func NewDeleteCommand(apiClient DeleteAPIClient, crypto *crypto.Crypto) *DeleteCommand {
	return &DeleteCommand{
		apiClient: apiClient,
		crypto:    crypto,
	}
}

// Execute выполняет команду удаления
func (c *DeleteCommand) Execute(ctx context.Context, args any) (any, error) {
	// Получаем название секрета
	secretName, ok := args.(string)
	if !ok {
		return nil, fmt.Errorf("неверный тип аргументов для команды удаления")
	}

	// Удаляем секрет через API клиент (сервер сам определит тип и выполнит соответствующее удаление)
	err := c.apiClient.DeleteSecret(ctx, secretName)
	if err != nil {
		return nil, fmt.Errorf("ошибка при удалении секрета: %w", err)
	}

	return nil, nil
}

// GetName возвращает имя команды
func (c *DeleteCommand) GetName() string {
	return "delete"
}

// GetDescription возвращает описание команды
func (c *DeleteCommand) GetDescription() string {
	return "Удалить секрет с сервера"
}
