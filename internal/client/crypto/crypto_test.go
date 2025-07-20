package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePasswordHash_Deterministic(t *testing.T) {
	crypto := NewCrypto()
	crypto.SetSalt("dGVzdC1zYWx0") // base64 encoded "test-salt"

	password := "test-password"

	// Генерируем хеш первый раз
	hash1, err := crypto.GeneratePasswordHash(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash1)

	// Генерируем хеш второй раз с теми же параметрами
	hash2, err := crypto.GeneratePasswordHash(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash2)

	// Хеши должны быть одинаковыми
	assert.Equal(t, hash1, hash2, "Хеши должны быть детерминированными")
}

func TestVerifyPasswordHash(t *testing.T) {
	crypto := NewCrypto()
	crypto.SetSalt("dGVzdC1zYWx0") // base64 encoded "test-salt"

	password := "test-password"

	// Генерируем хеш
	hash, err := crypto.GeneratePasswordHash(password)
	assert.NoError(t, err)

	// Проверяем правильный пароль
	valid, err := crypto.VerifyPasswordHash(password, "dGVzdC1zYWx0", hash)
	assert.NoError(t, err)
	assert.True(t, valid, "Правильный пароль должен быть валидным")

	// Проверяем неправильный пароль
	valid, err = crypto.VerifyPasswordHash("wrong-password", "dGVzdC1zYWx0", hash)
	assert.NoError(t, err)
	assert.False(t, valid, "Неправильный пароль должен быть невалидным")
}

func TestGeneratePasswordHash_NoSalt(t *testing.T) {
	crypto := NewCrypto()
	// Не устанавливаем соль

	_, err := crypto.GeneratePasswordHash("test-password")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "salt is not set")
}

func TestEncryptDecryptString(t *testing.T) {
	crypto := NewCrypto()
	crypto.SetSalt("dGVzdC1zYWx0") // base64 encoded "test-salt"
	err := crypto.GenerateAndStoreEncryptionKey("test_password")
	assert.NoError(t, err)

	original := "login:user123;password:secret456;url:example.com"

	encrypted, err := crypto.EncryptString(original)
	assert.NoError(t, err)
	assert.NotEqual(t, original, string(encrypted))

	decrypted, err := crypto.DecryptString(encrypted)
	assert.NoError(t, err)
	assert.Equal(t, original, decrypted)
}
