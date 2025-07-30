package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"gophkeeper/internal/models"
	"gophkeeper/internal/server/middleware"
	"gophkeeper/internal/server/usecase"

	"go.uber.org/zap"
)

// UpdateHandler обрабатывает HTTP-запросы для обновления секретов.
// Обработчик реализует endpoint /api/update и требует аутентификации пользователя.
// Использует usecase-слой для обновления секретов с проверкой версий.
type UpdateHandler struct {
	uc usecase.UpdateSecretUsecase
}

// NewUpdateHandler создает новый экземпляр UpdateHandler.
// Принимает usecase для обновления секретов.
// Функция используется для внедрения зависимостей и создания обработчика
// с конкретной реализацией бизнес-логики обновления секретов.
func NewUpdateHandler(uc usecase.UpdateSecretUsecase) *UpdateHandler {
	return &UpdateHandler{uc: uc}
}

// Update godoc
// @Summary Обновление секрета
// @Description Обновляет зашифрованный секрет на сервере с проверкой версии
// @Tags secrets
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.UpdateSecretIn true "Данные секрета для обновления"
// @Success 200 {object} models.UpdateSecretOut "Секрет успешно обновлен"
// @Failure 400 {object} models.JSONError "Некорректные данные запроса"
// @Failure 401 {object} models.JSONError "Неавторизованный доступ"
// @Failure 409 {object} models.JSONError "Конфликт версий"
// @Failure 500 {object} models.JSONError "Внутренняя ошибка сервера"
// @Router /api/update [put]
func (handler UpdateHandler) Update(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		models.EncodeError(res, "Only PUT requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	// Получаем userID из контекста (установлен middleware)
	userID, ok := middleware.GetUserIDFromContext(req.Context())
	if !ok {
		zap.L().Sugar().Debugln("User ID not found in context")
		models.EncodeError(res, "Authentication required", http.StatusUnauthorized)
		return
	}

	var in models.UpdateSecretIn
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		models.EncodeError(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &in); err != nil {
		models.EncodeError(res, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := handler.uc.UpdateSecret(req.Context(), in, userID)
	if err != nil {
		if err == models.ErrConflict {
			models.EncodeError(res, "Version conflict", http.StatusConflict)
			return
		}
		zap.L().Sugar().Debugln("Failed to update secret", zap.Error(err))
		models.EncodeError(res, "Failed to update secret", http.StatusInternalServerError)
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
