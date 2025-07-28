package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"gophkeeper/internal/models"
	"gophkeeper/internal/server/usecase"

	"go.uber.org/zap"
)

// AuthenticateHandler обрабатывает запросы аутентификации пользователей
type AuthenticateHandler struct {
	uc usecase.AuthenticateUsecase
}

// NewAuthenticateHandler создает новый экземпляр AuthenticateHandler
func NewAuthenticateHandler(uc usecase.AuthenticateUsecase) *AuthenticateHandler {
	return &AuthenticateHandler{
		uc: uc,
	}
}

// Authenticate godoc
// @Summary Аутентификация пользователя
// @Description Проверяет хеш пароля и возвращает токены доступа
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.AuthenticateIn true "Данные для аутентификации"
// @Success 200 {object} models.AuthenticateOut "Успешная аутентификация, возвращены токены"
// @Failure 400 {object} models.JSONError "Некорректные данные запроса"
// @Failure 401 {object} models.JSONError "Неверные учетные данные"
// @Failure 500 {object} models.JSONError "Внутренняя ошибка сервера"
// @Router /api/authenticate [post]
func (handler AuthenticateHandler) Authenticate(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		models.EncodeError(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	var in models.AuthenticateIn
	var buf bytes.Buffer
	// читаем тело запроса
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		models.EncodeError(res, err.Error(), http.StatusBadRequest)
		return
	}

	// десериализуем JSON
	if err = json.Unmarshal(buf.Bytes(), &in); err != nil {
		models.EncodeError(res, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := handler.uc.Authenticate(req.Context(), in.Login, in.PasswordHash)
	if err != nil {
		if errors.Is(err, models.ErrLoginOrPasswordEmpty) {
			zap.L().Sugar().Debugln("Login or password is empty", zap.Error(err))
			models.EncodeError(res, err.Error(), http.StatusBadRequest)
			return
		}

		if errors.Is(err, models.ErrUserNotFound) {
			zap.L().Sugar().Debugln("User not found", zap.Error(err))
			models.EncodeError(res, err.Error(), http.StatusUnauthorized)
			return
		}

		zap.L().Sugar().Debugln("Failed to authenticate user", zap.Error(err))
		models.EncodeError(res, "Failed to authenticate user", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(res).Encode(result); err != nil {
		zap.L().Sugar().Debugln("Failed to encode response", zap.Error(err))
		models.EncodeError(res, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
