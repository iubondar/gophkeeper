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

func (s *StorageTestSuite) TestCheckStatus() {
	ctx := context.Background()

	s.Run("successful status check", func() {
		err := s.storage.CheckStatus(ctx)
		s.Require().NoError(err)
	})
}

func (s *StorageTestSuite) TestRegister() {
	ctx := context.Background()
	userID := uuid.New()
	login := "testuser"
	passwordHash := "password123"
	salt := "salt123"

	s.Run("successful registration", func() {
		ok, err := s.storage.Register(ctx, userID, login, passwordHash, salt)
		s.Require().NoError(err)
		s.Require().True(ok)
	})

	s.Run("duplicate registration with same login", func() {
		userID2 := uuid.New()
		ok, err := s.storage.Register(ctx, userID2, login, passwordHash, salt)
		s.Require().Error(err)
		s.Require().False(ok)
		s.Require().Equal(models.ErrUserAlreadyExists, err)
	})

	s.Run("duplicate registration with same user ID", func() {
		ok, err := s.storage.Register(ctx, userID, "differentlogin", passwordHash, salt)
		s.Require().Error(err)
		s.Require().False(ok)
		// Это должно вызвать ошибку primary key violation
	})
}

func (s *StorageTestSuite) TestGetUserSalt() {
	ctx := context.Background()
	userID := uuid.New()
	login := "testuser"
	passwordHash := "password123"
	salt := "salt123"

	s.Run("get salt for existing user", func() {
		// Сначала регистрируем пользователя
		_, err := s.storage.Register(ctx, userID, login, passwordHash, salt)
		s.Require().NoError(err)

		// Получаем соль
		retrievedSalt, err := s.storage.GetUserSalt(ctx, login)
		s.Require().NoError(err)
		s.Require().Equal(salt, retrievedSalt)
	})

	s.Run("get salt for non-existing user", func() {
		retrievedSalt, err := s.storage.GetUserSalt(ctx, "nonexistent")
		s.Require().NoError(err)
		s.Require().Equal("", retrievedSalt)
	})
}

func (s *StorageTestSuite) TestGetUserByLoginAndPassword() {
	ctx := context.Background()
	userID := uuid.New()
	login := "testuser"
	passwordHash := "password123"
	salt := "salt123"

	s.Run("get user with correct credentials", func() {
		// Сначала регистрируем пользователя
		_, err := s.storage.Register(ctx, userID, login, passwordHash, salt)
		s.Require().NoError(err)

		// Получаем пользователя по логину и паролю
		retrievedUserID, err := s.storage.GetUserByLoginAndPassword(ctx, login, passwordHash)
		s.Require().NoError(err)
		s.Require().Equal(userID, retrievedUserID)
	})

	s.Run("get user with incorrect password", func() {
		retrievedUserID, err := s.storage.GetUserByLoginAndPassword(ctx, login, "wrongpassword")
		s.Require().NoError(err)
		s.Require().Equal(uuid.Nil, retrievedUserID)
	})

	s.Run("get user with non-existing login", func() {
		retrievedUserID, err := s.storage.GetUserByLoginAndPassword(ctx, "nonexistent", passwordHash)
		s.Require().NoError(err)
		s.Require().Equal(uuid.Nil, retrievedUserID)
	})
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

	s.Run("successful record insertion", func() {
		err = s.storage.InsertRecord(ctx, id, userID, label, recordType, metadata, encryptedData, fileKey, version, createdAt, updatedAt)
		s.Require().NoError(err)
	})

	s.Run("conflict on duplicate label", func() {
		id2 := uuid.New()
		err = s.storage.InsertRecord(ctx, id2, userID, label, recordType, metadata, encryptedData, fileKey, version, createdAt, updatedAt)
		s.Require().Error(err)
		s.Require().Equal(models.ErrConflict, err)
	})

	s.Run("conflict on duplicate record ID", func() {
		err = s.storage.InsertRecord(ctx, id, userID, "different-label", recordType, metadata, encryptedData, fileKey, version, createdAt, updatedAt)
		s.Require().Error(err)
		// Это должно вызвать ошибку primary key violation
	})

	s.Run("insert record with different user", func() {
		userID2 := uuid.New()
		login2 := "testuser2"
		_, err := s.storage.Register(ctx, userID2, login2, passwordHash, salt)
		s.Require().NoError(err)

		id3 := uuid.New()
		label2 := "test-label-2"
		err = s.storage.InsertRecord(ctx, id3, userID2, label2, recordType, metadata, encryptedData, fileKey, version, createdAt, updatedAt)
		s.Require().NoError(err)
	})

	s.Run("insert record with same label for different user", func() {
		userID3 := uuid.New()
		login3 := "testuser3"
		_, err := s.storage.Register(ctx, userID3, login3, passwordHash, salt)
		s.Require().NoError(err)

		id4 := uuid.New()
		// Используем тот же label, но для другого пользователя
		err = s.storage.InsertRecord(ctx, id4, userID3, label, recordType, metadata, encryptedData, fileKey, version, createdAt, updatedAt)
		s.Require().Error(err)
		s.Require().Equal(models.ErrConflict, err)
	})
}

func (s *StorageTestSuite) TestGetRecordByLabel() {
	ctx := context.Background()
	userID := uuid.New()
	login := "testuser"
	passwordHash := "hash"
	salt := "salt"
	_, err := s.storage.Register(ctx, userID, login, passwordHash, salt)
	s.Require().NoError(err)

	// Создаем тестовую запись
	recordID := uuid.New()
	label := "test-secret"
	recordType := models.SecretTypeText
	metadata := "test metadata"
	encryptedData := []byte("encrypted-data")
	fileKey := ""
	version := 1
	createdAt := time.Now()
	updatedAt := time.Now()

	err = s.storage.InsertRecord(ctx, recordID, userID, label, recordType, metadata, encryptedData, fileKey, version, createdAt, updatedAt)
	s.Require().NoError(err)

	s.Run("successful get record", func() {
		result, err := s.storage.GetRecordByLabel(ctx, label, userID)
		s.Require().NoError(err)
		s.Require().NotNil(result)
		s.Require().Equal(recordID.String(), result.ID)
		s.Require().Equal(label, result.Label)
		s.Require().Equal(recordType, result.Type)
		s.Require().Equal(metadata, result.Metadata)
		s.Require().Equal(encryptedData, result.EncryptedData)
		s.Require().Equal(fileKey, result.FileKey)
		s.Require().Equal(version, result.Version)
	})

	s.Run("record not found", func() {
		result, err := s.storage.GetRecordByLabel(ctx, "non-existent", userID)
		s.Require().Error(err)
		s.Require().Equal(models.ErrRecordNotFound, err)
		s.Require().Nil(result)
	})

	s.Run("wrong user", func() {
		result, err := s.storage.GetRecordByLabel(ctx, label, uuid.New())
		s.Require().Error(err)
		s.Require().Equal(models.ErrRecordNotFound, err)
		s.Require().Nil(result)
	})
}
