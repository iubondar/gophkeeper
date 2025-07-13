package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/alexedwards/argon2id"
)

const additionalData = "data"

type Crypto struct {
	salt          string
	encryptionKey string
}

func NewCrypto() *Crypto {
	return &Crypto{}
}

// GenerateSalt генерирует случайную соль для хеширования пароля в base64
// После регистрации мы генерируем соль и устанавливаем ее в Crypto
func (c *Crypto) GenerateAndSetSalt() string {
	salt := make([]byte, 32)
	rand.Read(salt)
	c.salt = base64.StdEncoding.EncodeToString(salt)
	return c.salt
}

// SetSalt устанавливает соль для хеширования пароля
// После логина мы получаем соль с сервера и устанавливаем ее в Crypto
func (c *Crypto) SetSalt(salt string) {
	c.salt = salt
}

// GeneratePasswordHash генерирует хеш пароля с использованием соли
// Для этого мы используем соль, которая была установлена в Crypto
func (c *Crypto) GeneratePasswordHash(password string) (string, error) {
	if c.salt == "" {
		return "", errors.New("salt is not set")
	}
	hash, err := argon2id.CreateHash(password+c.salt, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

// VerifyPasswordHash проверяет, совпадает ли пароль с хешем
func (c *Crypto) VerifyPasswordHash(password string, salt string, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password+salt, hash)
}

// GenerateAndStoreEncryptionKey генерирует ключ шифрования с использованием соли
// Для этого мы используем соль, которая была установлена в Crypto
func (c *Crypto) GenerateAndStoreEncryptionKey(password string) error {
	if c.salt == "" {
		return errors.New("salt is not set")
	}
	encryptionKey, err := argon2id.CreateHash(password+c.salt+additionalData, argon2id.DefaultParams)
	if err != nil {
		return err
	}
	c.encryptionKey = encryptionKey
	return nil
}
