// Package server предоставляет функциональность для создания и управления HTTP-сервером.
// Пакет поддерживает как локальную разработку с использованием статических сертификатов,
// так и продакшен-среду с автоматическим получением SSL-сертификатов через Let's Encrypt.
package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
)

// Server представляет HTTP-сервер с поддержкой graceful shutdown и автоматического получения SSL-сертификатов.
// Сервер может работать как в режиме разработки (localhost с HTTP/HTTPS), так и в продакшене
// с автоматическим получением SSL-сертификатов через Let's Encrypt.
type Server struct {
	address    string
	router     http.Handler
	httpServer *http.Server
}

// New создает новый экземпляр Server с указанным адресом и HTTP-обработчиком.
// Адрес должен быть в формате "host:port" (например, "localhost:8080" или "example.com:443").
// Для localhost сервер будет использовать статические сертификаты, для других адресов - autocert.
func New(addr string, handler http.Handler) *Server {
	return &Server{
		address: addr,
		router:  handler,
	}
}

// Start запускает сервер и начинает обработку запросов.
// Метод блокирует выполнение до получения сигнала завершения или ошибки сервера.
//
// Поведение сервера зависит от адреса:
// - Для localhost: использует HTTP или HTTPS с статическими сертификатами (если доступны)
// - Для других адресов: использует HTTPS с автоматическим получением сертификатов через Let's Encrypt
//
// Сервер корректно обрабатывает сигналы SIGINT и SIGTERM для graceful shutdown.
// Возвращает ошибку в случае проблем с запуском или завершением сервера.
func (s *Server) Start() error {
	// Канал для обработки ошибок сервера
	serverErrors := make(chan error, 1)

	// Извлекаем хост из адреса (без порта)
	isLocalhost := strings.Contains(s.address, "localhost")

	if isLocalhost {
		// Для localhost используем статические сертификаты
		s.httpServer = &http.Server{
			Addr:    s.address,
			Handler: s.router,
		}
	} else {
		// Для продакшена используем autocert
		m := &autocert.Manager{
			Cache:      autocert.DirCache("certs"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(s.address),
		}
		s.httpServer = &http.Server{
			Addr:      s.address,
			TLSConfig: m.TLSConfig(),
			Handler:   s.router,
		}
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		if isLocalhost {
			zap.L().Info("Starting server", zap.String("address", s.httpServer.Addr))
			// Проверяем наличие сертификатов
			if _, err := os.Stat("certs/cert.pem"); os.IsNotExist(err) {
				// Если сертификатов нет, используем HTTP
				zap.L().Info("No certificates found, using HTTP")
				serverErrors <- s.httpServer.ListenAndServe()
			} else {
				// Если сертификаты есть, используем HTTPS
				serverErrors <- s.httpServer.ListenAndServeTLS("certs/cert.pem", "certs/key.pem")
			}
		} else {
			zap.L().Info("Starting server", zap.String("address", s.httpServer.Addr))
			serverErrors <- s.httpServer.ListenAndServeTLS("", "")
		}
	}()

	// Канал для обработки сигналов завершения от ОС
	shutdown := make(chan os.Signal, 1)
	// Регистрируем обработчики для SIGINT (Ctrl+C) и SIGTERM
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Ожидаем либо ошибку сервера, либо сигнал завершения
	select {
	case err := <-serverErrors:
		zap.L().Error("Server error", zap.Error(err))
		return err

	case sig := <-shutdown:
		zap.L().Info("Start shutdown", zap.String("signal", sig.String()))
		return s.Shutdown()
	}
}

// Shutdown выполняет graceful shutdown сервера.
// Метод пытается корректно завершить все текущие соединения в течение 5 секунд.
// Если graceful shutdown не удается, сервер принудительно закрывается.
//
// Возвращает ошибку, если сервер не был запущен или возникли проблемы при завершении.
func (s *Server) Shutdown() error {
	// Проверяем, что сервер был инициализирован
	if s.httpServer == nil {
		return fmt.Errorf("server was not started")
	}

	// Устанавливаем таймаут 5 секунд для завершения текущих запросов
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Пытаемся корректно завершить работу сервера
	if err := s.httpServer.Shutdown(ctx); err != nil {
		zap.L().Error("Graceful shutdown did not complete", zap.Error(err))
		// Если плавное завершение не удалось, принудительно закрываем сервер
		if err := s.httpServer.Close(); err != nil {
			zap.L().Error("Could not stop server", zap.Error(err))
			return err
		}
		return err
	}

	return nil
}
