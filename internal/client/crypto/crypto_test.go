package crypto

import (
	"bytes"
	"crypto/aes"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Вспомогательные структуры для тестирования ошибок

type errorWriter struct {
	errorOnWrite bool
	writeCount   int
}

func (ew *errorWriter) Write(p []byte) (n int, err error) {
	ew.writeCount++
	if ew.writeCount == 1 {
		// Первая запись - это IV, пропускаем
		return len(p), nil
	}
	if ew.errorOnWrite {
		return 0, fmt.Errorf("write error")
	}
	return len(p), nil
}

type errorReader struct{}

func (er *errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("read IV error")
}

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

func TestStreamEncryptionDecryption(t *testing.T) {
	crypto := NewCrypto()
	// Генерируем корректную соль в base64
	crypto.GenerateAndSetSalt()
	err := crypto.GenerateAndStoreEncryptionKey("test-password")
	require.NoError(t, err)

	// Тестовые данные
	originalData := []byte("Это тестовые данные для потокового шифрования. " +
		"Они должны быть корректно зашифрованы и расшифрованы. " +
		"Потоковое шифрование позволяет обрабатывать большие файлы без загрузки их в память целиком.")

	// Создаем буфер для зашифрованных данных
	var encryptedBuffer bytes.Buffer

	// Создаем потоковый шифратор
	encryptWriter, err := crypto.EncryptStream(&encryptedBuffer)
	require.NoError(t, err)

	// Записываем данные через шифратор
	_, err = encryptWriter.Write(originalData)
	require.NoError(t, err)

	// Закрываем шифратор
	err = encryptWriter.Close()
	require.NoError(t, err)

	// Проверяем, что данные зашифрованы (не равны оригиналу)
	encryptedData := encryptedBuffer.Bytes()
	require.NotEqual(t, originalData, encryptedData)
	require.Greater(t, len(encryptedData), len(originalData)) // Должен быть IV в начале

	// Создаем потоковый дешифратор
	decryptReader, err := crypto.DecryptStream(bytes.NewReader(encryptedData))
	require.NoError(t, err)
	defer decryptReader.Close()

	// Читаем и расшифровываем данные
	decryptedData := make([]byte, len(originalData))
	n, err := decryptReader.Read(decryptedData)
	require.NoError(t, err)
	require.Equal(t, len(originalData), n)

	// Проверяем, что расшифрованные данные совпадают с оригиналом
	require.Equal(t, originalData, decryptedData)
}

func TestStreamEncryptionLargeData(t *testing.T) {
	crypto := NewCrypto()
	// Генерируем корректную соль в base64
	crypto.GenerateAndSetSalt()
	err := crypto.GenerateAndStoreEncryptionKey("test-password")
	require.NoError(t, err)

	// Создаем большие тестовые данные (1MB)
	originalData := make([]byte, 1024*1024)
	for i := range originalData {
		originalData[i] = byte(i % 256)
	}

	// Создаем буфер для зашифрованных данных
	var encryptedBuffer bytes.Buffer

	// Создаем потоковый шифратор
	encryptWriter, err := crypto.EncryptStream(&encryptedBuffer)
	require.NoError(t, err)

	// Записываем данные блоками
	blockSize := 4096
	for i := 0; i < len(originalData); i += blockSize {
		end := i + blockSize
		if end > len(originalData) {
			end = len(originalData)
		}
		_, err = encryptWriter.Write(originalData[i:end])
		require.NoError(t, err)
	}

	// Закрываем шифратор
	err = encryptWriter.Close()
	require.NoError(t, err)

	// Создаем потоковый дешифратор
	decryptReader, err := crypto.DecryptStream(bytes.NewReader(encryptedBuffer.Bytes()))
	require.NoError(t, err)
	defer decryptReader.Close()

	// Читаем и расшифровываем данные блоками
	decryptedData := make([]byte, len(originalData))
	totalRead := 0
	for totalRead < len(originalData) {
		n, err := decryptReader.Read(decryptedData[totalRead:])
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		totalRead += n
	}

	// Проверяем, что все данные прочитаны и расшифрованы корректно
	require.Equal(t, len(originalData), totalRead)
	require.Equal(t, originalData, decryptedData)
}

func TestStreamEncryptionEmptyData(t *testing.T) {
	crypto := NewCrypto()
	// Генерируем корректную соль в base64
	crypto.GenerateAndSetSalt()
	err := crypto.GenerateAndStoreEncryptionKey("test-password")
	require.NoError(t, err)

	// Тестируем с пустыми данными
	originalData := []byte{}

	var encryptedBuffer bytes.Buffer
	encryptWriter, err := crypto.EncryptStream(&encryptedBuffer)
	require.NoError(t, err)

	_, err = encryptWriter.Write(originalData)
	require.NoError(t, err)

	err = encryptWriter.Close()
	require.NoError(t, err)

	// Проверяем, что создан только IV
	encryptedData := encryptedBuffer.Bytes()
	require.Equal(t, aes.BlockSize, len(encryptedData))

	// Расшифровываем
	decryptReader, err := crypto.DecryptStream(bytes.NewReader(encryptedData))
	require.NoError(t, err)
	defer decryptReader.Close()

	decryptedData := make([]byte, 1)
	n, err := decryptReader.Read(decryptedData)
	require.Equal(t, io.EOF, err)
	require.Equal(t, 0, n)
}

// Тесты для покрытия ошибок

func TestCalcArgon2Hash_InvalidSalt(t *testing.T) {
	crypto := NewCrypto()

	// Тестируем с некорректной солью (не base64)
	_, err := crypto.calcArgon2Hash("password", "invalid-salt-not-base64")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "illegal base64")
}

