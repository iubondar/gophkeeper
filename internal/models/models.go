// Package models предоставляет структуры данных для работы с GophKeeper.
// Включает модели для аутентификации, различных типов секретов и API запросов/ответов.
package models

import (
	"errors"
)

// Константы для типов секретов
const (
	// SecretTypeText - текстовый секрет
	SecretTypeText = "text"
	// SecretTypeLoginPassword - секрет с логином и паролем
	SecretTypeLoginPassword = "login_password"
	// SecretTypeCard - секрет с данными банковской карты
	SecretTypeCard = "card"
	// SecretTypeFile - секрет с файлом
	SecretTypeFile = "file"
)

// RegisterIn представляет входные данные для регистрации пользователя.
type RegisterIn struct {
	Login        string `json:"login" example:"user123"`         // Логин пользователя
	PasswordHash string `json:"password_hash" example:"hash123"` // Хеш пароля
	Salt         string `json:"salt" example:"salt123"`          // Соль для хеширования
}

// LoginIn представляет входные данные для входа пользователя.
type LoginIn struct {
	Login string `json:"login" example:"user123"` // Логин пользователя
}

// LoginOut представляет выходные данные для входа пользователя.
type LoginOut struct {
	Salt string `json:"salt" example:"salt123"` // Соль для хеширования пароля
}

// AuthenticateIn представляет входные данные для аутентификации пользователя.
type AuthenticateIn struct {
	Login        string `json:"login" example:"user123"`         // Логин пользователя
	PasswordHash string `json:"password_hash" example:"hash123"` // Хеш пароля
}

// AuthenticateOut представляет выходные данные для аутентификации пользователя.
type AuthenticateOut struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`  // Токен доступа
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // Токен обновления
}

// AuthenticateResult представляет результат аутентификации.
type AuthenticateResult struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`  // Токен доступа
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // Токен обновления
	ExpiresIn    int    `json:"expires_in" example:"3600"`                                       // Время жизни токена в секундах
	Success      bool   `json:"success" example:"true"`                                          // Успешность операции
}

// UserCredentials представляет учетные данные пользователя.
type UserCredentials struct {
	Login    string `json:"login" valid:"required~логин не может быть пустым,minstringlength(3)~логин должен быть не короче 3 символов" example:"user123"`
	Password string `json:"password" valid:"required~пароль не может быть пустым,minstringlength(3)~пароль должен быть не короче 3 символов" example:"password123"`
}

// SecretData представляет данные секрета.
type SecretData struct {
	Name     string `json:"name" example:"my_secret"`    // Название секрета
	Type     string `json:"type" example:"text"`         // Тип секрета
	Data     string `json:"data" example:"secret_data"`  // Данные секрета
	Metadata string `json:"metadata" example:"metadata"` // Метаданные
}

// TextSecretData представляет данные текстового секрета.
type TextSecretData struct {
	Name     string `json:"name" valid:"required~название секрета не может быть пустым,minstringlength(3)~название секрета должно быть не короче 3 символов" example:"my_text_secret"`
	Text     string `json:"text" valid:"required~текст секрета не может быть пустым" example:"This is my secret text"`
	Metadata string `json:"metadata" example:"Personal notes"`
}

// LoginPasswordData представляет данные логина и пароля.
type LoginPasswordData struct {
	Name     string `json:"name" valid:"required~название секрета не может быть пустым,minstringlength(3)~название секрета должно быть не короче 3 символов" example:"my_login_password"`
	Login    string `json:"login" valid:"required~логин не может быть пустым,minstringlength(3)~логин должен быть не короче 3 символов" example:"user@example.com"`
	Password string `json:"password" valid:"required~пароль не может быть пустым,minstringlength(3)~пароль должен быть не короче 3 символов" example:"mypassword123"`
	URL      string `json:"url" example:"https://example.com"`
	Metadata string `json:"metadata" example:"Work account"`
}

// CardData представляет данные банковской карты.
type CardData struct {
	Name     string `json:"name" valid:"required~название секрета не может быть пустым,minstringlength(3)~название секрета должно быть не короче 3 символов" example:"my_card"`
	Number   string `json:"number" valid:"required~номер карты не может быть пустым" example:"1234567890123456"`
	Holder   string `json:"holder" valid:"required~имя владельца не может быть пустым" example:"JOHN DOE"`
	Expiry   string `json:"expiry" valid:"required~срок действия не может быть пустым" example:"12/25"`
	CVV      string `json:"cvv" valid:"required~CVV не может быть пустым" example:"123"`
	Metadata string `json:"metadata" example:"Main credit card"`
}

// FileData представляет данные файла.
type FileData struct {
	Name     string `json:"name" valid:"required~название секрета не может быть пустым,minstringlength(3)~название секрета должно быть не короче 3 символов" example:"my_file"`
	FilePath string `json:"file_path" valid:"required~путь к файлу не может быть пустым" example:"/path/to/file.txt"`
	Metadata string `json:"metadata" example:"Important document"`
}

// UploadSecretIn представляет входные данные для загрузки секрета.
// user id должен быть получен на стороне сервера из access_token.
type UploadSecretIn struct {
	Label         string `json:"label" example:"my_secret"`          // Метка секрета
	Type          string `json:"type" example:"text"`                // Тип секрета
	Metadata      string `json:"metadata" example:"Personal secret"` // Метаданные
	EncryptedData []byte `json:"encrypted_data"`                     // Зашифрованные данные
	FileKey       string `json:"file_key" example:"file_key_123"`    // Ключ файла (для файловых секретов)
}

// UploadSecretOut представляет выходные данные для загрузки секрета.
type UploadSecretOut struct {
	ID      string `json:"id" example:"uuid-123"` // Уникальный идентификатор секрета
	Version int    `json:"version" example:"1"`   // Версия секрета
}

// UpdateSecretIn представляет входные данные для обновления секрета.
type UpdateSecretIn struct {
	UploadSecretIn
	Version int `json:"version" example:"1"` // версия, которую ожидает клиент
}

// UpdateSecretOut представляет выходные данные для обновления секрета.
type UpdateSecretOut struct {
	Version int `json:"version" example:"2"` // Новая версия секрета
}

// GetSecretOut представляет выходные данные для получения секрета.
type GetSecretOut struct {
	ID            string `json:"id" example:"uuid-123"`              // Уникальный идентификатор секрета
	Label         string `json:"label" example:"my_secret"`          // Метка секрета
	Type          string `json:"type" example:"text"`                // Тип секрета
	Metadata      string `json:"metadata" example:"Personal secret"` // Метаданные
	EncryptedData []byte `json:"encrypted_data"`                     // Зашифрованные данные
	FileKey       string `json:"file_key" example:"file_key_123"`    // Ключ файла (для файловых секретов)
	FileName      string `json:"file_name" example:"document.txt"`   // Имя файла (для файловых секретов)
	Version       int    `json:"version" example:"1"`                // Версия секрета
}

// ErrConflict возвращается при попытке создать ресурс, который уже существует.
var ErrConflict = errors.New("conflict: resource already exists")

// ErrRecordNotFound возвращается при попытке найти несуществующую запись.
var ErrRecordNotFound = errors.New("record not found")
