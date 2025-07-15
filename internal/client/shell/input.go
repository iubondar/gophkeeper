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

	credentials := &models.UserCredentials{Login: login, Password: password}

	// Валидация входных данных
	if ok, err := models.ValidateUserCredentials(credentials); !ok {
		return nil, err
	}

	return credentials, nil
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

	data := &models.TextSecretData{
		Name: name,
		Text: text,
	}

	// Валидация входных данных
	if ok, err := models.ValidateTextSecretData(data); !ok {
		return nil, err
	}

	return data, nil
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

	data := &models.LoginPasswordData{
		Name:     name,
		Login:    login,
		Password: password,
		URL:      url,
	}

	// Валидация входных данных
	if ok, err := models.ValidateLoginPasswordData(data); !ok {
		return nil, err
	}

	return data, nil
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

	data := &models.CardData{
		Name:   name,
		Number: number,
		Holder: holder,
		Expiry: expiry,
		CVV:    cvv,
	}

	// Валидация входных данных
	if ok, err := models.ValidateCardData(data); !ok {
		return nil, err
	}

	return data, nil
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

	data := &models.FileData{
		Name:     name,
		FilePath: filePath,
	}

	// Валидация входных данных
	if ok, err := models.ValidateFileData(data); !ok {
		return nil, err
	}

	return data, nil
}

// GetSecretName получает название секрета для операций get и delete
func (h *InputHandler) GetSecretName() (string, error) {
	fmt.Print("Введите название секрета: ")
	name, err := h.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	name = strings.TrimSpace(name)

	// Валидация входных данных
	if ok, err := models.ValidateSecretName(name); !ok {
		return "", err
	}

	return name, nil
}
