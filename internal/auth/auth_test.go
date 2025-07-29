package auth

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func init() {
	// Initialize logger for tests
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}

func TestGenerateAccessToken(t *testing.T) {
	// Create a test user ID
	userID := uuid.New().String()

	// Call the function
	tokenString, err := GenerateAccessToken(userID)
	require.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Parse the token to verify its contents
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	require.NoError(t, err)
	assert.True(t, token.Valid)

	// Verify claims
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	require.True(t, ok)
	assert.Equal(t, userID, claims.Subject)
}

func TestGenerateRefreshToken(t *testing.T) {
	// Create a test user ID
	userID := uuid.New().String()

	// Call the function
	tokenString, err := GenerateRefreshToken(userID)
	require.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Parse the token to verify its contents
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	require.NoError(t, err)
	assert.True(t, token.Valid)

	// Verify claims
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	require.True(t, ok)
	assert.Equal(t, userID, claims.Subject)
}

func TestGetUserIDFromReq(t *testing.T) {
	// Create a test user ID
	userID := uuid.New()

	// Create a valid JWT token
	tokenString, err := GenerateAccessToken(userID.String())
	require.NoError(t, err)

	// Create a request with the cookie
	req := httptest.NewRequest("GET", "/", nil)
	cookie := &http.Cookie{
		Name:  AuthCookieName,
		Value: tokenString,
	}
	req.AddCookie(cookie)

	// Test valid cookie
	extractedUserID, err := GetUserIDFromReq(req)
	require.NoError(t, err)
	assert.Equal(t, userID, extractedUserID)

	// Test missing cookie
	req = httptest.NewRequest("GET", "/", nil)
	_, err = GetUserIDFromReq(req)
	assert.Error(t, err)
}

func TestGetUserID(t *testing.T) {
	// Create a test user ID
	userID := uuid.New()

	// Test valid token
	tokenString, err := GenerateAccessToken(userID.String())
	require.NoError(t, err)

	extractedUserID, err := getUserID(tokenString)
	require.NoError(t, err)
	assert.Equal(t, userID, extractedUserID)

	// Test invalid token
	_, err = getUserID("invalid.token.string")
	assert.Error(t, err)

	// Test expired token
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
	})
	expiredTokenString, err := expiredToken.SignedString([]byte(secretKey))
	require.NoError(t, err)
	_, err = getUserID(expiredTokenString)
	assert.Error(t, err)
}

func TestGetUserID_ExpiredToken(t *testing.T) {
	// Создаем токен с истекшим временем
	userID := uuid.New().String()

	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Токен истек час назад
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	require.NoError(t, err)

	// Пытаемся получить userID из истекшего токена
	result, err := getUserID(tokenString)

	// Должны получить ошибку и uuid.Nil
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, result)
}

func TestGetUserID_ValidToken(t *testing.T) {
	// Создаем валидный токен
	userID := uuid.New().String()

	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // Токен действителен час
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	require.NoError(t, err)

	// Пытаемся получить userID из валидного токена
	result, err := getUserID(tokenString)

	// Должны получить userID без ошибки
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, result)

	// Проверяем, что полученный userID соответствует ожидаемому
	expectedUserID, err := uuid.Parse(userID)
	require.NoError(t, err)
	assert.Equal(t, expectedUserID, result)
}

func TestJWTExpirationBehavior(t *testing.T) {
	// Создаем токен без ExpiresAt
	userID := uuid.New().String()

	claims := jwt.RegisteredClaims{
		Subject:  userID,
		IssuedAt: jwt.NewNumericDate(time.Now()),
		// Нет ExpiresAt
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	require.NoError(t, err)

	// Пытаемся получить userID из токена без ExpiresAt
	result, err := getUserID(tokenString)

	// Должны получить userID без ошибки, так как exp опциональный
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, result)

	// Создаем токен с истекшим ExpiresAt
	claims = jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Истек час назад
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(secretKey))
	require.NoError(t, err)

	// Пытаемся получить userID из истекшего токена
	result, err = getUserID(tokenString)

	// Должны получить ошибку
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, result)

	// Проверяем, что это именно ошибка expiration
	assert.Contains(t, err.Error(), "expired")
}

