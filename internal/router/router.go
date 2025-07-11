// Package router предоставляет функциональность для настройки маршрутизации HTTP-запросов.
// Использует библиотеку chi для создания маршрутизатора и настройки middleware.
package router

import (
	"gophkeeper/internal/compress"
	"net/http"

	"github.com/go-chi/chi"
)

// Возвращает настроенный маршрутизатор и ошибку, если она возникла.
func NewRouter() (chi.Router, error) {
	router := chi.NewRouter()
	router.Use(compress.WithGzipCompression)

	// API маршруты
	router.Route("/api", func(r chi.Router) {
		r.Post("/register", handleRegister)
		r.Post("/login", handleLogin)
		r.Post("/authenticate", handleAuthenticate)
		r.Post("/refresh", handleRefresh)

		r.Post("/passwords", handleCreatePassword)
		r.Post("/notes", handleCreateNote)
		r.Post("/cards", handleCreateCard)
		r.Post("/files", handleCreateFile)
		r.Get("/records/{label}", handleGetRecord)
		r.Put("/records/{label}", handleUpdateRecord)
		r.Get("/files/{label}/download", handleDownloadFile)

		r.Get("/health", handleHealth)
	})

	return router, nil
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

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
