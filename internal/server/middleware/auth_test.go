package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gophkeeper/internal/auth"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		setupRequest   func() *http.Request
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				// Создаем валидный токен
				token, err := auth.GenerateAccessToken(uuid.New().String())
				require.NoError(t, err)
				req.AddCookie(&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: token,
				})
				return req
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name: "Missing token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"message":"Authentication required","code":401}` + "\n",
		},
		{
			name: "Invalid token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.AddCookie(&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: "invalid-token",
				})
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"message":"Malformed authentication token","code":401}` + "\n",
		},
		{
			name: "Empty token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.AddCookie(&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: "",
				})
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"message":"Malformed authentication token","code":401}` + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый хэндлер, который будет вызван после middleware
			handlerCalled := false
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			})

			// Создаем middleware
			middleware := AuthMiddleware(testHandler)

			// Создаем запрос
			req := tt.setupRequest()
			w := httptest.NewRecorder()

			// Выполняем запрос
			middleware.ServeHTTP(w, req)

			// Проверяем результат
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())

			// Проверяем, был ли вызван хэндлер
			if tt.expectedStatus == http.StatusOK {
				assert.True(t, handlerCalled, "Handler should have been called")
			} else {
				assert.False(t, handlerCalled, "Handler should not have been called")
			}
		})
	}
}

func TestAuthMiddleware_WithExpiredToken(t *testing.T) {
	// Создаем токен с истекшим временем жизни
	userID := uuid.New().String()

	// Создаем claims с истекшим временем
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Токен истек час назад
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
	}

	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("supersecretkey"))
	require.NoError(t, err)

	// Создаем запрос с истекшим токеном
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  auth.AuthCookieName,
		Value: tokenString,
	})

	// Создаем тестовый хэндлер
	handlerCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Создаем middleware
	middleware := AuthMiddleware(testHandler)
	w := httptest.NewRecorder()

	// Выполняем запрос
	middleware.ServeHTTP(w, req)

	// Проверяем, что запрос был отклонен из-за истекшего токена
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authentication token expired")
	assert.False(t, handlerCalled, "Handler should not have been called")
}

func TestAuthMiddleware_WithFutureToken(t *testing.T) {
	// Создаем токен с будущим временем начала действия
	userID := uuid.New().String()

	// Создаем claims с будущим временем начала действия
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		NotBefore: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // Токен начнет действовать через час
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
	}

	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("supersecretkey"))
	require.NoError(t, err)

	// Создаем запрос с токеном, который еще не действует
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  auth.AuthCookieName,
		Value: tokenString,
	})

	// Создаем тестовый хэндлер
	handlerCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Создаем middleware
	middleware := AuthMiddleware(testHandler)
	w := httptest.NewRecorder()

	// Выполняем запрос
	middleware.ServeHTTP(w, req)

	// Проверяем, что запрос был отклонен из-за токена, который еще не действует
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authentication token not valid yet")
	assert.False(t, handlerCalled, "Handler should not have been called")
}

func TestAuthMiddleware_ExpiredTokenMessage(t *testing.T) {
	// Создаем токен с истекшим временем жизни
	userID := uuid.New().String()

	// Создаем claims с истекшим временем
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Токен истек час назад
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
	}

	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("supersecretkey"))
	require.NoError(t, err)

	// Создаем запрос с истекшим токеном
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  auth.AuthCookieName,
		Value: tokenString,
	})

	// Создаем тестовый хэндлер
	handlerCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Создаем middleware
	middleware := AuthMiddleware(testHandler)
	w := httptest.NewRecorder()

	// Выполняем запрос
	middleware.ServeHTTP(w, req)

	// Проверяем, что запрос был отклонен с правильным сообщением
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authentication token expired")
	assert.False(t, handlerCalled, "Handler should not have been called")
}

func TestAuthMiddleware_SpecificJWTErrors(t *testing.T) {
	tests := []struct {
		name           string
		tokenClaims    jwt.RegisteredClaims
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Expired token",
			tokenClaims: jwt.RegisteredClaims{
				Subject:   uuid.New().String(),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"message":"Authentication token expired","code":401}` + "\n",
		},
		{
			name: "Future token",
			tokenClaims: jwt.RegisteredClaims{
				Subject:   uuid.New().String(),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
				NotBefore: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"message":"Authentication token not valid yet","code":401}` + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем токен с указанными claims
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, tt.tokenClaims)
			tokenString, err := token.SignedString([]byte("supersecretkey"))
			require.NoError(t, err)

			// Создаем запрос с токеном
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.AddCookie(&http.Cookie{
				Name:  auth.AuthCookieName,
				Value: tokenString,
			})

			// Создаем тестовый хэндлер
			handlerCalled := false
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			})

			// Создаем middleware
			middleware := AuthMiddleware(testHandler)
			w := httptest.NewRecorder()

			// Выполняем запрос
			middleware.ServeHTTP(w, req)

			// Проверяем результат
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
			assert.False(t, handlerCalled, "Handler should not have been called")
		})
	}
}
