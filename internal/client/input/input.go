// Package input предоставляет функции для обработки пользовательского ввода.
// Включает методы для получения различных типов данных секретов,
// валидации ввода и интерактивных подсказок.
package input

import (
	"gophkeeper/internal/models"
	"gophkeeper/internal/validators"
)

// InputHandler обрабатывает пользовательский ввод для различных типов секретов.
// Предоставляет методы для получения данных от пользователя с валидацией.
type InputHandler struct{}

// NewInputHandler создает новый обработчик пользовательского ввода.
//
// Возвращает:
//   - *InputHandler: новый экземпляр обработчика ввода
func NewInputHandler() *InputHandler {
	return &InputHandler{}
}

// GetUserChoice получает выбор пользователя из консоли.
//
// Возвращает:
//   - string: выбор пользователя
func (h *InputHandler) GetUserChoice() string {
	input, _ := readLine()
	return input
}

// GetUserCredentials получает логин и пароль от пользователя с валидацией.
//
// Возвращает:
//   - *models.UserCredentials: данные пользователя
//   - error: ошибка в случае неудачи или невалидных данных
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

	if ok, err := validators.ValidateUserCredentials(credentials); !ok {
		return nil, err
	}

	return credentials, nil
}

// GetLoginCredentials получает логин и пароль для входа в систему.
// Алиас для GetUserCredentials для совместимости.
//
// Возвращает:
//   - *models.UserCredentials: данные пользователя
//   - error: ошибка в случае неудачи или невалидных данных
func (h *InputHandler) GetLoginCredentials() (*models.UserCredentials, error) {
	return h.GetUserCredentials()
}

// getSecretName получает название секрета от пользователя с валидацией.
//
// Возвращает:
//   - string: название секрета
//   - error: ошибка в случае неудачи или невалидного названия
func (h *InputHandler) getSecretName() (string, error) {
	promptSecretName()
	name, err := readLine()
	if err != nil {
		return "", err
	}

	if ok, err := validators.ValidateSecretName(name); !ok {
		return "", err
	}

	return name, nil
}

// getMetadata получает метаданные секрета от пользователя.
//
// Возвращает:
//   - string: метаданные секрета
//   - error: ошибка в случае неудачи
func (h *InputHandler) getMetadata() (string, error) {
	promptMetadata()
	metadata, err := readLine()
	if err != nil {
		return "", err
	}
	return metadata, nil
}

// getTextData получает текстовые данные секрета (базовый или обновленный).
//
// Параметры:
//   - name: название секрета
//   - isUpdate: true для обновления, false для создания
//
// Возвращает:
//   - *models.TextSecretData: данные текстового секрета
//   - error: ошибка в случае неудачи или невалидных данных
func (h *InputHandler) getTextData(name string, isUpdate bool) (*models.TextSecretData, error) {
	if isUpdate {
		promptUpdateTextSecret()
	} else {
		promptTextSecret()
	}

	text, err := readLine()
	if err != nil {
		return nil, err
	}

	metadata, err := h.getMetadata()
	if err != nil {
		return nil, err
	}

	data := &models.TextSecretData{
		Name:     name,
		Text:     text,
		Metadata: metadata,
	}

	if ok, err := validators.ValidateTextSecretData(data); !ok {
		return nil, err
	}

	return data, nil
}

// getLoginPasswordData получает данные логин/пароль (базовый или обновленный).
//
// Параметры:
//   - name: название секрета
//   - isUpdate: true для обновления, false для создания
//
// Возвращает:
//   - *models.LoginPasswordData: данные логина и пароля
//   - error: ошибка в случае неудачи или невалидных данных
func (h *InputHandler) getLoginPasswordData(name string, isUpdate bool) (*models.LoginPasswordData, error) {
	if isUpdate {
		promptUpdateLogin()
	} else {
		promptLogin()
	}
	login, err := readLine()
	if err != nil {
		return nil, err
	}

	if isUpdate {
		promptUpdatePassword()
	} else {
		promptPassword()
	}
	password, err := readLine()
	if err != nil {
		return nil, err
	}

	if isUpdate {
		promptUpdateURL()
	} else {
		promptURL()
	}
	url, err := readLine()
	if err != nil {
		return nil, err
	}

	metadata, err := h.getMetadata()
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

	if ok, err := validators.ValidateLoginPasswordData(data); !ok {
		return nil, err
	}

	return data, nil
}

// getCardData получает данные банковской карты (базовый или обновленный).
//
// Параметры:
//   - name: название секрета
//   - isUpdate: true для обновления, false для создания
//
// Возвращает:
//   - *models.CardData: данные банковской карты
//   - error: ошибка в случае неудачи или невалидных данных
func (h *InputHandler) getCardData(name string, isUpdate bool) (*models.CardData, error) {
	if isUpdate {
		promptUpdateCardNumber()
	} else {
		promptCardNumber()
	}
	number, err := readLine()
	if err != nil {
		return nil, err
	}

	if isUpdate {
		promptUpdateCardHolder()
	} else {
		promptCardHolder()
	}
	holder, err := readLine()
	if err != nil {
		return nil, err
	}

	if isUpdate {
		promptUpdateCardExpiry()
	} else {
		promptCardExpiry()
	}
	expiry, err := readLine()
	if err != nil {
		return nil, err
	}

	if isUpdate {
		promptUpdateCardCVV()
	} else {
		promptCardCVV()
	}
	cvv, err := readLine()
	if err != nil {
		return nil, err
	}

	metadata, err := h.getMetadata()
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

	if ok, err := validators.ValidateCardData(data); !ok {
		return nil, err
	}

	return data, nil
}

