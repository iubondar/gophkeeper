package storage

import (
	"context"
	"database/sql"
	"errors"
	"gophkeeper/internal/server/storage/queries"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
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
	_, err = s.db.ExecContext(ctx, queries.InsertUser, userID, login, passwordHash)
	if err != nil {
		// Если пользователь с логином уже существует - возвращаем не ок
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return false, nil
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
