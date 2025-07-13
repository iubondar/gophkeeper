package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/alexedwards/argon2id"
	"golang.org/x/crypto/argon2"
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

// calcArgon2Hash вычисляет base64-хеш пароля и соли с помощью Argon2
func (c *Crypto) calcArgon2Hash(password, salt string) (string, error) {
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), saltBytes, 1, 64*1024, 4, 32)
	return base64.StdEncoding.EncodeToString(hash), nil
}

// GeneratePasswordHash генерирует хеш пароля с использованием соли (детерминированно, Argon2)
// Для этого мы используем соль, которая была установлена в Crypto
func (c *Crypto) GeneratePasswordHash(password string) (string, error) {
	if c.salt == "" {
		return "", errors.New("salt is not set")
	}
	return c.calcArgon2Hash(password, c.salt)
}

// VerifyPasswordHash проверяет, совпадает ли пароль с хешем (Argon2)
func (c *Crypto) VerifyPasswordHash(password string, salt string, hash string) (bool, error) {
	recalc, err := c.calcArgon2Hash(password, salt)
	if err != nil {
		return false, err
	}
	return recalc == hash, nil
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
