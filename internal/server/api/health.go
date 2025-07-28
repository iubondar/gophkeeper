package api

import (
	"context"
	"net/http"

	"gophkeeper/internal/models"
)

// StatusChecker определяет интерфейс для проверки состояния хранилища.
// Используется для проверки доступности и работоспособности хранилища.
type StatusChecker interface {
	// CheckStatus проверяет состояние хранилища.
	// Возвращает ошибку, если хранилище недоступно или неработоспособно.
	CheckStatus(ctx context.Context) error
}

// HealthHandler обрабатывает запросы для проверки доступности сервиса.
// Используется для проверки работоспособности сервера и его подключения к хранилищу.
type HealthHandler struct {
	checker StatusChecker // интерфейс для проверки статуса хранилища
}

// NewHealthHandler создает новый экземпляр HealthHandler.
// Принимает интерфейс для проверки статуса хранилища.
func NewHealthHandler(checker StatusChecker) HealthHandler {
	return HealthHandler{
		checker: checker,
	}
}

// Health godoc
// @Summary Проверка здоровья сервиса
// @Description Проверяет доступность сервиса и подключение к хранилищу данных
// @Tags health
// @Accept json
// @Produce json
// @Success 200 "Сервис работает нормально"
// @Failure 500 {object} models.JSONError "Сервис недоступен"
// @Router /health [get]
func (handler HealthHandler) Health(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		models.EncodeError(res, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	err := handler.checker.CheckStatus(req.Context())
	if err != nil {
		models.EncodeError(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
}
