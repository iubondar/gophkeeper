package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/models"
)

// UpdateCryptoEncryptor интерфейс для шифрования данных при обновлении секретов.
type UpdateCryptoEncryptor interface {
	EncryptString(plaintext string) ([]byte, error)
}

// UpdateData представляет данные для обновления секрета.
// Содержит имя секрета, версию, тип и новые данные.
type UpdateData struct {
	SecretName string `json:"secret_name"` // Имя секрета для обновления
	Version    int    `json:"version"`     // Текущая версия секрета
	Type       string `json:"type"`        // Тип секрета
	Data       any    `json:"data"`        // Новые данные секрета
}

// UpdateAPIClient интерфейс для обновления секретов на сервере.
type UpdateAPIClient interface {
	UpdateSecret(ctx context.Context, in models.UpdateSecretIn) error
}

// UpdateCommand представляет команду обновления секретов.
// Поддерживает обновление различных типов секретов: текстовых, логинов/паролей,
// данных карт и файлов.
type UpdateCommand struct {
	apiClient UpdateAPIClient
	crypto    UpdateCryptoEncryptor
}

// NewUpdateCommand создает новую команду обновления секретов.
//
// Параметры:
//   - apiClient: API клиент для обновления секретов на сервере
//   - crypto: криптографический модуль для шифрования данных
//
// Возвращает:
//   - *UpdateCommand: новый экземпляр команды обновления
func NewUpdateCommand(apiClient UpdateAPIClient, crypto UpdateCryptoEncryptor) *UpdateCommand {
	return &UpdateCommand{
		apiClient: apiClient,
		crypto:    crypto,
	}
}

// Execute выполняет команду обновления секрета.
// Поддерживает следующие типы данных:
// - TextSecretData: текстовые секреты
// - LoginPasswordData: логины и пароли
// - CardData: данные банковских карт
// - FileData: файлы
//
// Параметры:
//   - ctx: контекст выполнения
//   - args: должен быть типа *UpdateData с данными для обновления
//
// Возвращает:
//   - any: nil при успешном обновлении
//   - error: ошибка в случае неудачи
func (c *UpdateCommand) Execute(ctx context.Context, args any) (any, error) {
	// Получаем данные для обновления
	updateData, ok := args.(*UpdateData)
	if !ok {
		return nil, fmt.Errorf("неверный тип аргументов для команды обновления")
	}

	// Подготавливаем данные для отправки
	var secretData models.SecretData
	switch data := updateData.Data.(type) {
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

	// Шифруем только конфиденциальные данные (поле Data)
	// Type и Metadata остаются открытыми для индексации и поиска
	encryptedData, err := c.crypto.EncryptString(secretData.Data)
	if err != nil {
		return nil, fmt.Errorf("ошибка при шифровании данных: %w", err)
	}

	// Выполняем обновление через API клиент
	in := models.UpdateSecretIn{
		UploadSecretIn: models.UploadSecretIn{
			Label:         secretData.Name,
			Type:          secretData.Type,
			Metadata:      secretData.Metadata,
			EncryptedData: encryptedData,
		},
		Version: updateData.Version,
	}
	err = c.apiClient.UpdateSecret(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("ошибка при обновлении секрета: %w", err)
	}

	return nil, nil
}

// GetName возвращает имя команды.
//
// Возвращает:
//   - string: "update"
func (c *UpdateCommand) GetName() string {
	return "update"
}

// GetDescription возвращает описание команды.
//
// Возвращает:
//   - string: описание команды обновления
func (c *UpdateCommand) GetDescription() string {
	return "Обновить секрет на сервере"
}
