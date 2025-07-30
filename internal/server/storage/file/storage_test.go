package file

import (
	"context"
	"io"
	"strings"
	"testing"

	"gophkeeper/internal/config"
	"gophkeeper/internal/server/storage/testhelpers"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type FileStorageTestSuite struct {
	suite.Suite
	storage *Storage
	cleanup func()
}

func (s *FileStorageTestSuite) SetupSuite() {
	ctx := context.Background()
	container, err := testhelpers.CreateMinioContainer(ctx)
	require.NoError(s.T(), err)

	config := &config.Config{
		MinioEndpoint:  container.Endpoint,
		MinioAccessKey: container.AccessKey,
		MinioSecretKey: container.SecretKey,
		MinioUseSSL:    false,
	}

	storage, err := NewStorage(config)
	require.NoError(s.T(), err)

	s.storage = storage
	s.cleanup = func() {
		container.Terminate(ctx)
	}
}

func (s *FileStorageTestSuite) TearDownSuite() {
	if s.cleanup != nil {
		s.cleanup()
	}
}

func TestFileStorageSuite(t *testing.T) {
	suite.Run(t, new(FileStorageTestSuite))
}

func (s *FileStorageTestSuite) TestNewStorage() {
	ctx := context.Background()
	container, err := testhelpers.CreateMinioContainer(ctx)
	require.NoError(s.T(), err)
	defer container.Terminate(ctx)

	config := &config.Config{
		MinioEndpoint:  container.Endpoint,
		MinioAccessKey: container.AccessKey,
		MinioSecretKey: container.SecretKey,
		MinioUseSSL:    false,
	}

	storage, err := NewStorage(config)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), storage)
	require.NotNil(s.T(), storage.minioClient)
}

func (s *FileStorageTestSuite) TestNewStorageWithError() {
	// Тест на ошибку при создании MinIO клиента
	config := &config.Config{
		MinioEndpoint:  "invalid-endpoint:9999",
		MinioAccessKey: "invalid-key",
		MinioSecretKey: "invalid-secret",
		MinioUseSSL:    false,
	}

	storage, err := NewStorage(config)
	require.Error(s.T(), err)
	require.Nil(s.T(), storage)
}

func (s *FileStorageTestSuite) TestUploadFile() {
	ctx := context.Background()
	fileKey := "test-file-upload"
	content := "test content for upload"
	reader := strings.NewReader(content)

	s.Run("successful upload", func() {
		err := s.storage.UploadFile(ctx, fileKey, reader, int64(len(content)))
		s.Require().NoError(err)

		// Проверяем, что файл существует
		exists, err := s.storage.FileExists(ctx, fileKey)
		s.Require().NoError(err)
		s.Require().True(exists)
	})

	s.Run("upload with nil client", func() {
		storage := &Storage{minioClient: nil}
		err := storage.UploadFile(ctx, "test-nil", reader, int64(len(content)))
		s.Require().Error(err)
		s.Require().Contains(err.Error(), "minio client is not initialized")
	})

	s.Run("upload with invalid data", func() {
		// Тест на ошибку загрузки с неверными данными
		// Создаем reader, который будет возвращать ошибку
		errorReader := &errorReader{}
		err := s.storage.UploadFile(ctx, "test-error", errorReader, 100)
		s.Require().Error(err)
	})
}

// errorReader - reader, который всегда возвращает ошибку
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func (s *FileStorageTestSuite) TestDownloadFile() {
	ctx := context.Background()
	fileKey := "test-file-download"
	content := "test content for download"
	reader := strings.NewReader(content)

	// Сначала загружаем файл
	err := s.storage.UploadFile(ctx, fileKey, reader, int64(len(content)))
	s.Require().NoError(err)

	s.Run("successful download", func() {
		downloadedReader, err := s.storage.DownloadFile(ctx, fileKey)
		s.Require().NoError(err)
		defer downloadedReader.Close()

		downloadedContent, err := io.ReadAll(downloadedReader)
		s.Require().NoError(err)
		s.Require().Equal(content, string(downloadedContent))
	})

	s.Run("download non-existent file", func() {
		_, err := s.storage.DownloadFile(ctx, "non-existent-file")
		s.Require().Error(err)
	})

	s.Run("download with nil client", func() {
		storage := &Storage{minioClient: nil}
		_, err := storage.DownloadFile(ctx, fileKey)
		s.Require().Error(err)
		s.Require().Contains(err.Error(), "minio client is not initialized")
	})

	s.Run("download with connection error", func() {
		// Создаем storage с неверным endpoint для тестирования ошибки соединения
		invalidConfig := &config.Config{
			MinioEndpoint:  "invalid-endpoint:9999",
			MinioAccessKey: "invalid-key",
			MinioSecretKey: "invalid-secret",
			MinioUseSSL:    false,
		}
		invalidStorage, err := NewStorage(invalidConfig)
		if err == nil {
			// Если удалось создать storage, тестируем скачивание
			_, err := invalidStorage.DownloadFile(ctx, "test-file")
			s.Require().Error(err)
		}
	})
}

