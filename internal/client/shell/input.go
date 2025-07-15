package shell

import (
	"gophkeeper/internal/models"
	"strings"
)

// InputHandler обрабатывает пользовательский ввод
type InputHandler struct{}

// NewInputHandler создает новый обработчик ввода
func NewInputHandler() *InputHandler {
	return &InputHandler{}
}

// GetUserChoice получает выбор пользователя
func (h *InputHandler) GetUserChoice() string {
	input, _ := readLine()
	return input
}

// GetUserCredentials получает логин и пароль от пользователя
func (h *InputHandler) GetUserCredentials() (*models.UserCredentials, error) {
	promptLogin()
	login, err := readLine()
	if err != nil {
		return nil, err
	}
	login = strings.TrimSpace(login)

	promptPassword()
	password, err := readLine()
	if err != nil {
		return nil, err
	}
	password = strings.TrimSpace(password)

	credentials := &models.UserCredentials{Login: login, Password: password}

	if ok, err := models.ValidateUserCredentials(credentials); !ok {
		return nil, err
	}

	return credentials, nil
}

func (h *InputHandler) GetLoginCredentials() (*models.UserCredentials, error) {
	return h.GetUserCredentials()
}

func (h *InputHandler) GetTextData() (*models.TextSecretData, error) {
	promptSecretName()
	name, err := readLine()
	if err != nil {
		return nil, err
	}
	name = strings.TrimSpace(name)

	promptTextSecret()
	text, err := readLine()
	if err != nil {
		return nil, err
	}
	text = strings.TrimSpace(text)

	data := &models.TextSecretData{
		Name: name,
		Text: text,
	}

	if ok, err := models.ValidateTextSecretData(data); !ok {
		return nil, err
	}

	return data, nil
}

func (h *InputHandler) GetLoginPasswordData() (*models.LoginPasswordData, error) {
	promptSecretName()
	name, err := readLine()
	if err != nil {
		return nil, err
	}
	name = strings.TrimSpace(name)

	promptLogin()
	login, err := readLine()
	if err != nil {
		return nil, err
	}
	login = strings.TrimSpace(login)

	promptPassword()
	password, err := readLine()
	if err != nil {
		return nil, err
	}
	password = strings.TrimSpace(password)

	promptURL()
	url, err := readLine()
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

	if ok, err := models.ValidateLoginPasswordData(data); !ok {
		return nil, err
	}

	return data, nil
}

func (h *InputHandler) GetCardData() (*models.CardData, error) {
	promptSecretName()
	name, err := readLine()
	if err != nil {
		return nil, err
	}
	name = strings.TrimSpace(name)

	promptCardNumber()
	number, err := readLine()
	if err != nil {
		return nil, err
	}
	number = strings.TrimSpace(number)

	promptCardHolder()
	holder, err := readLine()
	if err != nil {
		return nil, err
	}
	holder = strings.TrimSpace(holder)

	promptCardExpiry()
	expiry, err := readLine()
	if err != nil {
		return nil, err
	}
	expiry = strings.TrimSpace(expiry)

	promptCardCVV()
	cvv, err := readLine()
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

	if ok, err := models.ValidateCardData(data); !ok {
		return nil, err
	}

	return data, nil
}

func (h *InputHandler) GetFileData() (*models.FileData, error) {
	promptSecretName()
	name, err := readLine()
	if err != nil {
		return nil, err
	}
	name = strings.TrimSpace(name)

	promptFilePath()
	filePath, err := readLine()
	if err != nil {
		return nil, err
	}
	filePath = strings.TrimSpace(filePath)

	data := &models.FileData{
		Name:     name,
		FilePath: filePath,
	}

	if ok, err := models.ValidateFileData(data); !ok {
		return nil, err
	}

	return data, nil
}

func (h *InputHandler) GetSecretName() (string, error) {
	promptSecretName()
	name, err := readLine()
	if err != nil {
		return "", err
	}
	name = strings.TrimSpace(name)

	if ok, err := models.ValidateSecretName(name); !ok {
		return "", err
	}

	return name, nil
}
