package cmd

import (
	"context"
	"fmt"
	"gophkeeper/internal/models"
)

type LoginAPIClient interface {
	Login(ctx context.Context, in models.LoginIn) error
}

// LoginCommand обрабатывает команду входа
type LoginCommand struct {
	apiClient LoginAPIClient
}

// LoginCommand реализует интерфейс Command
var _ Command = (*LoginCommand)(nil)

// NewLoginCommand создает новую команду входа
func NewLoginCommand(apiClient LoginAPIClient) *LoginCommand {
	return &LoginCommand{
		apiClient: apiClient,
	}
}

// Execute выполняет процесс входа пользователя
func (c *LoginCommand) Execute(ctx context.Context, args any) error {
	credentials, ok := args.(models.UserCredentials)
	if !ok {
		return fmt.Errorf("неверный тип аргументов для команды входа")
	}

	// Создаем запрос для входа
	requestBody := models.LoginIn{
		Login:        credentials.Login,
		PasswordHash: credentials.Password, // TODO: добавить хеширование пароля
	}

	// Выполняем вход через API клиент
	err := c.apiClient.Login(ctx, requestBody)
	if err != nil {
		return fmt.Errorf("ошибка при входе: %w", err)
	}

	fmt.Println("✅ Вход выполнен успешно!")
	return nil
}

// GetName возвращает имя команды
func (c *LoginCommand) GetName() string {
	return "login"
}

// GetDescription возвращает описание команды
func (c *LoginCommand) GetDescription() string {
	return "Вход в существующий аккаунт"
}
