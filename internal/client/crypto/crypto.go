package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
)

// Crypto предоставляет методы для шифрования и хеширования данных
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
	key, err := c.calcArgon2Key(password, salt)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// calcArgon2Key вычисляет бинарный ключ шифрования с помощью Argon2
func (c *Crypto) calcArgon2Key(password, salt string) ([]byte, error) {
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return nil, err
	}
	// 32 байта для AES-256
	key := argon2.IDKey([]byte(password), saltBytes, 1, 64*1024, 4, 32)
	return key, nil
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

	// Декодируем соль из base64
	saltBytes, err := base64.StdEncoding.DecodeString(c.salt)
	if err != nil {
		return fmt.Errorf("invalid salt format: %w", err)
	}

	// Создаем новую соль для шифрования, добавляя additionalData к байтам
	encryptionSaltBytes := append(saltBytes, []byte(additionalData)...)
	encryptionSalt := base64.StdEncoding.EncodeToString(encryptionSaltBytes)

	key, err := c.calcArgon2Key(password, encryptionSalt)
	if err != nil {
		return err
	}

	c.encryptionKey = key
	return nil
}

// EncryptString шифрует строку с использованием AES-GCM
func (c *Crypto) EncryptString(plaintext string) ([]byte, error) {
	if c.encryptionKey == nil {
		return nil, errors.New("encryption key is not set")
	}

	block, err := aes.NewCipher(c.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("new GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	// Префикс: nonce || ciphertext
	return append(nonce, ciphertext...), nil
}

// DecryptString расшифровывает бинарные данные и возвращает строку
func (c *Crypto) DecryptString(encryptedData []byte) (string, error) {
	if c.encryptionKey == nil {
		return "", errors.New("encryption key is not set")
	}

	block, err := aes.NewCipher(c.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("new cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("new GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return "", fmt.Errorf("invalid encrypted data")
	}

	nonce := encryptedData[:nonceSize]
	ciphertext := encryptedData[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt: %w", err)
	}

	return string(plaintext), nil
}

// EncryptStream создает потоковый шифратор для файлов
func (c *Crypto) EncryptStream(writer io.Writer) (io.WriteCloser, error) {
	if c.encryptionKey == nil {
		return nil, errors.New("encryption key is not set")
	}

	block, err := aes.NewCipher(c.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}

	// Используем CTR режим для потокового шифрования
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("generate IV: %w", err)
	}

	// Записываем IV в начало потока
	if _, err := writer.Write(iv); err != nil {
		return nil, fmt.Errorf("write IV: %w", err)
	}

	stream := cipher.NewCTR(block, iv)

	return &encryptWriter{
		stream: stream,
		writer: writer,
	}, nil
}

// DecryptStream создает потоковый дешифратор для файлов
func (c *Crypto) DecryptStream(reader io.Reader) (io.ReadCloser, error) {
	if c.encryptionKey == nil {
		return nil, errors.New("encryption key is not set")
	}

	block, err := aes.NewCipher(c.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}

	// Читаем IV из начала потока
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(reader, iv); err != nil {
		return nil, fmt.Errorf("read IV: %w", err)
	}

	stream := cipher.NewCTR(block, iv)

	return &decryptReader{
		stream: stream,
		reader: reader,
	}, nil
}

// encryptWriter реализует потоковое шифрование
type encryptWriter struct {
	stream cipher.Stream
	writer io.Writer
}

func (ew *encryptWriter) Write(p []byte) (n int, err error) {
	// Создаем буфер для зашифрованных данных
	ciphertext := make([]byte, len(p))
	ew.stream.XORKeyStream(ciphertext, p)

	_, err = ew.writer.Write(ciphertext)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (ew *encryptWriter) Close() error {
	return nil
}

// decryptReader реализует потоковое дешифрование
type decryptReader struct {
	stream cipher.Stream
	reader io.Reader
}

func (dr *decryptReader) Read(p []byte) (n int, err error) {
	// Читаем зашифрованные данные
	n, err = dr.reader.Read(p)
	if err != nil {
		return n, err
	}

	// Расшифровываем данные на месте
	dr.stream.XORKeyStream(p[:n], p[:n])

	return n, nil
}

func (dr *decryptReader) Close() error {
	return nil
}
