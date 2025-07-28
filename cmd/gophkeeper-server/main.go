package main

import (
	"log"
	"os"

	"gophkeeper/internal/config"
	"gophkeeper/internal/server"
	"gophkeeper/internal/server/router"
	"gophkeeper/internal/server/storage/file"
	"gophkeeper/internal/server/storage/pg"

	"go.uber.org/zap"
)

// @title GophKeeper API
// @version 1.0
// @description API для безопасного хранения секретов и файлов
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Токен доступа в формате Bearer {token}

// @tag.name auth
// @tag.description Операции аутентификации и регистрации

// @tag.name secrets
// @tag.description Операции с секретами (текст, логины/пароли, карты)

// @tag.name files
// @tag.description Операции с файлами

// @tag.name health
// @tag.description Проверка состояния сервиса

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
}

func main() {
	config, err := config.NewConfig(os.Args[0], os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	storage, err := pg.NewStorage(config.GetDatabaseURI())
	if err != nil {
		log.Fatal(err)
	}

	// Создаем file storage
	fileStorage, err := file.NewStorage(config)
	if err != nil {
		log.Fatal(err)
	}

	router, err := router.NewRouter(storage, fileStorage)
	if err != nil {
		log.Fatal(err)
	}

	// Создаем и запускаем сервер
	srv := server.New(config.RunAddress, router)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
