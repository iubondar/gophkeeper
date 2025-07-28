// Package api предоставляет HTTP-обработчики для REST API сервера GophKeeper.
// Пакет содержит обработчики для всех основных операций: аутентификация,
// управление секретами, загрузка/скачивание файлов, проверка состояния сервиса.
// Все обработчики используют usecase-слой для бизнес-логики и возвращают
// структурированные JSON-ответы с соответствующими HTTP-статусами.
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

// RegisterHandler обрабатывает HTTP-запросы для регистрации новых пользователей.
// Обработчик реализует endpoint /api/register и использует usecase-слой
// для выполнения бизнес-логики регистрации пользователей.
type RegisterHandler struct {
	uc usecase.RegisterUsecase
}

// NewRegisterHandler создает новый экземпляр RegisterHandler.
// Принимает usecase для регистрации пользователей.
// Функция используется для внедрения зависимостей и создания обработчика
// с конкретной реализацией бизнес-логики регистрации.
func NewRegisterHandler(uc usecase.RegisterUsecase) *RegisterHandler {
	return &RegisterHandler{
		uc: uc,
	}
}

// Register godoc
// @Summary Регистрация нового пользователя
// @Description Регистрирует нового пользователя в системе
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterIn true "Данные для регистрации"
// @Success 200 {object} models.AuthenticateOut "Успешная регистрация"
// @Failure 400 {object} models.JSONError "Некорректные данные запроса"
// @Failure 409 {object} models.JSONError "Пользователь уже существует"
// @Failure 500 {object} models.JSONError "Внутренняя ошибка сервера"
// @Router /api/register [post]
func (handler RegisterHandler) Register(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		models.EncodeError(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	var in models.RegisterIn
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

	result, err := handler.uc.Register(req.Context(), in)
	if err != nil {
		if errors.Is(err, models.ErrLoginOrPasswordEmpty) {
			zap.L().Sugar().Debugln("Login or password is empty", zap.Error(err))
			models.EncodeError(res, models.ErrLoginOrPasswordEmpty.Error(), http.StatusBadRequest)
			return
		}

		if errors.Is(err, models.ErrUserAlreadyExists) {
			zap.L().Sugar().Debugln("User already exists", zap.Error(err))
			models.EncodeError(res, models.ErrUserAlreadyExists.Error(), http.StatusConflict)
			return
		}

		zap.L().Sugar().Debugln("Failed to register user", zap.Error(err))
		models.EncodeError(res, "Failed to register user", http.StatusInternalServerError)
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
