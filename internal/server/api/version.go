package api

import (
	"encoding/json"
	"net/http"

	"gophkeeper/internal/models"
	"gophkeeper/internal/server/middleware"
	"gophkeeper/internal/server/usecase"

	"go.uber.org/zap"
)

// VersionHandler обрабатывает HTTP-запросы для получения версий секретов.
// Обработчик реализует endpoint /api/version и требует аутентификации пользователя.
// Использует usecase-слой для получения информации о версии секрета.
type VersionHandler struct {
	uc usecase.GetSecretUsecase
}

// NewVersionHandler создает новый экземпляр VersionHandler.
// Принимает usecase для получения секретов.
// Функция используется для внедрения зависимостей и создания обработчика
// с конкретной реализацией бизнес-логики получения версий секретов.
func NewVersionHandler(uc usecase.GetSecretUsecase) *VersionHandler {
	return &VersionHandler{uc: uc}
}

// GetSecretVersionHandler godoc
// @Summary Получение версии секрета
// @Description Возвращает только версию секрета по имени
// @Tags secrets
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query string true "Имя секрета"
// @Success 200 {object} map[string]int "Версия секрета"
// @Failure 400 {object} models.JSONError "Некорректные данные запроса"
// @Failure 401 {object} models.JSONError "Неавторизованный доступ"
// @Failure 404 {object} models.JSONError "Секрет не найден"
// @Failure 500 {object} models.JSONError "Внутренняя ошибка сервера"
// @Router /api/version [get]
func (handler VersionHandler) GetSecretVersion(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		models.EncodeError(res, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	// Получаем userID из контекста (установлен middleware)
	userID, ok := middleware.GetUserIDFromContext(req.Context())
	if !ok {
		zap.L().Sugar().Debugln("User ID not found in context")
		models.EncodeError(res, "Authentication required", http.StatusUnauthorized)
		return
	}

	secretName := req.URL.Query().Get("name")
	if secretName == "" {
		models.EncodeError(res, "Secret name is required", http.StatusBadRequest)
		return
	}

	result, err := handler.uc.GetSecret(req.Context(), secretName, userID)
	if err != nil {
		zap.L().Sugar().Debugln("Failed to get secret version", zap.Error(err))
		models.EncodeError(res, "Failed to get secret version", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(struct {
		Version int `json:"version"`
	}{
		Version: result.Version,
	})
}
