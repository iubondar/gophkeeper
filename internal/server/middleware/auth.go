// Package middleware предоставляет HTTP middleware для сервера.
// Включает middleware для аутентификации и других общих функций.
package middleware

import (
	"context"
	"errors"
	"net/http"

	"gophkeeper/internal/auth"
	"gophkeeper/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UserIDKey - ключ для хранения userID в контексте
type UserIDKey struct{}

// GetUserIDFromContext извлекает userID из контекста запроса
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey{}).(uuid.UUID)
	return userID, ok
}

// AuthMiddleware проверяет наличие и валидность JWT токена в запросе.
// Если токен отсутствует или недействителен, возвращает 401 Unauthorized.
// Если токен валиден, добавляет userID в контекст запроса.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserIDFromReq(r)
		if err != nil {
			// Проверяем тип ошибки для более точного сообщения
			if errors.Is(err, jwt.ErrTokenExpired) {
				zap.L().Sugar().Debugln("JWT token expired", zap.Error(err))
				models.EncodeError(w, "Authentication token expired", http.StatusUnauthorized)
				return
			}
			if errors.Is(err, jwt.ErrTokenNotValidYet) {
				zap.L().Sugar().Debugln("JWT token not valid yet", zap.Error(err))
				models.EncodeError(w, "Authentication token not valid yet", http.StatusUnauthorized)
				return
			}
			if errors.Is(err, jwt.ErrTokenMalformed) {
				zap.L().Sugar().Debugln("JWT token malformed", zap.Error(err))
				models.EncodeError(w, "Malformed authentication token", http.StatusUnauthorized)
				return
			}
			if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				zap.L().Sugar().Debugln("JWT token signature invalid", zap.Error(err))
				models.EncodeError(w, "Invalid authentication token signature", http.StatusUnauthorized)
				return
			}

			zap.L().Sugar().Debugln("Authentication failed", zap.Error(err))
			models.EncodeError(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		if userID == uuid.Nil {
			zap.L().Sugar().Debugln("Invalid or missing authentication token")
			models.EncodeError(w, "Invalid or missing authentication token", http.StatusUnauthorized)
			return
		}

		// Добавляем userID в контекст запроса для использования в хэндлерах
		ctx := context.WithValue(r.Context(), UserIDKey{}, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
