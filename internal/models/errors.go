package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrLoginOrPasswordEmpty = errors.New("пустой логин или пароль")
	ErrUserNotFound         = errors.New("пользователь не найден")
	ErrUserAlreadyExists    = errors.New("пользователь уже существует")
	ErrAccessTokenExpired   = errors.New("access token expired")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
	ErrRefreshTokenInvalid  = errors.New("invalid refresh token")
	ErrUnauthorized         = errors.New("unauthorized")
)

// JSONError представляет ошибку в формате JSON для API ответов
type JSONError struct {
	Message string `json:"message" example:"Error message"` // Сообщение об ошибке
	Code    int    `json:"code" example:"400"`              // HTTP код ошибки
}

func (e *JSONError) Error() string {
	return e.Message
}

// NewJSONError создает новый экземпляр JSONError
func NewJSONError(message string, code int) *JSONError {
	return &JSONError{Message: message, Code: code}
}

// Encode кодирует JSONError в HTTP ответ
func (e *JSONError) Encode(res http.ResponseWriter) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(e.Code)
	if err := json.NewEncoder(res).Encode(e); err != nil {
		// Используем JSONError вместо http.Error
		jsonErr := NewJSONError("Failed to encode response", http.StatusInternalServerError)
		jsonErr.Encode(res)
		return
	}
}

// EncodeError кодирует ошибку в HTTP ответ
func EncodeError(res http.ResponseWriter, message string, code int) {
	jsonErr := NewJSONError(message, code)
	jsonErr.Encode(res)
}

// ParseJSONError разбирает JSONError из тела ответа ([]byte)
func ParseJSONError(body []byte) (*JSONError, error) {
	var jsonErr JSONError
	if err := json.Unmarshal(body, &jsonErr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON error: %w", err)
	}
	return &jsonErr, nil
}
