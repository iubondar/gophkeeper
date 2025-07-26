package config

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	RunAddress       string `env:"RUN_ADDRESS"`
	DatabaseURI      string `env:"DATABASE_URI"`
	DatabaseUser     string `env:"DATABASE_USER"`
	DatabasePassword string `env:"DATABASE_PASSWORD"`
	DatabaseHost     string `env:"DATABASE_HOST"`
	DatabaseName     string `env:"DATABASE_NAME"`
	DatabaseSSLMode  string `env:"DATABASE_SSL_MODE"`
	MinioEndpoint    string `env:"MINIO_ENDPOINT"`
	MinioAccessKey   string `env:"MINIO_ACCESS_KEY"`
	MinioSecretKey   string `env:"MINIO_SECRET_KEY"`
	MinioUseSSL      bool   `env:"MINIO_USE_SSL"`
}

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
		"DatabaseURI", c.DatabaseURI,
		"DatabaseUser", c.DatabaseUser,
		"DatabasePassword", c.DatabasePassword,
		"DatabaseHost", c.DatabaseHost,
		"DatabaseName", c.DatabaseName,
		"DatabaseSSLMode", c.DatabaseSSLMode,
		"MinioEndpoint", c.MinioEndpoint,
		"MinioAccessKey", c.MinioAccessKey,
		"MinioSecretKey", c.MinioSecretKey,
		"MinioUseSSL", c.MinioUseSSL,
	)

	return &c, nil
}

func (c *Config) GetDatabaseURI() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DatabaseHost,
		c.DatabaseUser,
		c.DatabasePassword,
		c.DatabaseName,
		c.DatabaseSSLMode,
	)
}
