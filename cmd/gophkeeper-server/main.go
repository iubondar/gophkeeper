// Package main предоставляет точку входа для серверного приложения GophKeeper.
// Серверное приложение представляет собой REST API сервер для безопасного
// хранения и управления секретами и файлами пользователей.
//
// Основные возможности:
// - Регистрация и авторизация пользователей
// - Безопасное хранение секретов (текст, логины/пароли, карты)
// - Загрузка и скачивание файлов с шифрованием
// - Управление данными пользователей (создание, чтение, обновление, удаление)
// - Проверка состояния сервиса
// - Swagger документация API
//
// Сервер использует PostgreSQL для хранения метаданных и файловую систему
// для хранения зашифрованных файлов. Все данные шифруются на стороне клиента
// перед отправкой на сервер.
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

// Swagger документация для API GophKeeper
//
// @title GophKeeper API
// @version 1.0
// @description API для безопасного хранения секретов и файлов

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

// init инициализирует глобальный логгер zap для серверного приложения
func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
}

// main является точкой входа в серверное приложение GophKeeper.
// Функция выполняет следующие действия:
// 1. Инициализирует конфигурацию сервера
// 2. Подключается к базе данных PostgreSQL
// 3. Инициализирует файловое хранилище для файлов
// 4. Создает и настраивает HTTP роутер с API эндпоинтами
// 5. Запускает HTTP сервер на указанном адресе
//
// Сервер предоставляет REST API для:
// - Аутентификации и регистрации пользователей
// - Управления секретами (CRUD операции)
// - Загрузки и скачивания файлов
// - Проверки состояния сервиса
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
