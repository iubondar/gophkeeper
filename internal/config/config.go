package config

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type MinioConfig struct {
	Endpoint        string `env:"MINIO_ENDPOINT"`
	AccessKeyID     string `env:"MINIO_ACCESS_KEY"`
	SecretAccessKey string `env:"MINIO_SECRET_KEY"`
	UseSSL          bool   `env:"MINIO_USE_SSL"`
}

type DatabaseConfig struct {
	URI      string `env:"DATABASE_URI"`
	User     string `env:"DATABASE_USER"`
	Password string `env:"DATABASE_PASSWORD"`
	Host     string `env:"DATABASE_HOST"`
	DBName   string `env:"DATABASE_NAME"`
	SSLMode  string `env:"DATABASE_SSL_MODE"`
}

type Config struct {
	RunAddress     string
	DatabaseConfig DatabaseConfig
	MinioConfig    MinioConfig
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
		"DatabaseURI", c.DatabaseConfig.URI,
		"DatabaseUser", c.DatabaseConfig.User,
		"DatabasePassword", c.DatabaseConfig.Password,
		"DatabaseHost", c.DatabaseConfig.Host,
		"DatabaseName", c.DatabaseConfig.DBName,
		"DatabaseSSLMode", c.DatabaseConfig.SSLMode,
		"MinioEndpoint", c.MinioConfig.Endpoint,
		"MinioAccessKeyID", c.MinioConfig.AccessKeyID,
		"MinioSecretAccessKey", c.MinioConfig.SecretAccessKey,
		"MinioUseSSL", c.MinioConfig.UseSSL,
	)

	return &c, nil
}

func (c *Config) GetDatabaseURI() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DatabaseConfig.Host,
		c.DatabaseConfig.User,
		c.DatabaseConfig.Password,
		c.DatabaseConfig.DBName,
		c.DatabaseConfig.SSLMode,
	)
}
