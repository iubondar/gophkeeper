package api

import (
	"encoding/json"
	"net/http"

	"gophkeeper/internal/auth"
	"gophkeeper/internal/models"
	"gophkeeper/internal/server/usecase"

	"go.uber.org/zap"
)

type VersionHandler struct {
	uc usecase.GetSecretUsecase
}

func NewVersionHandler(uc usecase.GetSecretUsecase) *VersionHandler {
	return &VersionHandler{uc: uc}
}

// GetSecretVersionHandler возвращает только версию секрета по имени
func (handler VersionHandler) GetSecretVersion(res http.ResponseWriter, req *http.Request) {
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

	secretName := req.URL.Query().Get("name")
	if secretName == "" {
		models.EncodeError(res, "Secret name is required", http.StatusBadRequest)
		return
	}

	result, err := handler.uc.GetSecret(req.Context(), secretName, userID)
	if err != nil {
		zap.L().Sugar().Debugln("Failed to get secret version", zap.Error(err))
		models.EncodeError(res, "Failed to get secret version", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(struct {
		Version int `json:"version"`
	}{
		Version: result.Version,
	})
}
