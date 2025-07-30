package cmd

import (
	"context"
	"fmt"
	"gophkeeper/internal/client/crypto"
)

// DeleteAPIClient интерфейс для удаления секретов с сервера.
type DeleteAPIClient interface {
	DeleteSecret(ctx context.Context, secretName string) error
}

// DeleteCommand представляет команду удаления секретов с сервера.
// Удаляет секрет любого типа (текст, логин/пароль, карта, файл) по его имени.
type DeleteCommand struct {
	apiClient DeleteAPIClient
	crypto    *crypto.Crypto
}

// NewDeleteCommand создает новую команду удаления секретов.
//
// Параметры:
//   - apiClient: API клиент для удаления секретов с сервера
//   - crypto: криптографический модуль (не используется в данной команде)
//
// Возвращает:
//   - *DeleteCommand: новый экземпляр команды удаления
func NewDeleteCommand(apiClient DeleteAPIClient, crypto *crypto.Crypto) *DeleteCommand {
	return &DeleteCommand{
		apiClient: apiClient,
		crypto:    crypto,
	}
}

// Execute выполняет команду удаления секрета.
// Удаляет секрет любого типа с сервера по его имени.
//
// Параметры:
//   - ctx: контекст выполнения
//   - args: имя секрета для удаления (string)
//
// Возвращает:
//   - any: nil при успешном удалении
//   - error: ошибка в случае неудачи
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

// GetName возвращает имя команды.
//
// Возвращает:
//   - string: "delete"
func (c *DeleteCommand) GetName() string {
	return "delete"
}

// GetDescription возвращает описание команды.
//
// Возвращает:
//   - string: описание команды удаления
func (c *DeleteCommand) GetDescription() string {
	return "Удалить секрет с сервера"
}
