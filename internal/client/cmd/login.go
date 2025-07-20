package cmd

import (
	"context"
	"fmt"
	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/models"
)

type LoginAPIClient interface {
	Login(ctx context.Context, in models.LoginIn) (salt string, err error)
	Authenticate(ctx context.Context, in models.AuthenticateIn) error
}

// LoginCommand обрабатывает команду входа
type LoginCommand struct {
	apiClient LoginAPIClient
	crypto    *crypto.Crypto
}

// LoginCommand реализует интерфейс Command
var _ Command = (*LoginCommand)(nil)

// NewLoginCommand создает новую команду входа
func NewLoginCommand(apiClient LoginAPIClient, crypto *crypto.Crypto) *LoginCommand {
	return &LoginCommand{
		apiClient: apiClient,
		crypto:    crypto,
	}
}

// Execute выполняет процесс входа пользователя
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

	return nil, nil
}

// GetName возвращает имя команды
func (c *LoginCommand) GetName() string {
	return "login"
}

// GetDescription возвращает описание команды
func (c *LoginCommand) GetDescription() string {
	return "Вход в существующий аккаунт"
}
