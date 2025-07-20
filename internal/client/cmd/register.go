package cmd

import (
	"context"
	"fmt"
	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/models"
)

type RegisterAPIClient interface {
	Register(ctx context.Context, in models.RegisterIn) error
}

// RegisterCommand обрабатывает команду регистрации
type RegisterCommand struct {
	apiClient RegisterAPIClient
	crypto    *crypto.Crypto
}

// RegisterCommand реализует интерфейс Command
var _ Command = (*RegisterCommand)(nil)

// NewRegisterCommand создает новую команду регистрации
func NewRegisterCommand(apiClient RegisterAPIClient, crypto *crypto.Crypto) *RegisterCommand {
	return &RegisterCommand{
		apiClient: apiClient,
		crypto:    crypto,
	}
}

// Execute выполняет процесс регистрации пользователя
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

	return nil, nil
}

// GetName возвращает имя команды
func (c *RegisterCommand) GetName() string {
	return "register"
}

// GetDescription возвращает описание команды
func (c *RegisterCommand) GetDescription() string {
	return "Регистрация нового пользователя"
}
