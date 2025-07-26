package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/client/display"
	"gophkeeper/internal/client/input"
	"gophkeeper/internal/models"
)

// UpdateAPIClient интерфейс для обновления секрета
type UpdateAPIClient interface {
	GetSecret(ctx context.Context, secretName string) (*models.GetSecretOut, error)
	GetSecretVersion(ctx context.Context, secretName string) (int, error)
	UpdateSecret(ctx context.Context, secret models.SecretData, version int) error
}

// UpdateCommand представляет команду обновления секрета
type UpdateCommand struct {
	apiClient    UpdateAPIClient
	crypto       *crypto.Crypto
	inputHandler *input.InputHandler
}

// NewUpdateCommand создает новую команду обновления
func NewUpdateCommand(apiClient UpdateAPIClient, crypto *crypto.Crypto) *UpdateCommand {
	return &UpdateCommand{
		apiClient:    apiClient,
		crypto:       crypto,
		inputHandler: input.NewInputHandler(),
	}
}

// Execute выполняет команду обновления
func (c *UpdateCommand) Execute(ctx context.Context, args any) (any, error) {
	// Получаем название секрета
	secretName, ok := args.(string)
	if !ok {
		return nil, fmt.Errorf("неверный тип аргументов для команды обновления")
	}

	// Получаем информацию о секрете с сервера
	secret, err := c.apiClient.GetSecret(ctx, secretName)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении информации о секрете: %w", err)
	}

	// Получаем версию секрета
	version, err := c.apiClient.GetSecretVersion(ctx, secretName)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении версии секрета: %w", err)
	}

	// Расшифровываем текущие данные для отображения
	decryptedData, err := c.crypto.DecryptString(secret.EncryptedData)
	if err != nil {
		return nil, fmt.Errorf("ошибка при расшифровке секрета: %w", err)
	}

	// Отображаем текущие данные секрета
	fmt.Printf("Текущая версия секрета: %d\n", version)
	fmt.Printf("Тип секрета: %s\n", secret.Type)
	fmt.Printf("Текущие данные:\n")

	// Отображаем текущие данные в зависимости от типа
	switch secret.Type {
	case models.SecretTypeText:
		var textSecret models.TextSecretData
		if err := json.Unmarshal([]byte(decryptedData), &textSecret); err != nil {
			return nil, fmt.Errorf("ошибка при разборе текстового секрета: %w", err)
		}
		display.DisplayTextSecret(&textSecret, secret.Metadata)
	case models.SecretTypeLoginPassword:
		var loginPassword models.LoginPasswordData
		if err := json.Unmarshal([]byte(decryptedData), &loginPassword); err != nil {
			return nil, fmt.Errorf("ошибка при разборе секрета логин/пароль: %w", err)
		}
		display.DisplayLoginPassword(&loginPassword, secret.Metadata)
	case models.SecretTypeCard:
		var cardData models.CardData
		if err := json.Unmarshal([]byte(decryptedData), &cardData); err != nil {
			return nil, fmt.Errorf("ошибка при разборе данных карты: %w", err)
		}
		display.DisplayCardData(&cardData, secret.Metadata)
	case models.SecretTypeFile:
		var fileData models.FileData
		if err := json.Unmarshal([]byte(decryptedData), &fileData); err != nil {
			return nil, fmt.Errorf("ошибка при разборе данных файла: %w", err)
		}
		display.DisplayFileData(&fileData, secret.Metadata)
	default:
		return nil, fmt.Errorf("неподдерживаемый тип секрета: %s", secret.Type)
	}

	fmt.Println("Введите новые данные:")

	// Получаем обновленные данные в зависимости от типа секрета
	var updatedData any
	switch secret.Type {
	case models.SecretTypeText:
		updatedData, err = c.inputHandler.GetUpdatedTextData(secretName)
	case models.SecretTypeLoginPassword:
		updatedData, err = c.inputHandler.GetUpdatedLoginPasswordData(secretName)
	case models.SecretTypeCard:
		updatedData, err = c.inputHandler.GetUpdatedCardData(secretName)
	case models.SecretTypeFile:
		updatedData, err = c.inputHandler.GetUpdatedFileData(secretName)
	default:
		return nil, fmt.Errorf("неподдерживаемый тип секрета: %s", secret.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("ошибка при получении обновленных данных: %w", err)
	}

	// Подготавливаем данные для отправки
	var secretData models.SecretData
	switch data := updatedData.(type) {
	case *models.TextSecretData:
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сериализации текстовых данных: %w", err)
		}
		secretData = models.SecretData{
			Name:     data.Name,
			Type:     models.SecretTypeText,
			Data:     string(jsonData),
			Metadata: data.Metadata,
		}

	case *models.LoginPasswordData:
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сериализации данных логина/пароля: %w", err)
		}
		secretData = models.SecretData{
			Name:     data.Name,
			Type:     models.SecretTypeLoginPassword,
			Data:     string(jsonData),
			Metadata: data.Metadata,
		}

	case *models.CardData:
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сериализации данных карты: %w", err)
		}
		secretData = models.SecretData{
			Name:     data.Name,
			Type:     models.SecretTypeCard,
			Data:     string(jsonData),
			Metadata: data.Metadata,
		}

	case *models.FileData:
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сериализации данных файла: %w", err)
		}
		secretData = models.SecretData{
			Name:     data.Name,
			Type:     models.SecretTypeFile,
			Data:     string(jsonData),
			Metadata: data.Metadata,
		}

	default:
		return nil, fmt.Errorf("неподдерживаемый тип данных для обновления")
	}

	// Выполняем обновление через API клиент
	err = c.apiClient.UpdateSecret(ctx, secretData, version)
	if err != nil {
		return nil, fmt.Errorf("ошибка при обновлении секрета: %w", err)
	}

	return nil, nil
}

// GetName возвращает имя команды
func (c *UpdateCommand) GetName() string {
	return "update"
}

// GetDescription возвращает описание команды
func (c *UpdateCommand) GetDescription() string {
	return "Обновить секрет на сервере"
}
