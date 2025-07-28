package cmd

import (
	"context"
	"fmt"
	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/models"
)

// RegisterAPIClient определяет интерфейс для API клиента, поддерживающего регистрацию.
type RegisterAPIClient interface {
	Register(ctx context.Context, in models.RegisterIn) error
}

// RegisterCommand обрабатывает команду регистрации нового пользователя.
// Выполняет регистрацию через API и настраивает криптографические ключи.
type RegisterCommand struct {
	apiClient RegisterAPIClient
	crypto    *crypto.Crypto
}

// RegisterCommand реализует интерфейс Command
var _ Command = (*RegisterCommand)(nil)

// NewRegisterCommand создает новую команду регистрации.
//
// Параметры:
//   - apiClient: API клиент для выполнения регистрации
//   - crypto: криптографический модуль для работы с паролями и ключами
//
// Возвращает:
//   - *RegisterCommand: новый экземпляр команды регистрации
func NewRegisterCommand(apiClient RegisterAPIClient, crypto *crypto.Crypto) *RegisterCommand {
	return &RegisterCommand{
		apiClient: apiClient,
		crypto:    crypto,
	}
}

// Execute выполняет процесс регистрации пользователя.
// Генерирует соль, хеширует пароль и отправляет данные на сервер.
// После успешной регистрации создает ключ шифрования.
//
// Параметры:
//   - ctx: контекст выполнения
//   - args: должен быть типа models.UserCredentials с логином и паролем
//
// Возвращает:
//   - any: nil при успешной регистрации
//   - error: ошибка в случае неудачи
func (c *RegisterCommand) Execute(ctx context.Context, args any) (any, error) {
	credentials, ok := args.(models.UserCredentials)
	if !ok {
		return nil, fmt.Errorf("неверный тип аргументов для команды регистрации")
	}

	salt := c.crypto.GenerateAndSetSalt()
	passwordHash, err := c.crypto.GeneratePasswordHash(credentials.Password)
	if err != nil {
		return nil, fmt.Errorf("ошибка при генерации хеша пароля: %w", err)
	}

	// Создаем запрос для регистрации
	requestBody := models.RegisterIn{
		Login:        credentials.Login,
		PasswordHash: passwordHash,
		Salt:         salt,
	}

	// Выполняем регистрацию через API клиент
	err = c.apiClient.Register(ctx, requestBody)
	if err != nil {
		return nil, fmt.Errorf("ошибка при регистрации: %w", err)
	}

	// Регистрация выполнена успешно, сохраняем encryptionKey
	err = c.crypto.GenerateAndStoreEncryptionKey(credentials.Password)
	if err != nil {
		return nil, fmt.Errorf("ошибка при генерации ключа шифрования: %w", err)
	}

	return nil, nil
}

// GetName возвращает имя команды.
//
// Возвращает:
//   - string: "register"
func (c *RegisterCommand) GetName() string {
	return "register"
}

// GetDescription возвращает описание команды.
//
// Возвращает:
//   - string: описание команды регистрации
func (c *RegisterCommand) GetDescription() string {
	return "Регистрация нового пользователя"
}
