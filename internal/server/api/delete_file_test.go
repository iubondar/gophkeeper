package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"gophkeeper/internal/auth"
	"gophkeeper/internal/models"
	"gophkeeper/internal/server/storage/mocks"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDeleteFileHandler_DeleteFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockDeleteFileUsecase(ctrl)
	handler := NewDeleteFileHandler(mockUsecase)

	userID := uuid.New()
	label := "test-file"

	t.Run("successful delete", func(t *testing.T) {
		// Создаем запрос с параметром label
		req := httptest.NewRequest(http.MethodDelete, "/api/files/"+label, nil)
		// Устанавливаем URL параметр для chi router
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"label"},
				Values: []string{label},
			},
		}))

		// Добавляем аутентификацию
		token, err := auth.GenerateAccessToken(userID.String())
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{
			Name:  auth.AuthCookieName,
			Value: token,
		})

		w := httptest.NewRecorder()

		// Ожидаем вызов usecase
		mockUsecase.EXPECT().
			DeleteFile(gomock.Any(), label, userID).
			Return(nil)

		// Выполняем запрос
		handler.DeleteFile(w, req)

		// Проверяем результат
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("wrong HTTP method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/files/"+label, nil)
		w := httptest.NewRecorder()

		handler.DeleteFile(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("missing user ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/files/"+label, nil)
		w := httptest.NewRecorder()

		handler.DeleteFile(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("missing label", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/files/", nil)
		// Устанавливаем пустой URL параметр для chi router
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"label"},
				Values: []string{""},
			},
		}))

		// Добавляем аутентификацию
		token, err := auth.GenerateAccessToken(userID.String())
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{
			Name:  auth.AuthCookieName,
			Value: token,
		})
		w := httptest.NewRecorder()

		handler.DeleteFile(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase not found error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/files/"+label, nil)
		// Устанавливаем URL параметр для chi router
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"label"},
				Values: []string{label},
			},
		}))

		// Добавляем аутентификацию
		token, err := auth.GenerateAccessToken(userID.String())
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{
			Name:  auth.AuthCookieName,
			Value: token,
		})
		w := httptest.NewRecorder()

		mockUsecase.EXPECT().
			DeleteFile(gomock.Any(), label, userID).
			Return(models.ErrRecordNotFound)

		handler.DeleteFile(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/files/"+label, nil)
		// Устанавливаем URL параметр для chi router
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"label"},
				Values: []string{label},
			},
		}))

		// Добавляем аутентификацию
		token, err := auth.GenerateAccessToken(userID.String())
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{
			Name:  auth.AuthCookieName,
			Value: token,
		})
		w := httptest.NewRecorder()

		mockUsecase.EXPECT().
			DeleteFile(gomock.Any(), label, userID).
			Return(assert.AnError)

		handler.DeleteFile(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
