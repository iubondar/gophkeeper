package cmd

import (
	"context"
	"fmt"
	"gophkeeper/internal/models"
)

type RegisterAPIClient interface {
	Register(ctx context.Context, in models.RegisterIn) error
}

// RegisterCommand обрабатывает команду регистрации
type RegisterCommand struct {
	apiClient RegisterAPIClient
}

// NewRegisterCommand создает новую команду регистрации
func NewRegisterCommand(apiClient RegisterAPIClient) *RegisterCommand {
	return &RegisterCommand{
		apiClient: apiClient,
	}
}

// Execute выполняет процесс регистрации пользователя
func (c *RegisterCommand) Execute(ctx context.Context, args any) error {
	credentials, ok := args.(models.UserCredentials)
	if !ok {
		return fmt.Errorf("неверный тип аргументов для команды регистрации")
	}

	// Создаем запрос для регистрации
	requestBody := models.RegisterIn{
		Login:        credentials.Login,
		PasswordHash: credentials.Password,
		Salt:         "test-salt", // TODO: реализовать генерацию соли
	}

	// Выполняем регистрацию через API клиент
	err := c.apiClient.Register(ctx, requestBody)
	if err != nil {
		return fmt.Errorf("ошибка при регистрации: %w", err)
	}

	fmt.Println("✅ Регистрация успешна!")
	return nil
}

// GetName возвращает имя команды
func (c *RegisterCommand) GetName() string {
	return "register"
}

// GetDescription возвращает описание команды
func (c *RegisterCommand) GetDescription() string {
	return "Регистрация нового пользователя"
}