func TestCalcArgon2Key_InvalidSalt(t *testing.T) {
	crypto := NewCrypto()

	// Тестируем с некорректной солью (не base64)
	_, err := crypto.calcArgon2Key("password", "invalid-salt-not-base64")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "illegal base64")
}

func TestVerifyPasswordHash_InvalidSalt(t *testing.T) {
	crypto := NewCrypto()

	// Тестируем с некорректной солью
	_, err := crypto.VerifyPasswordHash("password", "invalid-salt", "hash")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "illegal base64")
}

func TestGenerateAndStoreEncryptionKey_NoSalt(t *testing.T) {
	crypto := NewCrypto()
	// Не устанавливаем соль

	err := crypto.GenerateAndStoreEncryptionKey("password")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "salt is not set")
}

func TestGenerateAndStoreEncryptionKey_InvalidSalt(t *testing.T) {
	crypto := NewCrypto()
	crypto.SetSalt("invalid-salt-not-base64")

	err := crypto.GenerateAndStoreEncryptionKey("password")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid salt format")
}

func TestEncryptString_NoEncryptionKey(t *testing.T) {
	crypto := NewCrypto()
	// Не устанавливаем ключ шифрования

	_, err := crypto.EncryptString("test data")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "encryption key is not set")
}

func TestDecryptString_NoEncryptionKey(t *testing.T) {
	crypto := NewCrypto()
	// Не устанавливаем ключ шифрования

	_, err := crypto.DecryptString([]byte("test data"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "encryption key is not set")
}

func TestEncryptStream_NoEncryptionKey(t *testing.T) {
	crypto := NewCrypto()
	// Не устанавливаем ключ шифрования

	var buf bytes.Buffer
	_, err := crypto.EncryptStream(&buf)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "encryption key is not set")
}

func TestDecryptStream_NoEncryptionKey(t *testing.T) {
	crypto := NewCrypto()
	// Не устанавливаем ключ шифрования

	_, err := crypto.DecryptStream(bytes.NewReader([]byte("test data")))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "encryption key is not set")
}

