package usecase

import (
	"context"
	"gophkeeper/internal/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	Register(ctx context.Context, userID uuid.UUID, login string, passwordHash string, salt string) (ok bool, err error)
}

type RegisterUsecase interface {
	Register(ctx context.Context, in models.RegisterIn) (out models.AuthenticateOut, err error)
}

type registerUsecase struct {
	repo UserRepository
}

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
