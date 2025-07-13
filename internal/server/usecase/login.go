package usecase

import (
	"context"
	"fmt"
)

type LoginUserRepository interface {
	GetUserSalt(ctx context.Context, login string) (salt string, err error)
}

type LoginUsecase interface {
	GetSalt(ctx context.Context, login string) (salt string, err error)
}

type loginUsecase struct {
	repo LoginUserRepository
}

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
