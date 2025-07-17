package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
)

const additionalData = "data"

type Crypto struct {
	salt          string
	encryptionKey []byte
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
	saltBytes, err := base64.StdEncoding.DecodeString(c.salt)
	if err != nil {
		return err
	}
	key := argon2.IDKey([]byte(password), saltBytes, 1, 64*1024, 4, 32) // 32 байта для AES-256
	c.encryptionKey = key
	return nil
}

// EncryptJSON шифрует JSON-объект с использованием AES-GCM
// и возвращает зашифрованную строку в base64
func (c *Crypto) EncryptJSON(object any) (string, error) {
	if c.encryptionKey == nil {
		return "", errors.New("encryption key is not set")
	}

	plaintext, err := json.Marshal(object)
	if err != nil {
		return "", fmt.Errorf("marshal json: %w", err)
	}

	block, err := aes.NewCipher(c.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("new cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("new GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	// Префикс: nonce || ciphertext
	full := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(full), nil
}

// DecryptJSON расшифровывает зашифрованную строку в base64 и преобразует ее в JSON-объект
// и сохраняет результат в out
func (c *Crypto) DecryptJSON(encoded string, out any) error {
	if c.encryptionKey == nil {
		return errors.New("encryption key is not set")
	}

	full, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return fmt.Errorf("decode base64: %w", err)
	}

	block, err := aes.NewCipher(c.encryptionKey)
	if err != nil {
		return fmt.Errorf("new cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("new GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(full) < nonceSize {
		return fmt.Errorf("invalid encrypted data")
	}

	nonce := full[:nonceSize]
	ciphertext := full[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decrypt: %w", err)
	}

	if err := json.Unmarshal(plaintext, out); err != nil {
		return fmt.Errorf("unmarshal json: %w", err)
	}

	return nil
}
