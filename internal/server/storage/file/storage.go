// Package file предоставляет реализацию файлового хранилища на основе MinIO.
// Пакет содержит интерфейс FileStorage и его реализацию для загрузки,
// скачивания и удаления файлов в объектном хранилище MinIO.
package file

import (
	"context"
	"fmt"
	"gophkeeper/internal/config"
	"io"

	"github.com/minio/minio-go"
	"go.uber.org/zap"
)

const (
	// BucketName - имя бакета для хранения файлов
	BucketName = "gophkeeper-files"
)

// FileStorage определяет интерфейс для работы с файловым хранилищем.
// Интерфейс предоставляет методы для загрузки, скачивания, удаления файлов
// и проверки их существования в объектном хранилище.
type FileStorage interface {
	UploadFile(ctx context.Context, fileKey string, reader io.Reader, size int64) error
	DownloadFile(ctx context.Context, fileKey string) (io.ReadCloser, error)
	DeleteFile(ctx context.Context, fileKey string) error
	FileExists(ctx context.Context, fileKey string) (bool, error)
}

// Storage представляет реализацию файлового хранилища на основе MinIO.
// Структура содержит клиент MinIO для взаимодействия с объектным хранилищем.
type Storage struct {
	minioClient *minio.Client
}

// NewStorage создает новый экземпляр Storage с настройками MinIO.
// Принимает конфигурацию для подключения к MinIO серверу.
// Функция инициализирует клиент MinIO и создает бакет, если он не существует.
func NewStorage(config *config.Config) (*Storage, error) {
	minioClient, err := minio.New(
		config.MinioEndpoint,
		config.MinioAccessKey,
		config.MinioSecretKey,
		config.MinioUseSSL,
	)
	if err != nil {
		return nil, err
	}

	storage := &Storage{minioClient: minioClient}

	// Создаем бакет при инициализации, если он не существует
	if err := storage.ensureBucketExists(); err != nil {
		return nil, err
	}

	return storage, nil
}

// ensureBucketExists создает бакет, если он не существует
func (s *Storage) ensureBucketExists() error {
	exists, err := s.minioClient.BucketExists(BucketName)
	if err != nil {
		return err
	}

	if !exists {
		err = s.minioClient.MakeBucket(BucketName, "")
		if err != nil {
			zap.L().Sugar().Errorf("Failed to create bucket %s: %v", BucketName, err)
			return err
		}
		zap.L().Sugar().Infof("Bucket %s created successfully", BucketName)
	}

	return nil
}

// UploadFile загружает зашифрованный файл в MinIO
// fileKey используется как ключ объекта в бакете
func (s *Storage) UploadFile(ctx context.Context, fileKey string, reader io.Reader, size int64) error {
	if s.minioClient == nil {
		return fmt.Errorf("minio client is not initialized")
	}

	_, err := s.minioClient.PutObjectWithContext(ctx, BucketName, fileKey, reader, size, minio.PutObjectOptions{})
	if err != nil {
		zap.L().Sugar().Errorf("Failed to upload file with key %s: %v", fileKey, err)
		return err
	}

	zap.L().Sugar().Debugf("File uploaded successfully with key: %s", fileKey)
	return nil
}

// DownloadFile скачивает зашифрованный файл из MinIO
// fileKey используется как ключ объекта в бакете
func (s *Storage) DownloadFile(ctx context.Context, fileKey string) (io.ReadCloser, error) {
	if s.minioClient == nil {
		return nil, fmt.Errorf("minio client is not initialized")
	}

	obj, err := s.minioClient.GetObjectWithContext(ctx, BucketName, fileKey, minio.GetObjectOptions{})
	if err != nil {
		zap.L().Sugar().Errorf("Failed to download file with key %s: %v", fileKey, err)
		return nil, err
	}

	// Проверяем, что объект существует
	stat, err := obj.Stat()
	if err != nil {
		obj.Close()
		zap.L().Sugar().Errorf("Failed to get file stats for key %s: %v", fileKey, err)
		return nil, err
	}

	if stat.Size == 0 {
		obj.Close()
		zap.L().Sugar().Errorf("File with key %s is empty or does not exist", fileKey)
		return nil, err
	}

	zap.L().Sugar().Debugf("File downloaded successfully with key: %s, size: %d", fileKey, stat.Size)
	return obj, nil
}

// DeleteFile удаляет файл из MinIO
func (s *Storage) DeleteFile(ctx context.Context, fileKey string) error {
	if s.minioClient == nil {
		return fmt.Errorf("minio client is not initialized")
	}

	err := s.minioClient.RemoveObject(BucketName, fileKey)
	if err != nil {
		zap.L().Sugar().Errorf("Failed to delete file with key %s: %v", fileKey, err)
		return err
	}

	zap.L().Sugar().Debugf("File deleted successfully with key: %s", fileKey)
	return nil
}

// FileExists проверяет существование файла
func (s *Storage) FileExists(ctx context.Context, fileKey string) (bool, error) {
	if s.minioClient == nil {
		return false, fmt.Errorf("minio client is not initialized")
	}

	_, err := s.minioClient.StatObject(BucketName, fileKey, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
