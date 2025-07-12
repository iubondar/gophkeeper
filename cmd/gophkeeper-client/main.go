package main

import (
	"gophkeeper/internal/client/api"
	"gophkeeper/internal/client/cmd"
	"gophkeeper/internal/client/shell"
	"gophkeeper/internal/config"
	"log"
	"os"
)

func main() {
	config, err := config.NewConfig(os.Args[0], os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	apiClient := api.NewAPIClient(config.RunAddress)
	commandFactory := cmd.NewCommandFactory(apiClient)

	// Создаем реестр команд
	commandRegistry := commandFactory.CreateCommandRegistry()
	commandExecutor := cmd.NewCommandExecutor(commandRegistry)

	// Создаем и запускаем интерактивный интерфейс
	shell := shell.NewShell(commandExecutor)
	if err := shell.Run(); err != nil {
		log.Fatal(err)
	}
}
