package usecase

import (
	"context"

	"gophkeeper/internal/auth"
	"gophkeeper/internal/models"

	"github.com/google/uuid"
)

// AuthenticateUserRepository определяет интерфейс для аутентификации пользователей.
// Интерфейс используется для абстракции от конкретной реализации хранилища
// и позволяет тестировать usecase с помощью моков.
type AuthenticateUserRepository interface {
	GetUserByLoginAndPassword(ctx context.Context, login string, passwordHash string) (userID uuid.UUID, err error)
}

// AuthenticateUsecase определяет интерфейс для аутентификации пользователей.
// Интерфейс содержит бизнес-логику проверки учетных данных и генерации
// JWT токенов доступа при успешной аутентификации.
type AuthenticateUsecase interface {
	Authenticate(ctx context.Context, login string, passwordHash string) (result models.AuthenticateOut, err error)
}

type authenticateUsecase struct {
	repo AuthenticateUserRepository
}

// NewAuthenticateUsecase создает новый экземпляр AuthenticateUsecase.
// Принимает репозиторий для аутентификации пользователей.
// Функция используется для внедрения зависимостей и создания usecase
// с конкретной реализацией хранилища пользователей.
func NewAuthenticateUsecase(repo AuthenticateUserRepository) AuthenticateUsecase {
	return &authenticateUsecase{
		repo: repo,
	}
}

func (uc *authenticateUsecase) Authenticate(ctx context.Context, login string, passwordHash string) (result models.AuthenticateOut, err error) {
	if len(login) < 1 || len(passwordHash) < 1 {
		return models.AuthenticateOut{}, models.ErrLoginOrPasswordEmpty
	}

	userID, err := uc.repo.GetUserByLoginAndPassword(ctx, login, passwordHash)
	if err != nil {
		return models.AuthenticateOut{}, err
	}

	// Если пользователь не найден
	if userID == uuid.Nil {
		return models.AuthenticateOut{}, models.ErrUserNotFound
	}

	return MakeAuthenticateOut(userID)
}

// MakeAuthenticateOut создает структуру ответа для аутентификации.
// Функция генерирует JWT access и refresh токены для указанного пользователя.
// Используется как вспомогательная функция для создания ответа аутентификации.
func MakeAuthenticateOut(userID uuid.UUID) (out models.AuthenticateOut, err error) {
	// Генерируем access token
	accessToken, err := auth.GenerateAccessToken(userID.String())
	if err != nil {
		return models.AuthenticateOut{}, err
	}

	// Генерируем refresh token
	refreshToken, err := auth.GenerateRefreshToken(userID.String())
	if err != nil {
		return models.AuthenticateOut{}, err
	}

	return models.AuthenticateOut{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
