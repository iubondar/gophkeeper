// Package testhelpers предоставляет вспомогательные функции для тестирования хранилищ.
// Пакет содержит функции для создания тестовых контейнеров PostgreSQL и MinIO
// с использованием testcontainers для интеграционного тестирования.
package testhelpers

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// PostgresContainer представляет тестовый контейнер PostgreSQL.
// Структура содержит контейнер PostgreSQL и строку подключения к нему
// для использования в интеграционных тестах.
type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

// CreatePostgresContainer создает тестовый контейнер PostgreSQL.
// Функция запускает контейнер PostgreSQL с предустановленными настройками
// и возвращает структуру с контейнером и строкой подключения.
func CreatePostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	pgContainer, err := postgres.Run(ctx,
		"postgres:15.3-alpine",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connStr,
	}, nil
}
