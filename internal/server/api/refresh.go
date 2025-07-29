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

// RefreshHandler обрабатывает HTTP-запросы для обновления токенов.
// Обработчик реализует endpoint /api/refresh и валидирует refresh токен,
// возвращая новые JWT токены доступа при успешной валидации.
type RefreshHandler struct {
	uc usecase.RefreshUsecase
}

// NewRefreshHandler создает новый экземпляр RefreshHandler.
// Принимает usecase для обновления токенов.
// Функция используется для внедрения зависимостей и создания обработчика
// с конкретной реализацией бизнес-логики обновления токенов.
func NewRefreshHandler(uc usecase.RefreshUsecase) *RefreshHandler {
	return &RefreshHandler{
		uc: uc,
	}
}

// Refresh godoc
// @Summary Обновление токенов
// @Description Валидирует refresh токен и возвращает новые токены доступа
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RefreshIn true "Данные для обновления токенов"
// @Success 200 {object} models.RefreshOut "Успешное обновление токенов"
// @Failure 400 {object} models.JSONError "Некорректные данные запроса"
// @Failure 401 {object} models.JSONError "Недействительный refresh токен"
// @Failure 500 {object} models.JSONError "Внутренняя ошибка сервера"
// @Router /api/refresh [post]
func (handler RefreshHandler) Refresh(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		models.EncodeError(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	var in models.RefreshIn
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

	result, err := handler.uc.Refresh(req.Context(), in.RefreshToken)
	if err != nil {
		if errors.Is(err, models.ErrRefreshTokenInvalid) {
			zap.L().Sugar().Debugln("Invalid refresh token", zap.Error(err))
			models.EncodeError(res, err.Error(), http.StatusBadRequest)
			return
		}

		if errors.Is(err, models.ErrRefreshTokenExpired) {
			zap.L().Sugar().Debugln("Refresh token expired", zap.Error(err))
			models.EncodeError(res, err.Error(), http.StatusUnauthorized)
			return
		}

		zap.L().Sugar().Debugln("Failed to refresh tokens", zap.Error(err))
		models.EncodeError(res, "Failed to refresh tokens", http.StatusInternalServerError)
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
