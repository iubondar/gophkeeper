package cmd

import (
	"bufio"
	"context"
	"fmt"
	"gophkeeper/internal/models"
	"os"
	"strings"
)

type RegisterAPIClient interface {
	Register(ctx context.Context, in models.RegisterIn) error
}

// RegisterHandler обрабатывает команду регистрации
type RegisterHandler struct {
	apiClient RegisterAPIClient
}

// NewRegisterHandler создает новый обработчик регистрации
func NewRegisterHandler(apiClient RegisterAPIClient) *RegisterHandler {
	return &RegisterHandler{
		apiClient: apiClient,
	}
}

// Handle выполняет процесс регистрации пользователя
func (h *RegisterHandler) Run(ctx context.Context) error {
	fmt.Println("\n=== Регистрация ===")

	// Получаем данные от пользователя
	credentials, err := h.getUserCredentials()
	if err != nil {
		return fmt.Errorf("ошибка получения данных пользователя: %w", err)
	}

	// Создаем запрос для регистрации
	requestBody := models.RegisterIn{
		Login:        credentials.Login,
		PasswordHash: credentials.Password,
		Salt:         "test-salt", // TODO: реализовать генерацию соли
	}

	// Выполняем регистрацию через API клиент
	err = h.apiClient.Register(ctx, requestBody)
	if err != nil {
		return fmt.Errorf("ошибка при регистрации: %w", err)
	}

	fmt.Println("✅ Регистрация успешна!")
	return nil
}

// getUserCredentials получает логин и пароль от пользователя
func (h *RegisterHandler) getUserCredentials() (credentials *models.UserCredentials, err error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите логин: ")
	login, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	login = strings.TrimSpace(login)

	fmt.Print("Введите пароль: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	password = strings.TrimSpace(password)

	// Валидация входных данных
	if login == "" {
		return nil, fmt.Errorf("логин не может быть пустым")
	}
	if password == "" {
		return nil, fmt.Errorf("пароль не может быть пустым")
	}

	return &models.UserCredentials{Login: login, Password: password}, nil
}
