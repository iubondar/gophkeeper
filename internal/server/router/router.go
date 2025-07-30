// Package router предоставляет функциональность для настройки маршрутизации HTTP-запросов.
// Использует библиотеку chi для создания маршрутизатора и настройки middleware.
package router

import (
	"gophkeeper/internal/server/api"
	"gophkeeper/internal/server/compress"
	"gophkeeper/internal/server/middleware"
	"gophkeeper/internal/server/storage/file"
	"gophkeeper/internal/server/storage/pg"
	"gophkeeper/internal/server/templates"
	"gophkeeper/internal/server/usecase"
	"net/http"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

// NewRouter создает и настраивает HTTP-маршрутизатор с использованием библиотеки chi.
// Функция настраивает все необходимые маршруты для API, включая аутентификацию,
// управление секретами, загрузку/скачивание файлов, а также Swagger UI.
//
// Параметры:
// - storage: хранилище данных PostgreSQL для работы с секретами и пользователями
// - fileStorage: хранилище файлов для загрузки и скачивания файлов
//
// Возвращает настроенный маршрутизатор и ошибку, если она возникла.
func NewRouter(storage *pg.Storage, fileStorage file.FileStorage) (chi.Router, error) {
	router := chi.NewRouter()
	router.Use(compress.WithGzipCompression)

	// Главная страница с информацией о сервисе
	router.Get("/", handleHomePage)

	healthHandler := api.NewHealthHandler(storage)
	router.Get("/health", healthHandler.Health)

	// Раздача swagger файлов (swagger.json, docs.go и т.д.)
	router.Handle("/swagger/*", http.StripPrefix("/swagger/", http.FileServer(http.Dir("./docs"))))

	// Swagger UI (использует swagger.json из /swagger/)
	router.Get("/swagger-ui/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/swagger.json"),
	))

	registerHandler := api.NewRegisterHandler(usecase.NewRegisterUsecase(storage))
	loginHandler := api.NewLoginHandler(usecase.NewLoginUsecase(storage))
	authenticateHandler := api.NewAuthenticateHandler(usecase.NewAuthenticateUsecase(storage))
	refreshHandler := api.NewRefreshHandler(usecase.NewRefreshUsecase())
	uploadHandler := api.NewUploadHandler(usecase.NewUploadSecretUsecase(storage))
	getHandler := api.NewGetHandler(usecase.NewGetSecretUsecase(storage))
	deleteHandler := api.NewDeleteHandler(usecase.NewDeleteSecretUsecase(storage, fileStorage))
	updateHandler := api.NewUpdateHandler(usecase.NewUpdateSecretUsecase(storage))
	versionHandler := api.NewVersionHandler(usecase.NewGetSecretUsecase(storage))

	// Создаем file usecases и handlers для загрузки и скачивания
	uploadFileUsecase := usecase.NewUploadFileUsecase(storage, fileStorage)
	downloadFileUsecase := usecase.NewDownloadFileUsecase(storage, fileStorage)

	uploadFileHandler := api.NewUploadFileHandler(uploadFileUsecase)
	downloadFileHandler := api.NewDownloadFileHandler(downloadFileUsecase)

	// API маршруты
	router.Route("/api", func(r chi.Router) {
		// Публичные маршруты (не требуют аутентификации)
		r.Post("/register", registerHandler.Register)
		r.Post("/login", loginHandler.Login)
		r.Post("/authenticate", authenticateHandler.Authenticate)
		r.Post("/refresh", refreshHandler.Refresh)

		// Защищенные маршруты (требуют аутентификации)
		r.Route("/", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)

			r.Post("/upload", uploadHandler.Upload)
			r.Get("/get", getHandler.GetSecret)
			r.Delete("/delete", deleteHandler.DeleteSecret)
			r.Put("/update", updateHandler.Update)
			r.Get("/version", versionHandler.GetSecretVersion)

			// File routes
			r.Post("/files", uploadFileHandler.UploadFile)
			r.Get("/files/{label}/download", downloadFileHandler.DownloadFile)
		})
	})

	return router, nil
}

// handleHomePage отображает главную страницу с информацией о сервисе
func handleHomePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(templates.HomePageHTML()))
}
