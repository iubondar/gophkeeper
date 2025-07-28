package api

import (
	"net/http"

	"gophkeeper/internal/auth"
	"gophkeeper/internal/models"
	"gophkeeper/internal/server/usecase"

	"go.uber.org/zap"
)

// DeleteHandler обрабатывает HTTP-запросы для удаления секретов.
// Обработчик реализует endpoint /api/delete и требует аутентификации пользователя.
// Использует usecase-слой для удаления секретов из хранилища по имени.
type DeleteHandler struct {
	uc usecase.DeleteSecretUsecase
}

// NewDeleteHandler создает новый экземпляр DeleteHandler.
// Принимает usecase для удаления секретов.
// Функция используется для внедрения зависимостей и создания обработчика
// с конкретной реализацией бизнес-логики удаления секретов.
func NewDeleteHandler(uc usecase.DeleteSecretUsecase) *DeleteHandler {
	return &DeleteHandler{uc: uc}
}

// DeleteSecret godoc
// @Summary Удаление секрета
// @Description Удаляет секрет с сервера по имени
// @Tags secrets
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query string true "Имя секрета для удаления"
// @Success 204 "Секрет успешно удален"
// @Failure 400 {object} models.JSONError "Некорректные данные запроса"
// @Failure 401 {object} models.JSONError "Неавторизованный доступ"
// @Failure 404 {object} models.JSONError "Секрет не найден"
// @Failure 500 {object} models.JSONError "Внутренняя ошибка сервера"
// @Router /api/delete [delete]
func (handler DeleteHandler) DeleteSecret(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		models.EncodeError(res, "Only DELETE requests are allowed!", http.StatusMethodNotAllowed)
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

	err = handler.uc.DeleteSecret(req.Context(), secretName, userID)
	if err != nil {
		if err == models.ErrRecordNotFound {
			models.EncodeError(res, "Secret not found", http.StatusNotFound)
			return
		}
		zap.L().Sugar().Debugln("Failed to delete secret", zap.Error(err))
		models.EncodeError(res, "Failed to delete secret", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}
