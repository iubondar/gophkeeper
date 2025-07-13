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

// GetTextData получает текстовые данные
func (h *InputHandler) GetTextData() (*models.TextSecretData, error) {
	fmt.Print("Введите название секрета: ")
	name, err := h.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	name = strings.TrimSpace(name)

	fmt.Print("Введите текст секрета: ")
	text, err := h.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	text = strings.TrimSpace(text)

	if name == "" {
		return nil, fmt.Errorf("название секрета не может быть пустым")
	}
	if text == "" {
		return nil, fmt.Errorf("текст секрета не может быть пустым")
	}

	return &models.TextSecretData{
		Name: name,
		Text: text,
	}, nil
}

// GetLoginPasswordData получает данные логина и пароля
func (h *InputHandler) GetLoginPasswordData() (*models.LoginPasswordData, error) {
	fmt.Print("Введите название секрета: ")
	name, err := h.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	name = strings.TrimSpace(name)

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

	fmt.Print("Введите URL (опционально): ")
	url, err := h.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	url = strings.TrimSpace(url)

	if name == "" {
		return nil, fmt.Errorf("название секрета не может быть пустым")
	}
	if login == "" {
		return nil, fmt.Errorf("логин не может быть пустым")
	}
	if password == "" {
		return nil, fmt.Errorf("пароль не может быть пустым")
	}

	return &models.LoginPasswordData{
		Name:     name,
		Login:    login,
		Password: password,
		URL:      url,
	}, nil
}

// GetCardData получает данные банковской карты
func (h *InputHandler) GetCardData() (*models.CardData, error) {
	fmt.Print("Введите название секрета: ")
	name, err := h.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	name = strings.TrimSpace(name)

	fmt.Print("Введите номер карты: ")
	number, err := h.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	number = strings.TrimSpace(number)

	fmt.Print("Введите имя владельца: ")
	holder, err := h.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	holder = strings.TrimSpace(holder)

	fmt.Print("Введите срок действия (MM/YY): ")
	expiry, err := h.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	expiry = strings.TrimSpace(expiry)

	fmt.Print("Введите CVV: ")
	cvv, err := h.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	cvv = strings.TrimSpace(cvv)

	if name == "" {
		return nil, fmt.Errorf("название секрета не может быть пустым")
	}
	if number == "" {
		return nil, fmt.Errorf("номер карты не может быть пустым")
	}
	if holder == "" {
		return nil, fmt.Errorf("имя владельца не может быть пустым")
	}
	if expiry == "" {
		return nil, fmt.Errorf("срок действия не может быть пустым")
	}
	if cvv == "" {
		return nil, fmt.Errorf("CVV не может быть пустым")
	}

	return &models.CardData{
		Name:   name,
		Number: number,
		Holder: holder,
		Expiry: expiry,
		CVV:    cvv,
	}, nil
}

// GetFileData получает данные файла
func (h *InputHandler) GetFileData() (*models.FileData, error) {
	fmt.Print("Введите название секрета: ")
	name, err := h.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	name = strings.TrimSpace(name)

	fmt.Print("Введите путь к файлу: ")
	filePath, err := h.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	filePath = strings.TrimSpace(filePath)

	if name == "" {
		return nil, fmt.Errorf("название секрета не может быть пустым")
	}
	if filePath == "" {
		return nil, fmt.Errorf("путь к файлу не может быть пустым")
	}

	return &models.FileData{
		Name:     name,
		FilePath: filePath,
	}, nil
}

// GetSecretName получает название секрета для операций get и delete
func (h *InputHandler) GetSecretName() (string, error) {
	fmt.Print("Введите название секрета: ")
	name, err := h.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	name = strings.TrimSpace(name)

	if name == "" {
		return "", fmt.Errorf("название секрета не может быть пустым")
	}

	return name, nil
}
