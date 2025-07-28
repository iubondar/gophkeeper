package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"gophkeeper/internal/models"
	"gophkeeper/internal/server/usecase"

	"go.uber.org/zap"
)

// LoginHandler обрабатывает HTTP-запросы для входа пользователей в систему.
// Обработчик реализует endpoint /api/login и возвращает соль для хеширования пароля.
// Использует usecase-слой для получения соли пользователя из хранилища.
type LoginHandler struct {
	uc usecase.LoginUsecase
}

// NewLoginHandler создает новый экземпляр LoginHandler.
// Принимает usecase для входа пользователей.
// Функция используется для внедрения зависимостей и создания обработчика
// с конкретной реализацией бизнес-логики входа.
func NewLoginHandler(uc usecase.LoginUsecase) *LoginHandler {
	return &LoginHandler{
		uc: uc,
	}
}

// Login godoc
// @Summary Вход пользователя в систему
// @Description Возвращает соль для хеширования пароля
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginIn true "Данные для входа"
// @Success 200 {object} models.LoginOut "Успешный вход, возвращена соль"
// @Failure 400 {object} models.JSONError "Некорректные данные запроса"
// @Failure 404 {object} models.JSONError "Пользователь не найден"
// @Failure 500 {object} models.JSONError "Внутренняя ошибка сервера"
// @Router /api/login [post]
func (handler LoginHandler) Login(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		models.EncodeError(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	var in models.LoginIn
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

	salt, err := handler.uc.GetSalt(req.Context(), in.Login)
	if err != nil {
		zap.L().Sugar().Debugln("Failed to get user salt", zap.Error(err))
		models.EncodeError(res, "Failed to get user salt", http.StatusInternalServerError)
		return
	}

	// Если соль пустая, значит пользователь не найден
	if salt == "" {
		models.EncodeError(res, models.ErrUserNotFound.Error(), http.StatusNotFound)
		return
	}

	response := models.LoginOut{Salt: salt}
	res.Header().Set("Content-Type", "application/json")

	res.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(res).Encode(response); err != nil {
		zap.L().Sugar().Debugln("Failed to encode response", zap.Error(err))
		models.EncodeError(res, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
