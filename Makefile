.PHONY: build-server build-cli run-server run-cli clean test

# Сборка сервера
build-server:
	go build -o bin/gophkeeper-server ./cmd/gophkeeper-server

# Сборка CLI клиента
build-client:
	go build -o bin/gophkeeper-client ./cmd/gophkeeper-client

# Сборка всех компонентов
build: build-server build-cli

# Запуск тестов
test:
	go test ./...

# Запуск сервера
run-server: build-server
	./bin/gophkeeper-server

# Запуск CLI клиента
run-cli: build-cli
	./bin/gophkeeper-client

# Очистка
clean:
	rm -rf bin/

# Создание директории bin
bin:
	mkdir -p bin 