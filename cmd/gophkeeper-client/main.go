package main

import (
	"gophkeeper/internal/client/api"
	"gophkeeper/internal/client/cmd"
	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/client/shell"
	"gophkeeper/internal/config"
	"log"
	"os"

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

	apiClient := api.NewAPIClient(config.RunAddress)
	crypto := crypto.NewCrypto()
	commandRegistry := cmd.NewCommandRegistry(apiClient, crypto)

	// Создаем и запускаем интерактивный интерфейс
	shell := shell.NewShell(commandRegistry)
	if err := shell.Run(); err != nil {
		log.Fatal(err)
	}
}
