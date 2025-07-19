package storage

import (
	"context"
	"database/sql"
	"gophkeeper/internal/server/storage/testhelpers"
	"log"
	"testing"
	"time"

	"gophkeeper/internal/models"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StorageTestSuite struct {
	suite.Suite
	storage *Storage
	cleanup func()
}

func (s *StorageTestSuite) SetupSuite() {
	ctx := context.Background()
	container, err := testhelpers.CreatePostgresContainer(ctx)
	require.NoError(s.T(), err)

	db, err := sql.Open("pgx", container.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	goose.SetDialect("postgres")
	err = goose.Up(db, "./migrations")
	if err != nil {
		log.Fatal(err)
	}

	storage, err := NewStorage(container.ConnectionString)
	require.NoError(s.T(), err)

	s.storage = storage
	s.cleanup = func() {
		container.Terminate(ctx)
	}
}

func (s *StorageTestSuite) TearDownSuite() {
	if s.cleanup != nil {
		s.cleanup()
	}
}

func (s *StorageTestSuite) TearDownTest() {
	ctx := context.Background()
	_, err := s.storage.db.ExecContext(ctx, "TRUNCATE TABLE users CASCADE")
	require.NoError(s.T(), err)
	_, err = s.storage.db.ExecContext(ctx, "TRUNCATE TABLE records CASCADE")
	require.NoError(s.T(), err)
}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}

func (s *StorageTestSuite) TestRegister() {
	// ctx := context.Background()
	// userID := uuid.New()
	// login := "testuser"
	// password := "password123"
	// passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// require.NoError(s.T(), err)

	// s.Run("successful registration", func() {
	// 	ok, err := s.storage.Register(ctx, userID, login, string(passwordHash))
	// 	assert.NoError(s.T(), err)
	// 	assert.True(s.T(), ok)
	// })

	// s.Run("duplicate registration", func() {
	// 	ok, err := s.storage.Register(ctx, userID, login, string(passwordHash))
	// 	assert.NoError(s.T(), err)
	// 	assert.False(s.T(), ok)
	// })
}

func (s *StorageTestSuite) TestInsertRecord() {
	ctx := context.Background()
	userID := uuid.New()
	id := uuid.New()
	label := "test-label"
	recordType := "note"
	metadata := "meta"
	encryptedData := []byte("data")
	fileKey := "key"
	version := 1
	createdAt := time.Now()
	updatedAt := time.Now()

	// Добавляем пользователя, чтобы не было ошибки foreign key
	login := "testuser"
	passwordHash := "hash"
	salt := "salt"
	_, err := s.storage.Register(ctx, userID, login, passwordHash, salt)
	s.Require().NoError(err)

	// Тест успешной вставки
	err = s.storage.InsertRecord(ctx, id, userID, label, recordType, metadata, encryptedData, fileKey, version, createdAt, updatedAt)
	s.Require().NoError(err)

	// Тест конфликта - попытка вставить запись с тем же label
	id2 := uuid.New()
	err = s.storage.InsertRecord(ctx, id2, userID, label, recordType, metadata, encryptedData, fileKey, version, createdAt, updatedAt)
	s.Require().Error(err)
	s.Require().Equal(models.ErrConflict, err)
}
