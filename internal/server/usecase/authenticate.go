package usecase

import (
	"context"

	"gophkeeper/internal/auth"
	"gophkeeper/internal/models"

	"github.com/google/uuid"
)

type AuthenticateUserRepository interface {
	GetUserByLoginAndPassword(ctx context.Context, login string, passwordHash string) (userID uuid.UUID, err error)
}

type AuthenticateUsecase interface {
	Authenticate(ctx context.Context, login string, passwordHash string) (result models.AuthenticateOut, err error)
}

type authenticateUsecase struct {
	repo AuthenticateUserRepository
}

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
