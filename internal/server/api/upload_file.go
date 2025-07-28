package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"gophkeeper/internal/auth"
	"gophkeeper/internal/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UploadFileUsecase интерфейс для загрузки файлов
type UploadFileUsecase interface {
	UploadFile(ctx context.Context, label, metadata, fileName string, file io.Reader, size int64, userID uuid.UUID) (models.UploadSecretOut, error)
}

// UploadFileHandler обрабатывает запросы загрузки файлов
type UploadFileHandler struct {
	uc UploadFileUsecase
}

// NewUploadFileHandler создает новый экземпляр UploadFileHandler
func NewUploadFileHandler(uc UploadFileUsecase) *UploadFileHandler {
	return &UploadFileHandler{uc: uc}
}

// UploadFile godoc
// @Summary Загрузка файла
// @Description Загружает зашифрованный файл на сервер
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param file formData file true "Файл для загрузки"
// @Param label formData string true "Метка файла"
// @Param metadata formData string false "Метаданные файла"
// @Success 200 {object} models.UploadSecretOut "Файл успешно загружен"
// @Failure 400 {object} models.JSONError "Некорректные данные запроса"
// @Failure 401 {object} models.JSONError "Неавторизованный доступ"
// @Failure 409 {object} models.JSONError "Файл с такой меткой уже существует"
// @Failure 500 {object} models.JSONError "Внутренняя ошибка сервера"
// @Router /api/files [post]
func (h *UploadFileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		models.EncodeError(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	userID, err := auth.GetUserIDFromReq(r)
	if err != nil {
		zap.L().Sugar().Debugln("Failed to get user ID", zap.Error(err))
		models.EncodeError(w, "Failed to get user ID", http.StatusUnauthorized)
		return
	}

	// Парсим multipart форму
	err = r.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		models.EncodeError(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Получаем поля формы
	label := r.FormValue("label")
	if label == "" {
		models.EncodeError(w, "Label is required", http.StatusBadRequest)
		return
	}

	metadata := r.FormValue("metadata")

	// Получаем файл
	file, header, err := r.FormFile("file")
	if err != nil {
		models.EncodeError(w, "File is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Загружаем файл
	result, err := h.uc.UploadFile(r.Context(), label, metadata, header.Filename, file, header.Size, userID)
	if err != nil {
		if err == models.ErrConflict {
			models.EncodeError(w, "File with this label already exists", http.StatusConflict)
			return
		}
		zap.L().Sugar().Debugln("Failed to upload file", zap.Error(err))
		models.EncodeError(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		zap.L().Sugar().Debugln("Failed to encode response", zap.Error(err))
		models.EncodeError(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
