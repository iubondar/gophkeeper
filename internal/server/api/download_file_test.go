package api

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gophkeeper/internal/models"
	"gophkeeper/internal/server/storage/mocks"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDownloadFileHandler_DownloadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockDownloadFileUsecase(ctrl)
	handler := NewDownloadFileHandler(mockUsecase)

	userID := uuid.New()
	label := "test-file"
	content := "test content"

	t.Run("successful download", func(t *testing.T) {
		// Создаем запрос с параметром label
		req := httptest.NewRequest(http.MethodGet, "/api/files/"+label+"/download", nil)
		// Устанавливаем URL параметр для chi router
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"label"},
				Values: []string{label},
			},
		}))

		// Добавляем userID в контекст
		req = setUserIDInContext(req, userID)

		w := httptest.NewRecorder()

		// Ожидаем вызов usecase
		mockUsecase.EXPECT().
			DownloadFile(gomock.Any(), label, userID).
			Return(io.NopCloser(strings.NewReader(content)), nil)

		// Выполняем запрос
		handler.DownloadFile(w, req)

		// Проверяем результат
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/octet-stream", w.Header().Get("Content-Type"))
		assert.Equal(t, "attachment; filename="+label, w.Header().Get("Content-Disposition"))
		assert.Equal(t, content, w.Body.String())
	})

	t.Run("wrong HTTP method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/files/"+label+"/download", nil)
		w := httptest.NewRecorder()

		handler.DownloadFile(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("missing user ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/files/"+label+"/download", nil)
		w := httptest.NewRecorder()

		handler.DownloadFile(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("missing label", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/files//download", nil)
		// Устанавливаем пустой URL параметр для chi router
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"label"},
				Values: []string{""},
			},
		}))

		// Добавляем userID в контекст
		req = setUserIDInContext(req, userID)

		w := httptest.NewRecorder()

		handler.DownloadFile(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase not found error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/files/"+label+"/download", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"label"},
				Values: []string{label},
			},
		}))

		// Добавляем userID в контекст
		req = setUserIDInContext(req, userID)

		w := httptest.NewRecorder()

		mockUsecase.EXPECT().
			DownloadFile(gomock.Any(), label, userID).
			Return(nil, models.ErrRecordNotFound)

		handler.DownloadFile(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/files/"+label+"/download", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"label"},
				Values: []string{label},
			},
		}))

		// Добавляем userID в контекст
		req = setUserIDInContext(req, userID)

		w := httptest.NewRecorder()

		mockUsecase.EXPECT().
			DownloadFile(gomock.Any(), label, userID).
			Return(nil, assert.AnError)

		handler.DownloadFile(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
