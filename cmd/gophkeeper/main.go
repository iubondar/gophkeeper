package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// Создаем роутер
	r := mux.NewRouter()

	// Настраиваем маршруты
	setupRoutes(r)

	// Создаем HTTP сервер
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запускаем сервер в горутине
	go func() {
		log.Printf("Сервер запущен на порту 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	// Ждем сигнала для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Завершение работы сервера...")

	// Даем серверу время на завершение
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при завершении сервера: %v", err)
	}

	log.Println("Сервер успешно завершен")
}

func setupRoutes(r *mux.Router) {
	// API маршруты
	api := r.PathPrefix("/api").Subrouter()

	// Аутентификация
	api.HandleFunc("/register", handleRegister).Methods("POST")
	api.HandleFunc("/login", handleLogin).Methods("POST")
	api.HandleFunc("/authenticate", handleAuthenticate).Methods("POST")
	api.HandleFunc("/refresh", handleRefresh).Methods("POST")

	// Записи
	api.HandleFunc("/passwords", handleCreatePassword).Methods("POST")
	api.HandleFunc("/notes", handleCreateNote).Methods("POST")
	api.HandleFunc("/cards", handleCreateCard).Methods("POST")
	api.HandleFunc("/files", handleCreateFile).Methods("POST")
	api.HandleFunc("/records/{label}", handleGetRecord).Methods("GET")
	api.HandleFunc("/records/{label}", handleUpdateRecord).Methods("PUT")
	api.HandleFunc("/files/{label}/download", handleDownloadFile).Methods("GET")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
}

// Заглушки для обработчиков
func handleRegister(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Регистрация пока не реализована"}`))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"salt": "dGVzdC1zYWx0"}`))
}

func handleAuthenticate(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"access_token": "test-token", "refresh_token": "test-refresh", "expires_in": 1800}`))
}

func handleRefresh(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"access_token": "new-test-token", "refresh_token": "new-test-refresh", "expires_in": 1800}`))
}

func handleCreatePassword(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"id": "test-id", "version": 1}`))
}

func handleCreateNote(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"id": "test-id", "version": 1}`))
}

func handleCreateCard(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"id": "test-id", "version": 1}`))
}

func handleCreateFile(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"id": "test-id", "version": 1}`))
}

func handleGetRecord(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"id": "test-id", "type": "password", "label": "test", "metadata": "test", "data": "test", "version": 1}`))
}

func handleUpdateRecord(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"version": 2}`))
}

func handleDownloadFile(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("test file content"))
}
