package cmd

import (
	"context"
	"fmt"
	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/models"
)

// LoginAPIClient определяет интерфейс для API клиента, поддерживающего вход и аутентификацию.
type LoginAPIClient interface {
	Login(ctx context.Context, in models.LoginIn) (salt string, err error)
	Authenticate(ctx context.Context, in models.AuthenticateIn) error
}

// LoginCommand обрабатывает команду входа в систему.
// Выполняет двухэтапную аутентификацию: получение соли и аутентификация с хешем пароля.
type LoginCommand struct {
	apiClient LoginAPIClient
	crypto    *crypto.Crypto
}

// LoginCommand реализует интерфейс Command
var _ Command = (*LoginCommand)(nil)

// NewLoginCommand создает новую команду входа.
//
// Параметры:
//   - apiClient: API клиент для выполнения входа и аутентификации
//   - crypto: криптографический модуль для работы с паролями и ключами
//
// Возвращает:
//   - *LoginCommand: новый экземпляр команды входа
func NewLoginCommand(apiClient LoginAPIClient, crypto *crypto.Crypto) *LoginCommand {
	return &LoginCommand{
		apiClient: apiClient,
		crypto:    crypto,
	}
}

// Execute выполняет процесс входа пользователя.
// Процесс состоит из двух этапов:
// 1. Получение соли с сервера
// 2. Аутентификация с хешированным паролем
// После успешного входа создает ключ шифрования.
//
// Параметры:
//   - ctx: контекст выполнения
//   - args: должен быть типа models.UserCredentials с логином и паролем
//
// Возвращает:
//   - any: nil при успешном входе
//   - error: ошибка в случае неудачи
func (c *LoginCommand) Execute(ctx context.Context, args any) (any, error) {
	credentials, ok := args.(models.UserCredentials)
	if !ok {
		return nil, fmt.Errorf("неверный тип аргументов для команды входа")
	}

	// Создаем запрос для входа
	requestBody := models.LoginIn{
		Login: credentials.Login,
	}

	// Выполняем вход через API клиент
	salt, err := c.apiClient.Login(ctx, requestBody)
	if err != nil {
		return nil, fmt.Errorf("ошибка при входе: %w", err)
	}

	c.crypto.SetSalt(salt)
	passwordHash, err := c.crypto.GeneratePasswordHash(credentials.Password)
	if err != nil {
		return nil, fmt.Errorf("ошибка при генерации пароля: %w", err)
	}

	// Выполняем аутентификацию через API клиент
	authenticateBody := models.AuthenticateIn{
		Login:        credentials.Login,
		PasswordHash: passwordHash,
	}
	err = c.apiClient.Authenticate(ctx, authenticateBody)
	if err != nil {
		return nil, fmt.Errorf("ошибка при аутентификации: %w", err)
	}

	// Вход выполнен успешно, сохраняем encryptionKey
	err = c.crypto.GenerateAndStoreEncryptionKey(credentials.Password)
	if err != nil {
		return nil, fmt.Errorf("ошибка при генерации ключа шифрования: %w", err)
	}

	return nil, nil
}

// GetName возвращает имя команды.
//
// Возвращает:
//   - string: "login"
func (c *LoginCommand) GetName() string {
	return "login"
}

// GetDescription возвращает описание команды.
//
// Возвращает:
//   - string: описание команды входа
func (c *LoginCommand) GetDescription() string {
	return "Вход в существующий аккаунт"
}