// getFileData получает файловые данные (базовый или обновленный).
//
// Параметры:
//   - name: название секрета
//   - isUpdate: true для обновления, false для создания
//
// Возвращает:
//   - *models.FileData: данные файла
//   - error: ошибка в случае неудачи или невалидных данных
func (h *InputHandler) getFileData(name string, isUpdate bool) (*models.FileData, error) {
	if isUpdate {
		promptUpdateFilePath()
	} else {
		promptFilePath()
	}
	filePath, err := readLine()
	if err != nil {
		return nil, err
	}

	metadata, err := h.getMetadata()
	if err != nil {
		return nil, err
	}

	data := &models.FileData{
		Name:     name,
		FilePath: filePath,
		Metadata: metadata,
	}

	if ok, err := validators.ValidateFileData(data); !ok {
		return nil, err
	}

	return data, nil
}

// Публичные методы для базового ввода данных

// GetTextData получает данные текстового секрета от пользователя.
//
// Возвращает:
//   - *models.TextSecretData: данные текстового секрета
//   - error: ошибка в случае неудачи или невалидных данных
func (h *InputHandler) GetTextData() (*models.TextSecretData, error) {
	name, err := h.getSecretName()
	if err != nil {
		return nil, err
	}
	return h.getTextData(name, false)
}

// GetLoginPasswordData получает данные логина и пароля от пользователя.
//
// Возвращает:
//   - *models.LoginPasswordData: данные логина и пароля
//   - error: ошибка в случае неудачи или невалидных данных
func (h *InputHandler) GetLoginPasswordData() (*models.LoginPasswordData, error) {
	name, err := h.getSecretName()
	if err != nil {
		return nil, err
	}
	return h.getLoginPasswordData(name, false)
}

// GetCardData получает данные банковской карты от пользователя.
//
// Возвращает:
//   - *models.CardData: данные банковской карты
//   - error: ошибка в случае неудачи или невалидных данных
func (h *InputHandler) GetCardData() (*models.CardData, error) {
	name, err := h.getSecretName()
	if err != nil {
		return nil, err
	}
	return h.getCardData(name, false)
}

// GetFileData получает данные файла от пользователя.
//
// Возвращает:
//   - *models.FileData: данные файла
//   - error: ошибка в случае неудачи или невалидных данных
func (h *InputHandler) GetFileData() (*models.FileData, error) {
	name, err := h.getSecretName()
	if err != nil {
		return nil, err
	}
	return h.getFileData(name, false)
}

// GetSecretName получает название секрета от пользователя.
//
// Возвращает:
//   - string: название секрета
//   - error: ошибка в случае неудачи или невалидного названия
func (h *InputHandler) GetSecretName() (string, error) {
	return h.getSecretName()
}

// Публичные методы для обновления данных

// GetUpdatedTextData получает обновленные данные текстового секрета.
//
// Параметры:
//   - name: название секрета для обновления
//
// Возвращает:
//   - *models.TextSecretData: обновленные данные текстового секрета
//   - error: ошибка в случае неудачи или невалидных данных
func (h *InputHandler) GetUpdatedTextData(name string) (*models.TextSecretData, error) {
	return h.getTextData(name, true)
}

// GetUpdatedLoginPasswordData получает обновленные данные логина и пароля.
//
// Параметры:
//   - name: название секрета для обновления
//
// Возвращает:
//   - *models.LoginPasswordData: обновленные данные логина и пароля
//   - error: ошибка в случае неудачи или невалидных данных
func (h *InputHandler) GetUpdatedLoginPasswordData(name string) (*models.LoginPasswordData, error) {
	return h.getLoginPasswordData(name, true)
}

// GetUpdatedCardData получает обновленные данные банковской карты.
//
// Параметры:
//   - name: название секрета для обновления
//
// Возвращает:
//   - *models.CardData: обновленные данные банковской карты
//   - error: ошибка в случае неудачи или невалидных данных
func (h *InputHandler) GetUpdatedCardData(name string) (*models.CardData, error) {
	return h.getCardData(name, true)
}

// GetUpdatedFileData получает обновленные данные файла.
//
// Параметры:
//   - name: название секрета для обновления
//
// Возвращает:
//   - *models.FileData: обновленные данные файла
//   - error: ошибка в случае неудачи или невалидных данных
func (h *InputHandler) GetUpdatedFileData(name string) (*models.FileData, error) {
	return h.getFileData(name, true)
}

// GetSecretInfoForUpdate получает название секрета для обновления.
//
// Возвращает:
//   - string: название секрета для обновления
//   - error: ошибка в случае неудачи или невалидного названия
func (h *InputHandler) GetSecretInfoForUpdate() (string, error) {
	return h.getSecretName()
}

// PromptEnterNewData запрашивает ввод новых данных от пользователя.
func (h *InputHandler) PromptEnterNewData() {
	promptEnterNewData()
}
