package cmd

import (
	"context"
	"fmt"
	"gophkeeper/internal/models"
)

type RegisterAPIClient interface {
	Register(ctx context.Context, in models.RegisterIn) error
}

// RegisterHandler обрабатывает команду регистрации
type RegisterHandler struct {
	apiClient RegisterAPIClient
}

// NewRegisterHandler создает новый обработчик регистрации
func NewRegisterHandler(apiClient RegisterAPIClient) RegisterHandler {
	return RegisterHandler{
		apiClient: apiClient,
	}
}

// Handle выполняет процесс регистрации пользователя
func (h *RegisterHandler) Run(ctx context.Context, in models.UserCredentials) error {
	// Создаем запрос для регистрации
	requestBody := models.RegisterIn{
		Login:        in.Login,
		PasswordHash: in.Password,
		Salt:         "test-salt", // TODO: реализовать генерацию соли
	}

	// Выполняем регистрацию через API клиент
	err := h.apiClient.Register(ctx, requestBody)
	if err != nil {
		return fmt.Errorf("ошибка при регистрации: %w", err)
	}

	fmt.Println("✅ Регистрация успешна!")
	return nil
}
