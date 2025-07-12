package shell

import (
	"bufio"
	"fmt"
	"gophkeeper/internal/models"
	"os"
	"strings"
)

// InputHandler обрабатывает пользовательский ввод
type InputHandler struct {
	reader *bufio.Reader
}

// NewInputHandler создает новый обработчик ввода
func NewInputHandler() *InputHandler {
	return &InputHandler{
		reader: bufio.NewReader(os.Stdin),
	}
}

// GetUserChoice получает выбор пользователя
func (h *InputHandler) GetUserChoice() string {
	choice, _ := h.reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	choice = strings.TrimSuffix(choice, "\r") // Убираем carriage return для Windows
	return choice
}

// GetUserCredentials получает логин и пароль от пользователя
func (h *InputHandler) GetUserCredentials() (*models.UserCredentials, error) {
	fmt.Print("Введите логин: ")
	login, err := h.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	login = strings.TrimSpace(login)

	fmt.Print("Введите пароль: ")
	password, err := h.reader.ReadString('\n')
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

// GetLoginCredentials получает данные для входа
func (h *InputHandler) GetLoginCredentials() (*models.UserCredentials, error) {
	return h.GetUserCredentials()
}
