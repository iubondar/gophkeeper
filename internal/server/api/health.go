// Package api предоставляет HTTP-обработчики для REST API сервера GophKeeper.
// Пакет содержит обработчики для всех основных операций: аутентификация,
// управление секретами, загрузка/скачивание файлов, проверка состояния сервиса.
// Все обработчики используют usecase-слой для бизнес-логики и возвращают
// структурированные JSON-ответы с соответствующими HTTP-статусами.
package api

import (
	"context"
	"net/http"

	"gophkeeper/internal/models"
)

// StatusChecker определяет интерфейс для проверки состояния хранилища.
// Используется для проверки доступности и работоспособности хранилища.
// Интерфейс позволяет абстрагироваться от конкретной реализации хранилища
// и тестировать обработчики с помощью моков.
type StatusChecker interface {
	// CheckStatus проверяет состояние хранилища.
	// Возвращает ошибку, если хранилище недоступно или неработоспособно.
	CheckStatus(ctx context.Context) error
}

// HealthHandler обрабатывает HTTP-запросы для проверки доступности сервиса.
// Используется для проверки работоспособности сервера и его подключения к хранилищу.
// Обработчик реализует endpoint /health для мониторинга состояния системы.
type HealthHandler struct {
	checker StatusChecker // интерфейс для проверки статуса хранилища
}

// NewHealthHandler создает новый экземпляр HealthHandler.
// Принимает интерфейс для проверки статуса хранилища.
// Функция используется для внедрения зависимостей и создания обработчика
// с конкретной реализацией проверки состояния системы.
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
