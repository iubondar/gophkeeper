.PHONY: build-server build-cli run-server run-cli clean test

# Сборка сервера
build-server:
	go build -o bin/gophkeeper ./cmd/gophkeeper

# Сборка CLI клиента
build-cli:
	go build -o bin/cli ./cmd/cli

# Сборка всех компонентов
build: build-server build-cli

# Запуск тестов
test:
	go test ./...

# Запуск сервера
run-server: build-server
	./bin/gophkeeper

# Запуск CLI клиента
run-cli: build-cli
	./bin/cli

# Очистка
clean:
	rm -rf bin/

# Создание директории bin
bin:
	mkdir -p bin 