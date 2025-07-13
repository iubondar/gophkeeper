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

type RegisterHandler struct {
	uc usecase.RegisterUsecase
}

func NewRegisterHandler(uc usecase.RegisterUsecase) *RegisterHandler {
	return &RegisterHandler{
		uc: uc,
	}
}

func (handler RegisterHandler) Register(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	var in models.RegisterIn
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

	result, err := handler.uc.Register(req.Context(), in)
	if err != nil {
		if errors.Is(err, models.ErrLoginOrPasswordEmpty) {
			zap.L().Sugar().Debugln("Login or password is empty", zap.Error(err))
			http.Error(res, models.ErrLoginOrPasswordEmpty.Error(), http.StatusBadRequest)
			return
		}

		if errors.Is(err, models.ErrUserAlreadyExists) {
			zap.L().Sugar().Debugln("User already exists", zap.Error(err))
			http.Error(res, models.ErrUserAlreadyExists.Error(), http.StatusConflict)
			return
		}

		zap.L().Sugar().Debugln("Failed to register user", zap.Error(err))
		http.Error(res, "Failed to register user", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(res).Encode(result); err != nil {
		zap.L().Sugar().Debugln("Failed to encode response", zap.Error(err))
		http.Error(res, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}
