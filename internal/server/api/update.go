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

type UpdateHandler struct {
	uc usecase.UpdateSecretUsecase
}

func NewUpdateHandler(uc usecase.UpdateSecretUsecase) *UpdateHandler {
	return &UpdateHandler{uc: uc}
}

func (handler UpdateHandler) Update(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		models.EncodeError(res, "Only PUT requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	userID, err := auth.GetUserIDFromReq(req)
	if err != nil {
		zap.L().Sugar().Debugln("Failed to get user ID", zap.Error(err))
		models.EncodeError(res, "Failed to get user ID", http.StatusUnauthorized)
		return
	}

	var in models.UpdateSecretIn
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

	result, err := handler.uc.UpdateSecret(req.Context(), in, userID)
	if err != nil {
		if err == models.ErrConflict {
			models.EncodeError(res, "Version conflict", http.StatusConflict)
			return
		}
		zap.L().Sugar().Debugln("Failed to update secret", zap.Error(err))
		models.EncodeError(res, "Failed to update secret", http.StatusInternalServerError)
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
