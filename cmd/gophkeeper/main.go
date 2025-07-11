package main

import (
	"log"
	"os"

	"gophkeeper/internal/config"
	"gophkeeper/internal/router"
	"gophkeeper/internal/server"

	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
}

func main() {
	config, err := config.NewConfig(os.Args[0], os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	router, err := router.NewRouter()
	if err != nil {
		log.Fatal(err)
	}

	// Создаем и запускаем сервер
	srv := server.New(config.RunAddress, router)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
