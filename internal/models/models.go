package models

import (
	"errors"
)

type RegisterIn struct {
	Login        string `json:"login"`
	PasswordHash string `json:"password_hash"`
	Salt         string `json:"salt"`
}

type LoginIn struct {
	Login string `json:"login"`
}

type LoginOut struct {
	Salt string `json:"salt"`
}

type AuthenticateIn struct {
	Login        string `json:"login"`
	PasswordHash string `json:"password_hash"`
}

type AuthenticateOut struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthenticateResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Success      bool   `json:"success"`
}

type UserCredentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// SecretData представляет данные секрета
type SecretData struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Data     string `json:"data"`
	Metadata string `json:"metadata"`
}

// TextSecretData представляет данные текстового секрета
type TextSecretData struct {
	Name     string `json:"name"`
	Text     string `json:"text"`
	Metadata string `json:"metadata"`
}

// LoginPasswordData представляет данные логина и пароля
type LoginPasswordData struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	URL      string `json:"url"`
	Metadata string `json:"metadata"`
}

// CardData представляет данные банковской карты
type CardData struct {
	Name     string `json:"name"`
	Number   string `json:"number"`
	Holder   string `json:"holder"`
	Expiry   string `json:"expiry"`
	CVV      string `json:"cvv"`
	Metadata string `json:"metadata"`
}

// FileData представляет данные файла
type FileData struct {
	Name     string `json:"name"`
	FilePath string `json:"file_path"`
	Metadata string `json:"metadata"`
}

// UploadSecretIn представляет входные данные для загрузки секрета
// user_id должен быть получен из авторизации (например, из токена)
type UploadSecretIn struct {
	UserID        string `json:"user_id"` // или uuid.UUID, если нужно
	Label         string `json:"label"`
	Type          string `json:"type"`
	Metadata      string `json:"metadata"`
	EncryptedData string `json:"encrypted_data"`
	FileKey       string `json:"file_key"`
	Version       int    `json:"version"`
}

type UploadSecretOut struct {
	ID      string `json:"id"`
	Version int    `json:"version"`
}

var ErrConflict = errors.New("conflict: resource already exists")
