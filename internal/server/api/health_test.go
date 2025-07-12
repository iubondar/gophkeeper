package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gophkeeper/internal/storage/mocks"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

// ExamplePingHandler_Ping демонстрирует пример использования эндпоинта проверки доступности сервиса.
// Пример показывает, как проверить работоспособность сервера.
func ExampleHealthHandler_Health() {
	// Создаем тестовый HTTP запрос
	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	// Создаем мок для проверки статуса
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()
	mockChecker := mocks.NewMockStatusChecker(ctrl)
	mockChecker.EXPECT().CheckStatus(gomock.Any()).Return(nil)

	// Инициализируем обработчик
	handler := NewHealthHandler(mockChecker)

	// Вызываем обработчик
	handler.Health(w, request)

	// Получаем ответ
	res := w.Result()
	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	// Выводим статус ответа
	fmt.Println(res.Status)
	// Output: 200 OK
}

func TestHealthHandler_Health(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		setErr   error
		wantCode int
	}{
		{
			name:     "Positive test",
			method:   http.MethodGet,
			setErr:   nil,
			wantCode: http.StatusOK,
		},
		{
			name:     "Test POST method not allowed",
			method:   http.MethodPost,
			setErr:   nil,
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			name:     "Test PUT method not allowed",
			method:   http.MethodPut,
			setErr:   nil,
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			name:     "Test DELETE method not allowed",
			method:   http.MethodDelete,
			setErr:   nil,
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			name:     "Check error test",
			method:   http.MethodGet,
			setErr:   errors.New("Status is not ok"),
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockStatusChecker(ctrl)

			m.EXPECT().CheckStatus(gomock.Any()).Return(tt.setErr).AnyTimes()

			request := httptest.NewRequest(tt.method, "/ping", nil)

			w := httptest.NewRecorder()

			handler := NewHealthHandler(m)
			handler.Health(w, request)

			res := w.Result()
			defer func() {
				if err := res.Body.Close(); err != nil {
					t.Errorf("Error closing response body: %v", err)
				}
			}()

			assert.Equal(t, tt.wantCode, res.StatusCode)
		})
	}
}
