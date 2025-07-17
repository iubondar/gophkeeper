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

func TestEncryptDecryptJSON(t *testing.T) {
	crypto := NewCrypto()
	crypto.SetSalt("dGVzdC1zYWx0") // base64 encoded "test-salt"
	// Для ключа шифрования используем тот же пароль и соль
	err := crypto.GenerateAndStoreEncryptionKey("test-password")
	assert.NoError(t, err)

	// Тестируемый объект
	type testStruct struct {
		Field1 string
		Field2 int
	}
	original := testStruct{Field1: "value", Field2: 42}

	encrypted, err := crypto.EncryptJSON(original)
	assert.NoError(t, err)
	assert.NotEmpty(t, encrypted)

	// Расшифровываем
	var decrypted testStruct
	err = crypto.DecryptJSON(encrypted, &decrypted)
	assert.NoError(t, err)
	assert.Equal(t, original, decrypted)
}

func TestDecryptJSON_InvalidKey(t *testing.T) {
	crypto1 := NewCrypto()
	crypto1.SetSalt("dGVzdC1zYWx0")
	err := crypto1.GenerateAndStoreEncryptionKey("test-password")
	assert.NoError(t, err)

	type testStruct struct {
		Field string
	}
	original := testStruct{Field: "secret"}
	encrypted, err := crypto1.EncryptJSON(original)
	assert.NoError(t, err)

	// Новый объект с другим ключом
	crypto2 := NewCrypto()
	crypto2.SetSalt("dGVzdC1zYWx0")
	err = crypto2.GenerateAndStoreEncryptionKey("wrong-password")
	assert.NoError(t, err)

	var out testStruct
	err = crypto2.DecryptJSON(encrypted, &out)
	assert.Error(t, err)
}
