package config

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	RunAddress       string `env:"RUN_ADDRESS"`
	DatabaseHost     string `env:"DATABASE_HOST"`
	DatabaseUser     string `env:"DATABASE_USER"`
	DatabasePassword string `env:"DATABASE_PASSWORD"`
	DatabaseName     string `env:"DATABASE_NAME"`
	DatabaseSSLMode  string `env:"DATABASE_SSL_MODE"`
	MinioEndpoint    string `env:"MINIO_ENDPOINT"`
	MinioAccessKey   string `env:"MINIO_ACCESS_KEY_ID"`
	MinioSecretKey   string `env:"MINIO_SECRET_ACCESS_KEY"`
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

func (c *Config) GetDatabaseURI() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DatabaseHost,
		c.DatabaseUser,
		c.DatabasePassword,
		c.DatabaseName,
		c.DatabaseSSLMode,
	)
}
