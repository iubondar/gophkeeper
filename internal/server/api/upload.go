package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"gophkeeper/internal/auth"
	"gophkeeper/internal/models"
	"gophkeeper/internal/server/usecase"

	"go.uber.org/zap"
)

// UploadHandler обрабатывает HTTP-запросы для загрузки зашифрованных секретов.
// Обработчик реализует endpoint /api/upload и требует аутентификации пользователя.
// Использует usecase-слой для сохранения секретов в хранилище.
type UploadHandler struct {
	uc usecase.UploadSecretUsecase
}

// NewUploadHandler создает новый экземпляр UploadHandler.
// Принимает usecase для загрузки секретов.
// Функция используется для внедрения зависимостей и создания обработчика
// с конкретной реализацией бизнес-логики загрузки секретов.
func NewUploadHandler(uc usecase.UploadSecretUsecase) *UploadHandler {
	return &UploadHandler{uc: uc}
}

// Upload godoc
// @Summary Загрузка секрета
// @Description Загружает зашифрованный секрет на сервер
// @Tags secrets
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.UploadSecretIn true "Данные секрета для загрузки"
// @Success 200 {object} models.UploadSecretOut "Секрет успешно загружен"
// @Failure 400 {object} models.JSONError "Некорректные данные запроса"
// @Failure 401 {object} models.JSONError "Неавторизованный доступ"
// @Failure 500 {object} models.JSONError "Внутренняя ошибка сервера"
// @Router /api/upload [post]
func (handler UploadHandler) Upload(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		models.EncodeError(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	userID, err := auth.GetUserIDFromReq(req)
	if err != nil {
		zap.L().Sugar().Debugln("Failed to get user ID", zap.Error(err))
		models.EncodeError(res, "Failed to get user ID", http.StatusUnauthorized)
		return
	}

	var in models.UploadSecretIn
	var buf bytes.Buffer
	_, err = buf.ReadFrom(req.Body)
	if err != nil {
		models.EncodeError(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &in); err != nil {
		models.EncodeError(res, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := handler.uc.UploadSecret(req.Context(), in, userID)
	if err != nil {
		zap.L().Sugar().Debugln("Failed to upload secret", zap.Error(err))
		models.EncodeError(res, "Failed to upload secret", http.StatusInternalServerError)
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
