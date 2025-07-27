.PHONY: build-server build-client run-server run-client clean test

# Получаем версию из git (если есть тег) или оставляем "unknown"
VERSION := $(shell git describe --tags 2>/dev/null || echo "unknown")
# Текущая дата и время в UTC
BUILD_TIME := $(shell date -u '+%Y-%m-%d %H:%M:%S')

# Сборка сервера
build-server:
	go build -o bin/gophkeeper-server ./cmd/gophkeeper-server

# Сборка CLI клиента
build-client:
	go build -ldflags="-X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)'" -o bin/gophkeeper-client ./cmd/gophkeeper-client

# Сборка всех компонентов
build: build-server build-client

# Запуск тестов
test:
	go test ./...

# Запуск сервера
run-server: build-server
	goose -dir ./internal/server/storage/pg/migrations postgres "user=ibondar password=postgres dbname=gophkeeper sslmode=disable" up
	./bin/gophkeeper-server

# Запуск CLI клиента
run-client: build-client
	./bin/gophkeeper-client

# Очистка
clean:
	rm -rf bin/

# Создание директории bin
bin:
	mkdir -p bin 