package api

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"gophkeeper/internal/server/middleware"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// setUserIDInContext устанавливает userID в контекст запроса для тестирования
func setUserIDInContext(req *http.Request, userID uuid.UUID) *http.Request {
	ctx := context.WithValue(req.Context(), middleware.UserIDKey{}, userID)
	return req.WithContext(ctx)
}

// mustMarshal маршалит объект в JSON для тестов
func mustMarshal(t *testing.T, v any) []byte {
	t.Helper()
	data, err := json.Marshal(v)
	assert.NoError(t, err)
	return data
}
