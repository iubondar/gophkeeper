package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"gophkeeper/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockDeleteSecretUsecase struct {
	deleteFunc func(ctx context.Context, secretName string, userID uuid.UUID) error
}

func (m *mockDeleteSecretUsecase) DeleteSecret(ctx context.Context, secretName string, userID uuid.UUID) error {
	return m.deleteFunc(ctx, secretName, userID)
}

func TestDeleteHandler_DeleteSecret(t *testing.T) {
	userID := uuid.New()
	secretName := "test-secret"

	t.Run("success", func(t *testing.T) {
		uc := &mockDeleteSecretUsecase{
			deleteFunc: func(ctx context.Context, name string, uid uuid.UUID) error {
				assert.Equal(t, secretName, name)
				assert.Equal(t, userID, uid)
				return nil
			},
		}
		h := NewDeleteHandler(uc)
		req := httptest.NewRequest(http.MethodDelete, "/api/delete?name="+secretName, nil)
		addUserIDCookie(req, userID)
		w := httptest.NewRecorder()
		h.DeleteSecret(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("not found", func(t *testing.T) {
		uc := &mockDeleteSecretUsecase{
			deleteFunc: func(ctx context.Context, name string, uid uuid.UUID) error {
				return models.ErrRecordNotFound
			},
		}
		h := NewDeleteHandler(uc)
		req := httptest.NewRequest(http.MethodDelete, "/api/delete?name="+secretName, nil)
		addUserIDCookie(req, userID)
		w := httptest.NewRecorder()
		h.DeleteSecret(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("wrong method", func(t *testing.T) {
		uc := &mockDeleteSecretUsecase{}
		h := NewDeleteHandler(uc)
		req := httptest.NewRequest(http.MethodPost, "/api/delete?name="+secretName, nil)
		addUserIDCookie(req, userID)
		w := httptest.NewRecorder()
		h.DeleteSecret(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("missing name", func(t *testing.T) {
		uc := &mockDeleteSecretUsecase{}
		h := NewDeleteHandler(uc)
		req := httptest.NewRequest(http.MethodDelete, "/api/delete", nil)
		addUserIDCookie(req, userID)
		w := httptest.NewRecorder()
		h.DeleteSecret(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("missing userID", func(t *testing.T) {
		uc := &mockDeleteSecretUsecase{}
		h := NewDeleteHandler(uc)
		req := httptest.NewRequest(http.MethodDelete, "/api/delete?name="+secretName, nil)
		w := httptest.NewRecorder()
		h.DeleteSecret(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func addUserIDCookie(req *http.Request, userID uuid.UUID) {
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: userID.String(),
	})
}
