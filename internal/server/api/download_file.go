package api

import (
	"context"
	"io"
	"net/http"

	"gophkeeper/internal/auth"
	"gophkeeper/internal/models"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// DownloadFileUsecase интерфейс для скачивания файлов
type DownloadFileUsecase interface {
	DownloadFile(ctx context.Context, label string, userID uuid.UUID) (io.ReadCloser, error)
}

// DownloadFileHandler обрабатывает запросы скачивания файлов
type DownloadFileHandler struct {
	uc DownloadFileUsecase
}

// NewDownloadFileHandler создает новый экземпляр DownloadFileHandler
func NewDownloadFileHandler(uc DownloadFileUsecase) *DownloadFileHandler {
	return &DownloadFileHandler{uc: uc}
}

// DownloadFile godoc
// @Summary Скачивание файла
// @Description Скачивает зашифрованный файл с сервера по метке
// @Tags files
// @Accept json
// @Produce application/octet-stream
// @Security ApiKeyAuth
// @Param label path string true "Метка файла"
// @Success 200 {file} file "Файл успешно скачан"
// @Failure 400 {object} models.JSONError "Некорректные данные запроса"
// @Failure 401 {object} models.JSONError "Неавторизованный доступ"
// @Failure 404 {object} models.JSONError "Файл не найден"
// @Failure 500 {object} models.JSONError "Внутренняя ошибка сервера"
// @Router /api/files/{label}/download [get]
func (h *DownloadFileHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		models.EncodeError(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
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

	// Скачиваем файл
	reader, err := h.uc.DownloadFile(r.Context(), label, userID)
	if err != nil {
		if err == models.ErrRecordNotFound {
			models.EncodeError(w, "File not found", http.StatusNotFound)
			return
		}
		zap.L().Sugar().Debugln("Failed to download file", zap.Error(err))
		models.EncodeError(w, "Failed to download file", http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	// Отправляем файл
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+label)

	// Копируем данные из reader в response writer
	_, err = io.Copy(w, reader)
	if err != nil {
		zap.L().Sugar().Debugln("Failed to copy file data", zap.Error(err))
		// Не отправляем ошибку клиенту, так как ответ уже начался
		return
	}
}