func TestJWTDetailedBehavior(t *testing.T) {
	// Тест 1: Токен без ExpiresAt (должен быть валидным)
	userID := uuid.New().String()

	claims := jwt.RegisteredClaims{
		Subject:  userID,
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	require.NoError(t, err)

	// Парсим токен напрямую через JWT библиотеку
	parsedToken, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{},
		func(t *jwt.Token) (any, error) {
			return []byte(secretKey), nil
		})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	// Тест 2: Токен с истекшим ExpiresAt
	claims = jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(secretKey))
	require.NoError(t, err)

	// Парсим истекший токен - JWT библиотека возвращает ошибку
	parsedToken, err = jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{},
		func(t *jwt.Token) (any, error) {
			return []byte(secretKey), nil
		})

	assert.Error(t, err) // Парсинг неуспешен из-за expiration
	assert.Contains(t, err.Error(), "expired")

	// Проверяем, что наша функция getUserID правильно обрабатывает это
	result, err := getUserID(tokenString)
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, result)

	// Проверяем, что ошибка содержит информацию об expiration
	assert.Contains(t, err.Error(), "expired")
}

func TestTokenValidCheck(t *testing.T) {
	// Создаем токен с неправильным алгоритмом подписи
	userID := uuid.New().String()

	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	// Создаем токен с RSA алгоритмом, но подписываем HMAC ключом
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey)) // Это вызовет ошибку
	require.Error(t, err)                                     // Ожидаем ошибку, так как RSA токен нельзя подписать HMAC ключом

	// Но если бы мы получили такой токен каким-то образом,
	// то без проверки token.Valid мы могли бы его принять
	// Давайте создадим валидный токен и изменим его алгоритм
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(secretKey))
	require.NoError(t, err)

	// Теперь изменим алгоритм в заголовке токена (это демонстрирует проблему)
	// В реальности это было бы сложнее, но это показывает концепцию
	parts := strings.Split(tokenString, ".")
	require.Equal(t, 3, len(parts))

	// Декодируем заголовок
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	require.NoError(t, err)

	var header map[string]interface{}
	err = json.Unmarshal(headerBytes, &header)
	require.NoError(t, err)

	// Изменяем алгоритм
	header["alg"] = "RS256"

	// Кодируем обратно
	headerBytes, err = json.Marshal(header)
	require.NoError(t, err)
	parts[0] = base64.RawURLEncoding.EncodeToString(headerBytes)

	// Собираем токен обратно
	modifiedTokenString := strings.Join(parts, ".")

	// Теперь у нас есть токен с неправильным алгоритмом
	// Без проверки token.Valid мы могли бы его принять
	result, err := getUserID(modifiedTokenString)

	// Должны получить ошибку из-за неправильного алгоритма
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, result)
}

func TestTokenValidCheckWithValidToken(t *testing.T) {
	// Создаем валидный токен
	userID := uuid.New().String()

	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	require.NoError(t, err)

	// Парсим токен напрямую
	parsedToken, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{},
		func(t *jwt.Token) (any, error) {
			return []byte(secretKey), nil
		})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid) // Токен валиден

	// Проверяем через нашу функцию
	result, err := getUserID(tokenString)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, result)

	// Теперь создадим токен, который парсится, но не валиден
	// Это сложно сделать с текущей JWT библиотекой, но мы можем
	// продемонстрировать концепцию с помощью токена с неправильной подписью

	// Создаем токен с неправильным ключом
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	wrongTokenString, err := token.SignedString([]byte("wrongkey"))
	require.NoError(t, err)

	// Пытаемся парсить с правильным ключом
	parsedToken, err = jwt.ParseWithClaims(wrongTokenString, &jwt.RegisteredClaims{},
		func(t *jwt.Token) (any, error) {
			return []byte(secretKey), nil
		})

	// Должны получить ошибку подписи
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "signature")
}
