package models

import (
	"errors"
)

// Константы для типов секретов
const (
	SecretTypeText          = "text"
	SecretTypeLoginPassword = "login_password"
	SecretTypeCard          = "card"
	SecretTypeFile          = "file"
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
	Login    string `json:"login" valid:"required~логин не может быть пустым,minstringlength(3)~логин должен быть не короче 3 символов"`
	Password string `json:"password" valid:"required~пароль не может быть пустым,minstringlength(3)~пароль должен быть не короче 3 символов"`
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
	Name     string `json:"name" valid:"required~название секрета не может быть пустым,minstringlength(3)~название секрета должно быть не короче 3 символов"`
	Text     string `json:"text" valid:"required~текст секрета не может быть пустым"`
	Metadata string `json:"metadata"`
}

// LoginPasswordData представляет данные логина и пароля
type LoginPasswordData struct {
	Name     string `json:"name" valid:"required~название секрета не может быть пустым,minstringlength(3)~название секрета должно быть не короче 3 символов"`
	Login    string `json:"login" valid:"required~логин не может быть пустым,minstringlength(3)~логин должен быть не короче 3 символов"`
	Password string `json:"password" valid:"required~пароль не может быть пустым,minstringlength(3)~пароль должен быть не короче 3 символов"`
	URL      string `json:"url"`
	Metadata string `json:"metadata"`
}

// CardData представляет данные банковской карты
type CardData struct {
	Name     string `json:"name" valid:"required~название секрета не может быть пустым,minstringlength(3)~название секрета должно быть не короче 3 символов"`
	Number   string `json:"number" valid:"required~номер карты не может быть пустым"`
	Holder   string `json:"holder" valid:"required~имя владельца не может быть пустым"`
	Expiry   string `json:"expiry" valid:"required~срок действия не может быть пустым"`
	CVV      string `json:"cvv" valid:"required~CVV не может быть пустым"`
	Metadata string `json:"metadata"`
}

// FileData представляет данные файла
type FileData struct {
	Name     string `json:"name" valid:"required~название секрета не может быть пустым,minstringlength(3)~название секрета должно быть не короче 3 символов"`
	FilePath string `json:"file_path" valid:"required~путь к файлу не может быть пустым"`
	Metadata string `json:"metadata"`
}

// UploadSecretIn представляет входные данные для загрузки секрета
// user id должен быть получен на стороне сервера из access_token
type UploadSecretIn struct {
	Label         string `json:"label"`
	Type          string `json:"type"`
	Metadata      string `json:"metadata"`
	EncryptedData []byte `json:"encrypted_data"`
	FileKey       string `json:"file_key"`
}

type UploadSecretOut struct {
	ID      string `json:"id"`
	Version int    `json:"version"`
}

type UpdateSecretIn struct {
	Label         string `json:"label"`
	Type          string `json:"type"`
	Metadata      string `json:"metadata"`
	EncryptedData []byte `json:"encrypted_data"`
	FileKey       string `json:"file_key"`
	Version       int    `json:"version"` // версия, которую ожидает клиент
}

type UpdateSecretOut struct {
	Version int `json:"version"`
}

type GetSecretOut struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	Type          string `json:"type"`
	Metadata      string `json:"metadata"`
	EncryptedData []byte `json:"encrypted_data"`
	FileKey       string `json:"file_key"`
	Version       int    `json:"version"`
}

var ErrConflict = errors.New("conflict: resource already exists")
var ErrRecordNotFound = errors.New("record not found")
