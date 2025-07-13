package usecase

import "errors"

var (
	ErrLoginOrPasswordEmpty = errors.New("login or password is empty")
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user already exists")
)