func (s *FileStorageTestSuite) TestDeleteFile() {
	ctx := context.Background()
	fileKey := "test-file-delete"
	content := "test content for delete"
	reader := strings.NewReader(content)

	// Сначала загружаем файл
	err := s.storage.UploadFile(ctx, fileKey, reader, int64(len(content)))
	s.Require().NoError(err)

	s.Run("successful delete", func() {
		err := s.storage.DeleteFile(ctx, fileKey)
		s.Require().NoError(err)

		// Проверяем, что файл удален
		exists, err := s.storage.FileExists(ctx, fileKey)
		s.Require().NoError(err)
		s.Require().False(exists)
	})

	s.Run("delete non-existent file", func() {
		err := s.storage.DeleteFile(ctx, "non-existent-file")
		s.Require().NoError(err) // MinIO не возвращает ошибку при удалении несуществующего файла
	})

	s.Run("delete with nil client", func() {
		storage := &Storage{minioClient: nil}
		err := storage.DeleteFile(ctx, fileKey)
		s.Require().Error(err)
		s.Require().Contains(err.Error(), "minio client is not initialized")
	})

	s.Run("delete with connection error", func() {
		// Создаем storage с неверным endpoint для тестирования ошибки соединения
		invalidConfig := &config.Config{
			MinioEndpoint:  "invalid-endpoint:9999",
			MinioAccessKey: "invalid-key",
			MinioSecretKey: "invalid-secret",
			MinioUseSSL:    false,
		}
		invalidStorage, err := NewStorage(invalidConfig)
		if err == nil {
			// Если удалось создать storage, тестируем удаление
			err := invalidStorage.DeleteFile(ctx, "test-file")
			s.Require().Error(err)
		}
	})
}

func (s *FileStorageTestSuite) TestFileExists() {
	ctx := context.Background()
	fileKey := "test-file-exists"
	content := "test content for exists"
	reader := strings.NewReader(content)

	s.Run("file exists", func() {
		// Сначала загружаем файл
		err := s.storage.UploadFile(ctx, fileKey, reader, int64(len(content)))
		s.Require().NoError(err)

		exists, err := s.storage.FileExists(ctx, fileKey)
		s.Require().NoError(err)
		s.Require().True(exists)
	})

	s.Run("file does not exist", func() {
		exists, err := s.storage.FileExists(ctx, "non-existent-file")
		s.Require().NoError(err)
		s.Require().False(exists)
	})

	s.Run("exists with nil client", func() {
		storage := &Storage{minioClient: nil}
		exists, err := storage.FileExists(ctx, fileKey)
		s.Require().Error(err)
		s.Require().False(exists)
		s.Require().Contains(err.Error(), "minio client is not initialized")
	})

	s.Run("exists with connection error", func() {
		// Создаем storage с неверным endpoint для тестирования ошибки соединения
		invalidConfig := &config.Config{
			MinioEndpoint:  "invalid-endpoint:9999",
			MinioAccessKey: "invalid-key",
			MinioSecretKey: "invalid-secret",
			MinioUseSSL:    false,
		}
		invalidStorage, err := NewStorage(invalidConfig)
		if err == nil {
			// Если удалось создать storage, тестируем проверку существования
			exists, err := invalidStorage.FileExists(ctx, "test-file")
			s.Require().Error(err)
			s.Require().False(exists)
		}
	})
}

func (s *FileStorageTestSuite) TestEnsureBucketExists() {
	ctx := context.Background()
	container, err := testhelpers.CreateMinioContainer(ctx)
	require.NoError(s.T(), err)
	defer container.Terminate(ctx)

	config := &config.Config{
		MinioEndpoint:  container.Endpoint,
		MinioAccessKey: container.AccessKey,
		MinioSecretKey: container.SecretKey,
		MinioUseSSL:    false,
	}

	storage, err := NewStorage(config)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), storage)

	// Проверяем, что бакет создан
	exists, err := storage.minioClient.BucketExists(BucketName)
	require.NoError(s.T(), err)
	require.True(s.T(), exists)
}