func TestDecryptString_InvalidData(t *testing.T) {
	crypto := NewCrypto()
	crypto.SetSalt("dGVzdC1zYWx0")
	err := crypto.GenerateAndStoreEncryptionKey("test-password")
	require.NoError(t, err)

	// Тестируем с данными слишком короткими для nonce
	_, err = crypto.DecryptString([]byte("short"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid encrypted data")
}

func TestDecryptString_CorruptedData(t *testing.T) {
	crypto := NewCrypto()
	crypto.SetSalt("dGVzdC1zYWx0")
	err := crypto.GenerateAndStoreEncryptionKey("test-password")
	require.NoError(t, err)

	// Создаем некорректные данные с правильной длиной nonce, но испорченным ciphertext
	nonce := make([]byte, 12) // GCM nonce size
	rand.Read(nonce)
	corruptedData := append(nonce, []byte("corrupted-ciphertext")...)

	_, err = crypto.DecryptString(corruptedData)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "decrypt")
}

func TestDecryptStream_ReadError(t *testing.T) {
	crypto := NewCrypto()
	crypto.GenerateAndSetSalt()
	err := crypto.GenerateAndStoreEncryptionKey("test-password")
	require.NoError(t, err)

	// Создаем reader, который всегда возвращает ошибку
	errorReader := &errorReader{}

	_, err = crypto.DecryptStream(errorReader)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "read IV")
}

func TestDecryptReader_Close(t *testing.T) {
	dr := &decryptReader{
		stream: nil,
		reader: &bytes.Buffer{},
	}

	// Close должен работать без ошибок
	err := dr.Close()
	assert.NoError(t, err)
}

func TestStreamEncryptionWithErrors(t *testing.T) {
	crypto := NewCrypto()
	crypto.GenerateAndSetSalt()
	err := crypto.GenerateAndStoreEncryptionKey("test-password")
	require.NoError(t, err)

	// Тестируем с данными, которые вызывают ошибки при записи
	errorWriter := &errorWriter{errorOnWrite: true}

	// Сначала создаем поток (IV записывается здесь)
	encryptWriter, err := crypto.EncryptStream(errorWriter)
	require.NoError(t, err)

	// Попытка записи должна вернуть ошибку
	_, err = encryptWriter.Write([]byte("test data"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "write error")

	// Close должен работать
	err = encryptWriter.Close()
	assert.NoError(t, err)
}

func TestStreamDecryptionWithErrors(t *testing.T) {
	crypto := NewCrypto()
	crypto.GenerateAndSetSalt()
	err := crypto.GenerateAndStoreEncryptionKey("test-password")
	require.NoError(t, err)

	// Создаем reader, который возвращает ошибку при чтении
	errorReader := &errorReader{}

	_, err = crypto.DecryptStream(errorReader)
	assert.Error(t, err)
}

func TestEncryptString_WithLargeData(t *testing.T) {
	crypto := NewCrypto()
	crypto.SetSalt("dGVzdC1zYWx0")
	err := crypto.GenerateAndStoreEncryptionKey("test-password")
	require.NoError(t, err)

	// Создаем большие данные
	largeData := strings.Repeat("test data ", 1000)

	encrypted, err := crypto.EncryptString(largeData)
	assert.NoError(t, err)
	assert.NotEqual(t, largeData, string(encrypted))

	decrypted, err := crypto.DecryptString(encrypted)
	assert.NoError(t, err)
	assert.Equal(t, largeData, decrypted)
}

func TestEncryptString_WithSpecialCharacters(t *testing.T) {
	crypto := NewCrypto()
	crypto.SetSalt("dGVzdC1zYWx0")
	err := crypto.GenerateAndStoreEncryptionKey("test-password")
	require.NoError(t, err)

	// Тестируем с специальными символами
	specialData := "Привет мир! 🌍 你好世界! Hello World! 1234567890 !@#$%^&*()_+-=[]{}|;':\",./<>?"

	encrypted, err := crypto.EncryptString(specialData)
	assert.NoError(t, err)
	assert.NotEqual(t, specialData, string(encrypted))

	decrypted, err := crypto.DecryptString(encrypted)
	assert.NoError(t, err)
	assert.Equal(t, specialData, decrypted)
}
