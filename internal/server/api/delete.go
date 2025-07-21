package api

import (
	"net/http"

	"gophkeeper/internal/auth"
	"gophkeeper/internal/models"
	"gophkeeper/internal/server/usecase"

	"go.uber.org/zap"
)

type DeleteHandler struct {
	uc usecase.DeleteSecretUsecase
}

func NewDeleteHandler(uc usecase.DeleteSecretUsecase) *DeleteHandler {
	return &DeleteHandler{uc: uc}
}

func (handler DeleteHandler) DeleteSecret(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		models.EncodeError(res, "Only DELETE requests are allowed!", http.StatusMethodNotAllowed)
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

	err = handler.uc.DeleteSecret(req.Context(), secretName, userID)
	if err != nil {
		if err == models.ErrRecordNotFound {
			models.EncodeError(res, "Secret not found", http.StatusNotFound)
			return
		}
		zap.L().Sugar().Debugln("Failed to delete secret", zap.Error(err))
		models.EncodeError(res, "Failed to delete secret", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}
