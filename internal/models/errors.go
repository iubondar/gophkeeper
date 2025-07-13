package models

import "errors"

var (
	ErrLoginOrPasswordEmpty = errors.New("пустой логин или пароль")
	ErrUserNotFound         = errors.New("пользователь не найден")
	ErrUserAlreadyExists    = errors.New("пользователь уже существует")
)
