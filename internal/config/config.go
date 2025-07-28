// Package config предоставляет функциональность для загрузки и управления конфигурацией приложения.
// Включает загрузку переменных окружения из .env файла и предоставляет структурированный доступ
// к настройкам базы данных, MinIO и сервера.
package config

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// Config представляет конфигурацию приложения.
// Содержит настройки для подключения к базе данных, MinIO и сервера.
type Config struct {
	RunAddress       string `env:"RUN_ADDRESS"`             // Адрес для запуска сервера
	DatabaseHost     string `env:"DATABASE_HOST"`           // Хост базы данных
	DatabaseUser     string `env:"DATABASE_USER"`           // Пользователь базы данных
	DatabasePassword string `env:"DATABASE_PASSWORD"`       // Пароль базы данных
	DatabaseName     string `env:"DATABASE_NAME"`           // Имя базы данных
	DatabaseSSLMode  string `env:"DATABASE_SSL_MODE"`       // Режим SSL для базы данных
	MinioEndpoint    string `env:"MINIO_ENDPOINT"`          // Эндпоинт MinIO
	MinioAccessKey   string `env:"MINIO_ACCESS_KEY_ID"`     // Ключ доступа MinIO
	MinioSecretKey   string `env:"MINIO_SECRET_ACCESS_KEY"` // Секретный ключ MinIO
	MinioUseSSL      bool   `env:"MINIO_USE_SSL"`           // Использовать SSL для MinIO
}

// NewConfig создает новый экземпляр конфигурации.
// Загружает переменные окружения из .env файла и парсит их в структуру Config.
//
// Параметры:
//   - progname: имя программы (не используется)
//   - args: аргументы командной строки (не используются)
//
// Возвращает:
//   - *Config: загруженная конфигурация
//   - error: ошибка в случае неудачи загрузки или парсинга
func NewConfig(progname string, args []string) (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	var c Config
	err = env.Parse(&c)
	if err != nil {
		return nil, err
	}

	zap.L().Sugar().Debugln(
		"Config: ",
		"RunAddress", c.RunAddress,
		"DatabaseUser", c.DatabaseUser,
		"DatabasePassword", c.DatabasePassword,
		"DatabaseHost", c.DatabaseHost,
		"DatabaseName", c.DatabaseName,
		"DatabaseSSLMode", c.DatabaseSSLMode,
		"MinioEndpoint", c.MinioEndpoint,
		"MinioAccessKeyID", c.MinioAccessKey,
		"MinioSecretAccessKey", c.MinioSecretKey,
		"MinioUseSSL", c.MinioUseSSL,
	)

	return &c, nil
}

// GetDatabaseURI возвращает строку подключения к базе данных PostgreSQL.
// Формирует URI в формате libpq для подключения к PostgreSQL.
//
// Возвращает:
//   - string: строка подключения к базе данных
func (c *Config) GetDatabaseURI() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DatabaseHost,
		c.DatabaseUser,
		c.DatabasePassword,
		c.DatabaseName,
		c.DatabaseSSLMode,
	)
}
