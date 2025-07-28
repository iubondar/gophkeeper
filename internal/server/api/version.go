package api

import (
	"encoding/json"
	"net/http"

	"gophkeeper/internal/auth"
	"gophkeeper/internal/models"
	"gophkeeper/internal/server/usecase"

	"go.uber.org/zap"
)

// VersionHandler обрабатывает запросы получения версий секретов
type VersionHandler struct {
	uc usecase.GetSecretUsecase
}

// NewVersionHandler создает новый экземпляр VersionHandler
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

	userID, err := auth.GetUserIDFromReq(req)
	if err != nil {
		zap.L().Sugar().Debugln("Failed to get user ID", zap.Error(err))
		models.EncodeError(res, "Failed to get user ID", http.StatusUnauthorized)
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
