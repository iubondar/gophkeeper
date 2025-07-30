package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func init() {
	// Initialize logger for tests
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}

func TestNewConfig_ValidEnvironment(t *testing.T) {

	tmpDir, err := os.MkdirTemp("", "test_config")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create .env file in the temporary directory
	envContent := `RUN_ADDRESS=:8080
DATABASE_USER=testuser
DATABASE_PASSWORD=testpass
DATABASE_HOST=localhost
DATABASE_NAME=testdb
DATABASE_SSL_MODE=disable
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY_ID=minioadmin
MINIO_SECRET_ACCESS_KEY=minioadmin
MINIO_USE_SSL=false`

	envFile := tmpDir + "/.env"
	err = os.WriteFile(envFile, []byte(envContent), 0644)
	require.NoError(t, err)

	// Temporarily change working directory to where the .env file is
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	os.Chdir(tmpDir)

	// Test NewConfig
	config, err := NewConfig("test", []string{})
	require.NoError(t, err)
	assert.NotNil(t, config)

	// Verify configuration values
	assert.Equal(t, ":8080", config.RunAddress)
	assert.Equal(t, "testuser", config.DatabaseUser)
	assert.Equal(t, "testpass", config.DatabasePassword)
	assert.Equal(t, "localhost", config.DatabaseHost)
	assert.Equal(t, "testdb", config.DatabaseName)
	assert.Equal(t, "disable", config.DatabaseSSLMode)
	assert.Equal(t, "localhost:9000", config.MinioEndpoint)
	assert.Equal(t, "minioadmin", config.MinioAccessKey)
	assert.Equal(t, "minioadmin", config.MinioSecretKey)
	assert.Equal(t, false, config.MinioUseSSL)
}

func TestNewConfig_MissingEnvFile(t *testing.T) {
	// Change to a directory without .env file
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	tmpDir, err := os.MkdirTemp("", "test_no_env")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	// Test NewConfig without .env file
	config, err := NewConfig("test", []string{})
	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestNewConfig_InvalidEnvFile(t *testing.T) {
	// Create temporary directory for .env file
	tmpDir, err := os.MkdirTemp("", "test_invalid_config")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create invalid .env file
	envContent := `INVALID_ENV_FORMAT
DATABASE_USER=testuser
DATABASE_PASSWORD=testpass
DATABASE_HOST=localhost
DATABASE_NAME=testdb
DATABASE_SSL_MODE=disable
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY_ID=minioadmin
MINIO_SECRET_ACCESS_KEY=minioadmin
MINIO_USE_SSL=false`

	envFile := tmpDir + "/.env"
	err = os.WriteFile(envFile, []byte(envContent), 0644)
	require.NoError(t, err)

	// Temporarily change working directory
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	os.Chdir(tmpDir)

	// Test NewConfig with invalid .env file
	config, err := NewConfig("test", []string{})
	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestGetDatabaseURI(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected string
	}{
		{
			name: "complete database config",
			config: &Config{
				DatabaseHost:     "localhost",
				DatabaseUser:     "testuser",
				DatabasePassword: "testpass",
				DatabaseName:     "testdb",
				DatabaseSSLMode:  "disable",
			},
			expected: "host=localhost user=testuser password=testpass dbname=testdb sslmode=disable",
		},
		{
			name: "database config with empty values",
			config: &Config{
				DatabaseHost:     "",
				DatabaseUser:     "",
				DatabasePassword: "",
				DatabaseName:     "",
				DatabaseSSLMode:  "",
			},
			expected: "host= user= password= dbname= sslmode=",
		},
		{
			name: "database config with special characters",
			config: &Config{
				DatabaseHost:     "localhost:5432",
				DatabaseUser:     "user@domain",
				DatabasePassword: "pass@word!",
				DatabaseName:     "test-db",
				DatabaseSSLMode:  "require",
			},
			expected: "host=localhost:5432 user=user@domain password=pass@word! dbname=test-db sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetDatabaseURI()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfig_StructFields(t *testing.T) {
	// Test that all struct fields are properly defined
	config := &Config{}

	// Test database fields
	config.DatabaseUser = "test-user"
	config.DatabasePassword = "test-password"
	config.DatabaseHost = "test-host"
	config.DatabaseName = "test-dbname"
	config.DatabaseSSLMode = "test-sslmode"

	assert.Equal(t, "test-user", config.DatabaseUser)
	assert.Equal(t, "test-password", config.DatabasePassword)
	assert.Equal(t, "test-host", config.DatabaseHost)
	assert.Equal(t, "test-dbname", config.DatabaseName)
	assert.Equal(t, "test-sslmode", config.DatabaseSSLMode)

	// Test MinIO fields
	config.MinioEndpoint = "test-endpoint"
	config.MinioAccessKey = "test-access-key"
	config.MinioSecretKey = "test-secret-key"
	config.MinioUseSSL = true

	assert.Equal(t, "test-endpoint", config.MinioEndpoint)
	assert.Equal(t, "test-access-key", config.MinioAccessKey)
	assert.Equal(t, "test-secret-key", config.MinioSecretKey)
	assert.Equal(t, true, config.MinioUseSSL)
}

func TestNewConfig_WithCommandLineArgs(t *testing.T) {
	// Create temporary directory for .env file
	tmpDir, err := os.MkdirTemp("", "test_args_config")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create .env file in the temporary directory
	envContent := `RUN_ADDRESS=:8080
DATABASE_USER=testuser
DATABASE_PASSWORD=testpass
DATABASE_HOST=localhost
DATABASE_NAME=testdb
DATABASE_SSL_MODE=disable
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY_ID=minioadmin
MINIO_SECRET_ACCESS_KEY=minioadmin
MINIO_USE_SSL=false`

	envFile := tmpDir + "/.env"
	err = os.WriteFile(envFile, []byte(envContent), 0644)
	require.NoError(t, err)

	// Temporarily change working directory
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	os.Chdir(tmpDir)

	// Test NewConfig with command line arguments
	args := []string{"--help", "--version"}
	config, err := NewConfig("test", args)
	require.NoError(t, err)
	assert.NotNil(t, config)

	// Verify that configuration is still loaded correctly
	assert.Equal(t, ":8080", config.RunAddress)
	assert.Equal(t, "testuser", config.DatabaseUser)
	assert.Equal(t, "localhost:9000", config.MinioEndpoint)
}
