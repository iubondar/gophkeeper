// Package main предоставляет точку входа для клиентского приложения GophKeeper.
// Клиентское приложение представляет собой интерактивную оболочку для работы
// с сервером GophKeeper, позволяющую пользователям регистрироваться, авторизоваться,
// загружать, просматривать, обновлять и удалять секреты и файлы.
//
// Основные возможности:
// - Регистрация и авторизация пользователей
// - Загрузка и управление секретами (текст, логины/пароли, карты)
// - Загрузка и скачивание файлов
// - Просмотр и обновление сохраненных данных
// - Удаление данных
// - Проверка состояния сервера
//
// Приложение использует шифрование для защиты данных на стороне клиента
// и взаимодействует с сервером через REST API.
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

// Version содержит версию приложения, устанавливаемую при сборке
var Version string

// BuildTime содержит время сборки приложения
var BuildTime string

// init инициализирует глобальный логгер zap для приложения
func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
}

// main является точкой входа в клиентское приложение GophKeeper.
// Функция выполняет следующие действия:
// 1. Инициализирует конфигурацию приложения
// 2. Создает API клиент для взаимодействия с сервером
// 3. Инициализирует криптографический модуль
// 4. Регистрирует все доступные команды
// 5. Запускает интерактивную оболочку
func main() {
	config, err := config.NewConfig(os.Args[0], os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	apiClient := api.NewAPIClient(config.RunAddress)
	crypto := crypto.NewCrypto()
	commandRegistry := makeCommandRegistry(apiClient, crypto)

	// Создаем и запускаем интерактивный интерфейс
	shell := shell.NewShellWithDefaults(commandRegistry, Version, BuildTime)
	if err := shell.Run(); err != nil {
		log.Fatal(err)
	}
}

// makeCommandRegistry создает и настраивает реестр команд для клиентского приложения.
// Функция регистрирует все доступные команды:
// - register: регистрация нового пользователя
// - login: авторизация пользователя
// - upload: загрузка секретов и файлов
// - show: просмотр сохраненных данных
// - update: обновление существующих данных
// - get: получение данных с сервера
// - delete: удаление данных
// - health: проверка состояния сервера
// - version: отображение версии приложения
//
// Параметры:
//   - apiClient: клиент для взаимодействия с API сервера
//   - crypto: модуль для шифрования данных
//
// Возвращает настроенный реестр команд
func makeCommandRegistry(apiClient *api.APIClient, crypto *crypto.Crypto) *cmd.CommandRegistry {
	registry := cmd.NewCommandRegistry(crypto)
	registry.RegisterCommand(cmd.NewRegisterCommand(apiClient, crypto))
	registry.RegisterCommand(cmd.NewLoginCommand(apiClient, crypto))
	registry.RegisterCommand(cmd.NewUploadCommand(apiClient, crypto))
	registry.RegisterCommand(cmd.NewShowCommand(apiClient, crypto))
	registry.RegisterCommand(cmd.NewUpdateCommand(apiClient, crypto))
	registry.RegisterCommand(cmd.NewGetCommand(apiClient, crypto))
	registry.RegisterCommand(cmd.NewDeleteCommand(apiClient, crypto))
	registry.RegisterCommand(cmd.NewHealthCommand(apiClient))
	registry.RegisterCommand(cmd.NewVersionCommand(Version, BuildTime))

	return registry
}
