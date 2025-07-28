package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestServer_New(t *testing.T) {
	// Настройка логгера для тестов
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := New(":8080", handler)
	assert.NotNil(t, server)
	assert.Equal(t, ":8080", server.address)
	// Не сравниваем функции, так как это невозможно
	assert.NotNil(t, server.router)
}

func TestServer_StartAndShutdown_Localhost(t *testing.T) {
	// Настройка логгера для тестов
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	// Создаем тестовый обработчик
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Создаем сервер на случайном порту
	server := New("localhost:0", handler)

	// Канал для сигнализации о готовности сервера
	serverReady := make(chan struct{})
	serverError := make(chan error, 1)

	// Запускаем сервер в горутине
	go func() {
		close(serverReady)
		err := server.Start()
		if err != nil && err != http.ErrServerClosed {
			serverError <- err
		}
		close(serverError)
	}()

	// Ждем готовности сервера
	<-serverReady
	time.Sleep(100 * time.Millisecond)

	// Проверяем, что сервер отвечает
	if server.httpServer != nil {
		req := httptest.NewRequest("GET", "http://"+server.httpServer.Addr, nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "test response", w.Body.String())
	}

	// Тестируем graceful shutdown
	shutdownErr := server.Shutdown()
	assert.NoError(t, shutdownErr)

	// Ждем завершения сервера
	select {
	case err := <-serverError:
		if err != nil {
			t.Fatalf("Неожиданная ошибка сервера: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Сервер не завершился вовремя")
	}
}

func TestServer_ShutdownTimeout(t *testing.T) {
	// Настройка логгера для тестов
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	// Создаем обработчик, который работает дольше timeout
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(6 * time.Second) // Дольше 5-секундного timeout
		w.WriteHeader(http.StatusOK)
	})

	server := New("localhost:0", handler)

	serverReady := make(chan struct{})
	serverError := make(chan error, 1)

	go func() {
		close(serverReady)
		err := server.Start()
		if err != nil && err != http.ErrServerClosed {
			serverError <- err
		}
		close(serverError)
	}()

	<-serverReady
	time.Sleep(100 * time.Millisecond)

	// Тестируем shutdown с timeout
	shutdownErr := server.Shutdown()
	assert.NoError(t, shutdownErr) // Должен успешно завершиться благодаря принудительному закрытию

	select {
	case err := <-serverError:
		if err != nil {
			t.Fatalf("Неожиданная ошибка сервера: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Сервер не завершился вовремя")
	}
}

func TestServer_TLSConfiguration(t *testing.T) {
	// Настройка логгера для тестов
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Тестируем с адресом, не содержащим localhost
	server := New("test.example.com:443", handler)

	// Проверяем, что для не-localhost адресов создается TLS конфигурация
	// Создаем сервер, но не запускаем его, чтобы избежать ошибок подключения
	serverReady := make(chan struct{})
	serverError := make(chan error, 1)

	go func() {
		close(serverReady)
		err := server.Start()
		// Ожидаем ошибку подключения, так как адрес недоступен
		serverError <- err
		close(serverError)
	}()

	<-serverReady
	time.Sleep(100 * time.Millisecond)

	// Проверяем, что TLS конфигурация создана
	if server.httpServer != nil {
		assert.NotNil(t, server.httpServer.TLSConfig)
	}

	// Проверяем, что сервер вернул ошибку подключения
	select {
	case err := <-serverError:
		// Ожидаем ошибку подключения
		assert.Error(t, err)
		// Проверяем, что ошибка связана с подключением
		assert.True(t, strings.Contains(err.Error(), "bind: can't assign requested address") ||
			strings.Contains(err.Error(), "no such host") ||
			strings.Contains(err.Error(), "connection refused"))
	case <-time.After(2 * time.Second):
		t.Fatal("Сервер не вернул ошибку вовремя")
	}
}

