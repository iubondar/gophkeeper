package api

import (
	"encoding/json"
	"net/http"

	"gophkeeper/internal/auth"
	"gophkeeper/internal/models"
	"gophkeeper/internal/server/usecase"

	"go.uber.org/zap"
)

// GetHandler обрабатывает HTTP-запросы для получения зашифрованных секретов.
// Обработчик реализует endpoint /api/get и требует аутентификации пользователя.
// Использует usecase-слой для получения секретов из хранилища по имени.
type GetHandler struct {
	uc usecase.GetSecretUsecase
}

// NewGetHandler создает новый экземпляр GetHandler.
// Принимает usecase для получения секретов.
// Функция используется для внедрения зависимостей и создания обработчика
// с конкретной реализацией бизнес-логики получения секретов.
func NewGetHandler(uc usecase.GetSecretUsecase) *GetHandler {
	return &GetHandler{uc: uc}
}

// GetSecret godoc
// @Summary Получение секрета
// @Description Получает зашифрованный секрет с сервера по имени
// @Tags secrets
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query string true "Имя секрета"
// @Success 200 {object} models.GetSecretOut "Секрет успешно получен"
// @Failure 400 {object} models.JSONError "Некорректные данные запроса"
// @Failure 401 {object} models.JSONError "Неавторизованный доступ"
// @Failure 404 {object} models.JSONError "Секрет не найден"
// @Failure 500 {object} models.JSONError "Внутренняя ошибка сервера"
// @Router /api/get [get]
func (handler GetHandler) GetSecret(res http.ResponseWriter, req *http.Request) {
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

	// Получаем имя секрета из query параметра
	secretName := req.URL.Query().Get("name")
	if secretName == "" {
		models.EncodeError(res, "Secret name is required", http.StatusBadRequest)
		return
	}

	result, err := handler.uc.GetSecret(req.Context(), secretName, userID)
	if err != nil {
		zap.L().Sugar().Debugln("Failed to get secret", zap.Error(err))
		models.EncodeError(res, "Failed to get secret", http.StatusInternalServerError)
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
