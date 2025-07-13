package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"gophkeeper/internal/models"
	"gophkeeper/internal/server/usecase"

	"go.uber.org/zap"
)

type LoginHandler struct {
	uc usecase.LoginUsecase
}

func NewLoginHandler(uc usecase.LoginUsecase) *LoginHandler {
	return &LoginHandler{
		uc: uc,
	}
}

func (handler LoginHandler) Login(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	var in models.LoginIn
	var buf bytes.Buffer
	// читаем тело запроса
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// десериализуем JSON
	if err = json.Unmarshal(buf.Bytes(), &in); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	salt, err := handler.uc.GetSalt(req.Context(), in.Login)
	if err != nil {
		zap.L().Sugar().Debugln("Failed to get user salt", zap.Error(err))
		http.Error(res, "Failed to get user salt", http.StatusInternalServerError)
		return
	}

	// Если соль пустая, значит пользователь не найден
	if salt == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	response := models.LoginOut{Salt: salt}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(res).Encode(response); err != nil {
		zap.L().Sugar().Debugln("Failed to encode response", zap.Error(err))
		http.Error(res, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
