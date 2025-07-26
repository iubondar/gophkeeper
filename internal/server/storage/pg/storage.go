package pg

import (
	"context"
	"database/sql"
	"errors"
	"gophkeeper/internal/models"
	"gophkeeper/internal/server/storage/queries"

	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(dsn string) (storage *Storage, err error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}

// CheckStatus проверяет состояние хранилища.
// Возвращает ошибку, если база данных недоступна.
func (s *Storage) CheckStatus(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *Storage) Register(ctx context.Context, userID uuid.UUID, login string, passwordHash string, salt string) (ok bool, err error) {
	_, err = s.db.ExecContext(ctx, queries.InsertUser, userID, login, passwordHash, salt)
	if err != nil {
		// Если пользователь с логином уже существует - возвращаем не ок
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return false, models.ErrUserAlreadyExists
			}
		}

		// Другая ошибка
		zap.L().Sugar().Debugln("Error insert new user:", err.Error())
		return false, err
	}

	return true, nil
}

func (s *Storage) GetUserSalt(ctx context.Context, login string) (salt string, err error) {
	err = s.db.QueryRowContext(ctx, queries.GetUserSalt, login).Scan(&salt)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		zap.L().Sugar().Debugln("Error getting user salt:", err.Error())
		return "", err
	}
	return salt, nil
}

func (s *Storage) GetUserByLoginAndPassword(ctx context.Context, login string, passwordHash string) (userID uuid.UUID, err error) {
	err = s.db.QueryRowContext(ctx, queries.GetUserByLoginAndPassword, login, passwordHash).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.Nil, nil
		}
		zap.L().Sugar().Debugln("Error getting user by login and password:", err.Error())
		return uuid.Nil, err
	}
	return userID, nil
}

func (s *Storage) InsertRecord(ctx context.Context, id, userID uuid.UUID, label, recordType, metadata string, encryptedData []byte, fileKey string, version int, createdAt, updatedAt time.Time) error {
	createdAtNull := sql.NullTime{Valid: true, Time: createdAt}
	updatedAtNull := sql.NullTime{Valid: true, Time: updatedAt}
	_, err := s.db.ExecContext(ctx, queries.InsertRecord, id, userID, label, recordType, metadata, encryptedData, fileKey, version, createdAtNull, updatedAtNull)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return models.ErrConflict
			}
		}
		zap.L().Sugar().Debugln("Error inserting record:", err.Error())
		return err
	}
	return nil
}

func (s *Storage) GetRecordByLabel(ctx context.Context, label string, userID uuid.UUID) (*models.GetSecretOut, error) {
	var record models.GetSecretOut
	var id uuid.UUID
	var createdAt, updatedAt sql.NullTime

	err := s.db.QueryRowContext(ctx, queries.GetRecordByLabel, label, userID).Scan(
		&id, &record.Label, &record.Type, &record.Metadata, &record.EncryptedData, &record.FileKey, &record.Version, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrRecordNotFound
		}
		zap.L().Sugar().Debugln("Error getting record by label:", err.Error())
		return nil, err
	}

	record.ID = id.String()
	return &record, nil
}

// DeleteRecordByLabel удаляет запись по label и userID
func (s *Storage) DeleteRecordByLabel(ctx context.Context, label string, userID uuid.UUID) error {
	result, err := s.db.ExecContext(ctx, queries.DeleteRecordByLabel, label, userID)
	if err != nil {
		zap.L().Sugar().Debugln("Error deleting record by label:", err.Error())
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return models.ErrRecordNotFound
	}
	return nil
}

// UpdateRecordByLabel обновляет запись с проверкой версии (OCC)
func (s *Storage) UpdateRecordByLabel(ctx context.Context, label string, userID uuid.UUID, recordType, metadata string, encryptedData []byte, fileKey string, expectedVersion int, updatedAt time.Time) (int, error) {
	result, err := s.db.ExecContext(
		ctx,
		queries.UpdateRecordByLabelAndVersion,
		recordType,
		metadata,
		encryptedData,
		fileKey,
		updatedAt,
		label,
		userID,
		expectedVersion,
	)
	if err != nil {
		zap.L().Sugar().Debugln("Error updating record by label:", err.Error())
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	if rowsAffected == 0 {
		return 0, models.ErrConflict // версия не совпала
	}
	// Получаем новую версию (expectedVersion + 1)
	return expectedVersion + 1, nil
}
