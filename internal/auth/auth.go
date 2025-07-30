// Package auth предоставляет функциональность для работы с JWT токенами аутентификации.
// Включает генерацию access и refresh токенов, а также извлечение информации о пользователе
// из HTTP запросов.
package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const secretKey = "supersecretkey"

const accessTokenDuration = 15 * time.Minute
const refreshTokenDuration = 7 * 24 * time.Hour

// AuthCookieName - имя cookie для хранения токена аутентификации
const AuthCookieName = "Authorization"

// GenerateAccessToken создает JWT access токен для указанного пользователя.
// Токен действителен в течение 15 минут.
//
// Параметры:
//   - userID: уникальный идентификатор пользователя
//
// Возвращает:
//   - string: подписанный JWT токен
//   - error: ошибка в случае неудачи
func GenerateAccessToken(userID string) (string, error) {
	return generateToken(userID, accessTokenDuration)
}

// GenerateRefreshToken создает JWT refresh токен для указанного пользователя.
// Токен действителен в течение 7 дней.
//
// Параметры:
//   - userID: уникальный идентификатор пользователя
//
// Возвращает:
//   - string: подписанный JWT токен
//   - error: ошибка в случае неудачи
func GenerateRefreshToken(userID string) (string, error) {
	return generateToken(userID, refreshTokenDuration)
}

func generateToken(userID string, duration time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// GetUserIDFromReq извлекает идентификатор пользователя из HTTP запроса.
// Ищет токен в cookie с именем AuthCookieName и валидирует его.
//
// Параметры:
//   - req: HTTP запрос для извлечения токена
//
// Возвращает:
//   - uuid.UUID: идентификатор пользователя или uuid.Nil если токен отсутствует/недействителен
//   - error: ошибка в случае неудачи
func GetUserIDFromReq(req *http.Request) (userID uuid.UUID, err error) {
	authCookie, err := req.Cookie(AuthCookieName)
	if err != nil {
		zap.L().Sugar().Debugln("No auth cookie found")
		return uuid.Nil, err
	}

	userID, err = getUserID(authCookie.Value)
	if err != nil {
		zap.L().Sugar().Debugln("Error getting user id from cookie. Message: ", err.Error())
		return uuid.Nil, err
	}

	return userID, nil
}

func getUserID(tokenString string) (userID uuid.UUID, err error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		return uuid.Nil, err
	}

	// Дополнительная проверка валидности токена
	if !token.Valid {
		return uuid.Nil, fmt.Errorf("token is not valid")
	}

	// Parse UUID from Subject claim
	userID, err = uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID in token: %v", err)
	}

	return userID, nil
}

// ValidateRefreshToken валидирует refresh token и возвращает userID.
// Функция проверяет подпись токена и его срок действия.
//
// Параметры:
//   - tokenString: строка refresh токена для валидации
//
// Возвращает:
//   - uuid.UUID: идентификатор пользователя или uuid.Nil если токен недействителен
//   - error: ошибка в случае неудачи
func ValidateRefreshToken(tokenString string) (userID uuid.UUID, err error) {
	return getUserID(tokenString)
}
