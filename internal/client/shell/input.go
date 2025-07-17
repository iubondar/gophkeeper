package shell

import (
	"gophkeeper/internal/models"
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

	promptPassword()
	password, err := readLine()
	if err != nil {
		return nil, err
	}

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

	promptTextSecret()
	text, err := readLine()
	if err != nil {
		return nil, err
	}

	promptMetadata()
	metadata, err := readLine()
	if err != nil {
		return nil, err
	}

	data := &models.TextSecretData{
		Name:     name,
		Text:     text,
		Metadata: metadata,
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

	promptLogin()
	login, err := readLine()
	if err != nil {
		return nil, err
	}

	promptPassword()
	password, err := readLine()
	if err != nil {
		return nil, err
	}

	promptURL()
	url, err := readLine()
	if err != nil {
		return nil, err
	}

	promptMetadata()
	metadata, err := readLine()
	if err != nil {
		return nil, err
	}

	data := &models.LoginPasswordData{
		Name:     name,
		Login:    login,
		Password: password,
		URL:      url,
		Metadata: metadata,
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

	promptCardNumber()
	number, err := readLine()
	if err != nil {
		return nil, err
	}

	promptCardHolder()
	holder, err := readLine()
	if err != nil {
		return nil, err
	}

	promptCardExpiry()
	expiry, err := readLine()
	if err != nil {
		return nil, err
	}

	promptCardCVV()
	cvv, err := readLine()
	if err != nil {
		return nil, err
	}

	promptMetadata()
	metadata, err := readLine()
	if err != nil {
		return nil, err
	}

	data := &models.CardData{
		Name:     name,
		Number:   number,
		Holder:   holder,
		Expiry:   expiry,
		CVV:      cvv,
		Metadata: metadata,
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

	promptFilePath()
	filePath, err := readLine()
	if err != nil {
		return nil, err
	}

	promptMetadata()
	metadata, err := readLine()
	if err != nil {
		return nil, err
	}

	data := &models.FileData{
		Name:     name,
		FilePath: filePath,
		Metadata: metadata,
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

	if ok, err := models.ValidateSecretName(name); !ok {
		return "", err
	}

	return name, nil
}
