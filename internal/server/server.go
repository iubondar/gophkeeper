package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
)

type Server struct {
	address    string
	router     http.Handler
	httpServer *http.Server
}

func New(addr string, handler http.Handler) *Server {
	return &Server{
		address: addr,
		router:  handler,
	}
}

// Start запускает сервер и начинает обработку запросов
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

// Shutdown выполняет graceful shutdown сервера
func (s *Server) Shutdown() error {
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
