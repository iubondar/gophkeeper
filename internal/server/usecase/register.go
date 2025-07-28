// Package usecase предоставляет бизнес-логику для сервера GophKeeper.
// Пакет содержит usecase-слой, который реализует основные операции:
// аутентификация, управление секретами, загрузка/скачивание файлов.
// Все usecase используют интерфейсы репозиториев для абстракции от хранилища.
package usecase

import (
	"context"
	"gophkeeper/internal/models"

	"github.com/google/uuid"
)

// UserRepository определяет интерфейс для работы с пользователями в хранилище.
// Интерфейс используется для абстракции от конкретной реализации хранилища
// и позволяет тестировать usecase с помощью моков.
type UserRepository interface {
	Register(ctx context.Context, userID uuid.UUID, login string, passwordHash string, salt string) (ok bool, err error)
}

// RegisterUsecase определяет интерфейс для регистрации новых пользователей.
// Интерфейс содержит бизнес-логику регистрации пользователей с валидацией
// входных данных и генерацией JWT токенов при успешной регистрации.
type RegisterUsecase interface {
	Register(ctx context.Context, in models.RegisterIn) (out models.AuthenticateOut, err error)
}

type registerUsecase struct {
	repo UserRepository
}

// NewRegisterUsecase создает новый экземпляр RegisterUsecase.
// Принимает репозиторий для работы с пользователями.
// Функция используется для внедрения зависимостей и создания usecase
// с конкретной реализацией хранилища пользователей.
func NewRegisterUsecase(repo UserRepository) RegisterUsecase {
	return &registerUsecase{
		repo: repo,
	}
}

func (uc *registerUsecase) Register(ctx context.Context, in models.RegisterIn) (out models.AuthenticateOut, err error) {
	if len(in.Login) < 1 || len(in.PasswordHash) < 1 {
		return models.AuthenticateOut{}, models.ErrLoginOrPasswordEmpty
	}

	userID := uuid.New()
	ok, err := uc.repo.Register(ctx, userID, in.Login, in.PasswordHash, in.Salt)
	if err != nil {
		return models.AuthenticateOut{}, err
	}

	if !ok {
		return models.AuthenticateOut{}, models.ErrUserAlreadyExists
	}

	return MakeAuthenticateOut(userID)
}
