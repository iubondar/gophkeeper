package api

import (
	"encoding/json"
	"net/http"

	"gophkeeper/internal/auth"
	"gophkeeper/internal/models"
	"gophkeeper/internal/server/usecase"

	"go.uber.org/zap"
)

type GetHandler struct {
	uc usecase.GetSecretUsecase
}

func NewGetHandler(uc usecase.GetSecretUsecase) *GetHandler {
	return &GetHandler{uc: uc}
}

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
	if err = json.NewEncoder(res).Encode(result); err != nil {
		zap.L().Sugar().Debugln("Failed to encode response", zap.Error(err))
		models.EncodeError(res, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
