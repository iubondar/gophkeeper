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

type UploadHandler struct {
	uc usecase.UploadSecretUsecase
}

func NewUploadHandler(uc usecase.UploadSecretUsecase) *UploadHandler {
	return &UploadHandler{uc: uc}
}

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