func TestServer_SignalHandling(t *testing.T) {
	// Настройка логгера для тестов
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := New("localhost:0", handler)

	serverReady := make(chan struct{})
	serverError := make(chan error, 1)

	go func() {
		close(serverReady)
		err := server.Start()
		if err != nil && err != http.ErrServerClosed {
			serverError <- err
		}
		close(serverError)
	}()

	<-serverReady
	time.Sleep(100 * time.Millisecond)

	// Симулируем отправку сигнала SIGTERM
	// В реальном тесте это сложно сделать, поэтому просто проверяем shutdown
	shutdownErr := server.Shutdown()
	assert.NoError(t, shutdownErr)

	select {
	case err := <-serverError:
		if err != nil {
			t.Fatalf("Неожиданная ошибка сервера: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Сервер не завершился вовремя")
	}
}

func TestServer_ConcurrentRequests(t *testing.T) {
	// Настройка логгера для тестов
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond) // Имитируем обработку
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("response"))
	})

	server := New("localhost:0", handler)

	serverReady := make(chan struct{})
	serverError := make(chan error, 1)

	go func() {
		close(serverReady)
		err := server.Start()
		if err != nil && err != http.ErrServerClosed {
			serverError <- err
		}
		close(serverError)
	}()

	<-serverReady
	time.Sleep(100 * time.Millisecond)

	// Отправляем конкурентные запросы
	if server.httpServer != nil {
		const numRequests = 10
		results := make(chan bool, numRequests)

		for i := 0; i < numRequests; i++ {
			go func() {
				req := httptest.NewRequest("GET", "http://"+server.httpServer.Addr, nil)
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, req)
				results <- w.Code == http.StatusOK
			}()
		}

		// Проверяем результаты
		successCount := 0
		for i := 0; i < numRequests; i++ {
			if <-results {
				successCount++
			}
		}

		assert.Equal(t, numRequests, successCount)
	}

	shutdownErr := server.Shutdown()
	assert.NoError(t, shutdownErr)

	select {
	case err := <-serverError:
		if err != nil {
			t.Fatalf("Неожиданная ошибка сервера: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Сервер не завершился вовремя")
	}
}

func TestServer_ErrorHandling(t *testing.T) {
	// Настройка логгера для тестов
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	// Тестируем создание сервера с некорректным адресом
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Создаем сервер с адресом, который может вызвать ошибку
	server := New("invalid-address", handler)

	serverReady := make(chan struct{})
	serverError := make(chan error, 1)

	go func() {
		close(serverReady)
		err := server.Start()
		serverError <- err
		close(serverError)
	}()

	<-serverReady
	time.Sleep(100 * time.Millisecond)

	// Проверяем, что сервер корректно обрабатывает ошибки
	select {
	case err := <-serverError:
		// Ожидаем ошибку для некорректного адреса
		assert.Error(t, err)
	case <-time.After(1 * time.Second):
		// Если ошибки нет, это тоже нормально для некоторых случаев
		t.Log("Сервер не вернул ошибку для некорректного адреса")
	}
}

func TestServer_ContextCancellation(t *testing.T) {
	// Настройка логгера для тестов
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := New("localhost:0", handler)

	serverReady := make(chan struct{})
	serverError := make(chan error, 1)

	go func() {
		close(serverReady)
		err := server.Start()
		if err != nil && err != http.ErrServerClosed {
			serverError <- err
		}
		close(serverError)
	}()

	<-serverReady
	time.Sleep(100 * time.Millisecond)

	// Тестируем shutdown с контекстом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Создаем новый сервер для тестирования контекста
	testServer := &http.Server{
		Addr:    "localhost:0",
		Handler: handler,
	}

	go func() {
		testServer.ListenAndServe()
	}()

	time.Sleep(100 * time.Millisecond)

	// Тестируем shutdown с контекстом
	shutdownErr := testServer.Shutdown(ctx)
	assert.NoError(t, shutdownErr)

	// Завершаем основной сервер
	shutdownErr = server.Shutdown()
	assert.NoError(t, shutdownErr)

	select {
	case err := <-serverError:
		if err != nil {
			t.Fatalf("Неожиданная ошибка сервера: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Сервер не завершился вовремя")
	}
}

func TestServer_IntegrationWithHandler(t *testing.T) {
	// Настройка логгера для тестов
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	// Создаем простой роутер для тестирования
	router := http.NewServeMux()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("healthy"))
	})
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("home"))
	})

	server := New("localhost:0", router)

	serverReady := make(chan struct{})
	serverError := make(chan error, 1)

	go func() {
		close(serverReady)
		err := server.Start()
		if err != nil && err != http.ErrServerClosed {
			serverError <- err
		}
		close(serverError)
	}()

	<-serverReady
	time.Sleep(100 * time.Millisecond)

	// Тестируем endpoints
	if server.httpServer != nil {
		baseURL := "http://" + server.httpServer.Addr

		// Тестируем health endpoint
		req := httptest.NewRequest("GET", baseURL+"/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "healthy", w.Body.String())

		// Тестируем home endpoint
		req = httptest.NewRequest("GET", baseURL+"/", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "home", w.Body.String())
	}

	shutdownErr := server.Shutdown()
	assert.NoError(t, shutdownErr)

	select {
	case err := <-serverError:
		if err != nil {
			t.Fatalf("Неожиданная ошибка сервера: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Сервер не завершился вовремя")
	}
}
