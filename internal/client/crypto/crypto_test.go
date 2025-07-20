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

func TestGenerateAndStoreEncryptionKey_EncryptDecryptCycle(t *testing.T) {
	crypto := NewCrypto()
	crypto.SetSalt("dGVzdC1zYWx0") // base64 encoded "test-salt"

	password := "test_password"

	// Генерируем ключ шифрования
	err := crypto.GenerateAndStoreEncryptionKey(password)
	assert.NoError(t, err)
	assert.NotNil(t, crypto.encryptionKey)
	assert.Len(t, crypto.encryptionKey, 32, "Ключ должен быть 32 байта для AES-256")

	// Тестируем различные типы данных
	testCases := []string{
		"Простая строка",
		"Строка с кириллицей: привет мир!",
		"String with English: hello world!",
		"Строка с символами: !@#$%^&*()_+-=[]{}|;':\",./<>?",
		"Очень длинная строка " + string(make([]byte, 1000)), // 1000 символов
		"", // пустая строка
		"Строка с переносами\nстрок\nи табуляцией\t",
	}

	for _, testCase := range testCases {
		t.Run(testCase, func(t *testing.T) {
			// Шифруем данные
			encrypted, err := crypto.EncryptString(testCase)
			assert.NoError(t, err, "Шифрование должно пройти успешно")
			assert.NotEmpty(t, encrypted, "Зашифрованные данные не должны быть пустыми")
			assert.NotEqual(t, testCase, string(encrypted), "Зашифрованные данные должны отличаться от оригинала")

			// Дешифруем данные
			decrypted, err := crypto.DecryptString(encrypted)
			assert.NoError(t, err, "Дешифрование должно пройти успешно")
			assert.Equal(t, testCase, decrypted, "Дешифрованные данные должны совпадать с оригиналом")
		})
	}
}

func TestGenerateAndStoreEncryptionKey_Deterministic(t *testing.T) {
	crypto1 := NewCrypto()
	crypto2 := NewCrypto()

	salt := "dGVzdC1zYWx0" // base64 encoded "test-salt"
	password := "test_password"

	crypto1.SetSalt(salt)
	crypto2.SetSalt(salt)

	// Генерируем ключи с одинаковыми параметрами
	err1 := crypto1.GenerateAndStoreEncryptionKey(password)
	err2 := crypto2.GenerateAndStoreEncryptionKey(password)

	assert.NoError(t, err1)
	assert.NoError(t, err2)

	// Ключи должны быть одинаковыми (детерминированными)
	assert.Equal(t, crypto1.encryptionKey, crypto2.encryptionKey, "Ключи должны быть детерминированными")

	// Шифруем одинаковые данные - должны получить разные результаты из-за случайного nonce
	testData := "test data"
	encrypted1, err := crypto1.EncryptString(testData)
	assert.NoError(t, err)

	encrypted2, err := crypto2.EncryptString(testData)
	assert.NoError(t, err)

	// Зашифрованные данные должны быть разными из-за случайного nonce в AES-GCM
	assert.NotEqual(t, encrypted1, encrypted2, "Зашифрованные данные должны быть разными из-за случайного nonce")

	// Но дешифрование должно работать для обоих вариантов
	decrypted1, err := crypto1.DecryptString(encrypted1)
	assert.NoError(t, err)
	assert.Equal(t, testData, decrypted1)

	decrypted2, err := crypto2.DecryptString(encrypted2)
	assert.NoError(t, err)
	assert.Equal(t, testData, decrypted2)

	// Также можно дешифровать данные, зашифрованные другим экземпляром
	decrypted1WithCrypto2, err := crypto2.DecryptString(encrypted1)
	assert.NoError(t, err)
	assert.Equal(t, testData, decrypted1WithCrypto2)

	decrypted2WithCrypto1, err := crypto1.DecryptString(encrypted2)
	assert.NoError(t, err)
	assert.Equal(t, testData, decrypted2WithCrypto1)
}

func TestEncryptionKeyConsistency(t *testing.T) {
	// Создаем два экземпляра с одинаковыми параметрами
	crypto1 := NewCrypto()
	crypto2 := NewCrypto()

	salt := "dGVzdC1zYWx0" // base64 encoded "test-salt"
	password := "my_secret_password"

	crypto1.SetSalt(salt)
	crypto2.SetSalt(salt)

	// Генерируем ключи - они должны быть одинаковыми
	err1 := crypto1.GenerateAndStoreEncryptionKey(password)
	err2 := crypto2.GenerateAndStoreEncryptionKey(password)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, crypto1.encryptionKey, crypto2.encryptionKey, "Ключи должны быть одинаковыми")

	// Шифруем данные несколько раз
	testData := "sensitive_data:password123"

	// Первое шифрование
	encrypted1, err := crypto1.EncryptString(testData)
	assert.NoError(t, err)

	// Второе шифрование (должно дать другой результат из-за случайного nonce)
	encrypted2, err := crypto1.EncryptString(testData)
	assert.NoError(t, err)

	// Результаты шифрования разные, но оба корректно расшифровываются
	assert.NotEqual(t, encrypted1, encrypted2, "Зашифрованные данные должны быть разными")

	// Дешифруем первый результат
	decrypted1, err := crypto1.DecryptString(encrypted1)
	assert.NoError(t, err)
	assert.Equal(t, testData, decrypted1)

	// Дешифруем второй результат
	decrypted2, err := crypto1.DecryptString(encrypted2)
	assert.NoError(t, err)
	assert.Equal(t, testData, decrypted2)

	// Дешифруем первый результат с помощью второго экземпляра (тот же ключ)
	decrypted1WithCrypto2, err := crypto2.DecryptString(encrypted1)
	assert.NoError(t, err)
	assert.Equal(t, testData, decrypted1WithCrypto2)

	// Дешифруем второй результат с помощью второго экземпляра
	decrypted2WithCrypto2, err := crypto2.DecryptString(encrypted2)
	assert.NoError(t, err)
	assert.Equal(t, testData, decrypted2WithCrypto2)

	t.Logf("Ключ шифрования: %x", crypto1.encryptionKey)
	t.Logf("Длина зашифрованных данных 1: %d байт", len(encrypted1))
	t.Logf("Длина зашифрованных данных 2: %d байт", len(encrypted2))
}

func TestGenerateAndStoreEncryptionKey_Base64Issue(t *testing.T) {
	crypto := NewCrypto()

	// Используем реальную соль, которая могла бы вызвать проблему
	salt := "dGVzdC1zYWx0" // base64 encoded "test-salt"
	crypto.SetSalt(salt)

	password := "test_password"

	// Это должно работать без ошибок base64
	err := crypto.GenerateAndStoreEncryptionKey(password)
	assert.NoError(t, err)
	assert.NotNil(t, crypto.encryptionKey)
	assert.Len(t, crypto.encryptionKey, 32, "Ключ должен быть 32 байта для AES-256")

	// Проверяем, что шифрование и дешифрование работают
	testData := "test data"
	encrypted, err := crypto.EncryptString(testData)
	assert.NoError(t, err)

	decrypted, err := crypto.DecryptString(encrypted)
	assert.NoError(t, err)
	assert.Equal(t, testData, decrypted)
}
