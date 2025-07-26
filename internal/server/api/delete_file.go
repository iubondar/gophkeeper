package api

import (
	"context"
	"net/http"

	"gophkeeper/internal/auth"
	"gophkeeper/internal/models"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// DeleteFileUsecase интерфейс для удаления файлов
type DeleteFileUsecase interface {
	DeleteFile(ctx context.Context, label string, userID uuid.UUID) error
}

type DeleteFileHandler struct {
	uc DeleteFileUsecase
}

func NewDeleteFileHandler(uc DeleteFileUsecase) *DeleteFileHandler {
	return &DeleteFileHandler{uc: uc}
}

// DeleteFile обрабатывает удаление файла
func (h *DeleteFileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		models.EncodeError(w, "Only DELETE requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	userID, err := auth.GetUserIDFromReq(r)
	if err != nil {
		zap.L().Sugar().Debugln("Failed to get user ID", zap.Error(err))
		models.EncodeError(w, "Failed to get user ID", http.StatusUnauthorized)
		return
	}

	// Получаем label из URL
	label := chi.URLParam(r, "label")
	if label == "" {
		models.EncodeError(w, "Label is required", http.StatusBadRequest)
		return
	}

	// Удаляем файл
	err = h.uc.DeleteFile(r.Context(), label, userID)
	if err != nil {
		if err == models.ErrRecordNotFound {
			models.EncodeError(w, "File not found", http.StatusNotFound)
			return
		}
		zap.L().Sugar().Debugln("Failed to delete file", zap.Error(err))
		models.EncodeError(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
