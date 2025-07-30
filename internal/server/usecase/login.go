package usecase

import (
	"context"
	"fmt"
)

// LoginUserRepository определяет интерфейс для получения соли пользователя.
// Интерфейс используется для абстракции от конкретной реализации хранилища
// и позволяет тестировать usecase с помощью моков.
type LoginUserRepository interface {
	GetUserSalt(ctx context.Context, login string) (salt string, err error)
}

// LoginUsecase определяет интерфейс для получения соли пользователя при входе.
// Интерфейс содержит бизнес-логику для первого шага аутентификации,
// который возвращает соль для хеширования пароля на клиенте.
type LoginUsecase interface {
	GetSalt(ctx context.Context, login string) (salt string, err error)
}

type loginUsecase struct {
	repo LoginUserRepository
}

// NewLoginUsecase создает новый экземпляр LoginUsecase.
// Принимает репозиторий для работы с пользователями при входе.
// Функция используется для внедрения зависимостей и создания usecase
// с конкретной реализацией хранилища пользователей.
func NewLoginUsecase(repo LoginUserRepository) LoginUsecase {
	return &loginUsecase{
		repo: repo,
	}
}

func (uc *loginUsecase) GetSalt(ctx context.Context, login string) (salt string, err error) {
	if len(login) < 1 {
		return "", fmt.Errorf("пустой логин")
	}

	return uc.repo.GetUserSalt(ctx, login)
}
