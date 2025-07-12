package storage

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	// _ "github.com/jackc/pgx/v5/stdlib"
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

func (s *Storage) Register(ctx context.Context, userID uuid.UUID, login string, passwordHash string) (ok bool, err error) {
	// _, err = s.db.ExecContext(ctx, queries.InsertUser, userID, login, passwordHash)
	// if err != nil {
	// 	// Если пользователь с логином уже существует - возвращаем не ок
	// 	var pgErr *pgconn.PgError
	// 	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
	// 		return false, nil
	// 	}

	// 	// Другая ошибка
	// 	zap.L().Sugar().Debugln("Error insert new user:", err.Error())
	// 	return false, err
	// }

	return true, nil
}
