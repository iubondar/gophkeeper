package usecase

import (
	"context"

	"gophkeeper/internal/auth"
	"gophkeeper/internal/models"

	"github.com/google/uuid"
)

// RefreshUsecase определяет интерфейс для обновления токенов.
// Интерфейс содержит бизнес-логику валидации refresh токена и генерации
// новых access и refresh токенов.
type RefreshUsecase interface {
	Refresh(ctx context.Context, refreshToken string) (result models.RefreshOut, err error)
}

type refreshUsecase struct{}

// NewRefreshUsecase создает новый экземпляр RefreshUsecase.
// Функция используется для внедрения зависимостей и создания usecase
// с конкретной реализацией бизнес-логики обновления токенов.
func NewRefreshUsecase() RefreshUsecase {
	return &refreshUsecase{}
}

func (uc *refreshUsecase) Refresh(ctx context.Context, refreshToken string) (result models.RefreshOut, err error) {
	if len(refreshToken) < 1 {
		return models.RefreshOut{}, models.ErrRefreshTokenInvalid
	}

	// Валидируем refresh token и получаем userID
	userID, err := auth.ValidateRefreshToken(refreshToken)
	if err != nil {
		return models.RefreshOut{}, models.ErrRefreshTokenExpired
	}

	// Если пользователь не найден
	if userID == uuid.Nil {
		return models.RefreshOut{}, models.ErrRefreshTokenInvalid
	}

	return makeRefreshOut(userID)
}

// MakeRefreshOut создает структуру ответа для обновления токенов.
// Функция генерирует новые JWT access и refresh токены для указанного пользователя.
// Используется как вспомогательная функция для создания ответа обновления токенов.
func makeRefreshOut(userID uuid.UUID) (out models.RefreshOut, err error) {
	// Генерируем новый access token
	accessToken, err := auth.GenerateAccessToken(userID.String())
	if err != nil {
		return models.RefreshOut{}, err
	}

	// Генерируем новый refresh token
	refreshToken, err := auth.GenerateRefreshToken(userID.String())
	if err != nil {
		return models.RefreshOut{}, err
	}

	return models.RefreshOut{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
