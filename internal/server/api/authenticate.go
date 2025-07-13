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

type AuthenticateHandler struct {
	uc usecase.AuthenticateUsecase
}

func NewAuthenticateHandler(uc usecase.AuthenticateUsecase) *AuthenticateHandler {
	return &AuthenticateHandler{
		uc: uc,
	}
}

func (handler AuthenticateHandler) Authenticate(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	var in models.AuthenticateIn
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

	result, err := handler.uc.Authenticate(req.Context(), in.Login, in.PasswordHash)
	if err != nil {
		if errors.Is(err, usecase.ErrLoginOrPasswordEmpty) {
			zap.L().Sugar().Debugln("Login or password is empty", zap.Error(err))
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		if errors.Is(err, usecase.ErrUserNotFound) {
			zap.L().Sugar().Debugln("User not found", zap.Error(err))
			http.Error(res, err.Error(), http.StatusUnauthorized)
			return
		}

		zap.L().Sugar().Debugln("Failed to authenticate user", zap.Error(err))
		http.Error(res, "Failed to authenticate user", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(res).Encode(result); err != nil {
		zap.L().Sugar().Debugln("Failed to encode response", zap.Error(err))
		http.Error(res, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
