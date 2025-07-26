// Package router предоставляет функциональность для настройки маршрутизации HTTP-запросов.
// Использует библиотеку chi для создания маршрутизатора и настройки middleware.
package router

import (
	"gophkeeper/internal/server/api"
	"gophkeeper/internal/server/compress"
	"gophkeeper/internal/server/storage/file"
	"gophkeeper/internal/server/storage/pg"
	"gophkeeper/internal/server/usecase"
	"net/http"

	"github.com/go-chi/chi"
)

// Возвращает настроенный маршрутизатор и ошибку, если она возникла.
func NewRouter(storage *pg.Storage, fileStorage file.FileStorage) (chi.Router, error) {
	router := chi.NewRouter()
	router.Use(compress.WithGzipCompression)

	healthHandler := api.NewHealthHandler(storage)
	router.Get("/health", healthHandler.Health)

	registerHandler := api.NewRegisterHandler(usecase.NewRegisterUsecase(storage))
	loginHandler := api.NewLoginHandler(usecase.NewLoginUsecase(storage))
	authenticateHandler := api.NewAuthenticateHandler(usecase.NewAuthenticateUsecase(storage))
	uploadHandler := api.NewUploadHandler(usecase.NewUploadSecretUsecase(storage))
	getHandler := api.NewGetHandler(usecase.NewGetSecretUsecase(storage))
	deleteHandler := api.NewDeleteHandler(usecase.NewDeleteSecretUsecase(storage))
	updateHandler := api.NewUpdateHandler(usecase.NewUpdateSecretUsecase(storage))
	versionHandler := api.NewVersionHandler(usecase.NewGetSecretUsecase(storage))

	// Создаем file usecases и handlers
	uploadFileUsecase := usecase.NewUploadFileUsecase(storage, fileStorage)
	downloadFileUsecase := usecase.NewDownloadFileUsecase(storage, fileStorage)
	deleteFileUsecase := usecase.NewDeleteFileUsecase(storage, fileStorage)

	uploadFileHandler := api.NewUploadFileHandler(uploadFileUsecase)
	downloadFileHandler := api.NewDownloadFileHandler(downloadFileUsecase)
	deleteFileHandler := api.NewDeleteFileHandler(deleteFileUsecase)

	// API маршруты
	router.Route("/api", func(r chi.Router) {
		r.Post("/register", registerHandler.Register)
		r.Post("/login", loginHandler.Login)
		r.Post("/authenticate", authenticateHandler.Authenticate)
		r.Post("/refresh", handleRefresh)
		r.Post("/upload", uploadHandler.Upload)
		r.Get("/get", getHandler.GetSecret)
		r.Delete("/delete", deleteHandler.DeleteSecret)
		r.Put("/update", updateHandler.Update)
		r.Get("/version", versionHandler.GetSecretVersion)

		// File routes
		r.Post("/files", uploadFileHandler.UploadFile)
		r.Get("/files/{label}/download", downloadFileHandler.DownloadFile)
		r.Delete("/files/{label}", deleteFileHandler.DeleteFile)
	})

	return router, nil
}

func handleRefresh(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"access_token": "new-test-token", "refresh_token": "new-test-refresh", "expires_in": 1800}`))
}
