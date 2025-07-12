package cmd

import (
	"context"
	"gophkeeper/internal/models"
)

// GophKeeperClient это интерфейс для работы с API сервера, который предоставляет все API методы
type GophKeeperClient interface {
	Register(ctx context.Context, in models.RegisterIn) error
	Login(ctx context.Context, in models.LoginIn) error
}

// CommandFactory это фабрика команд
type CommandFactory struct {
	apiClient GophKeeperClient
}

// NewCommandFactory создает новый экземпляр CommandFactory
func NewCommandFactory(apiClient GophKeeperClient) *CommandFactory {
	return &CommandFactory{apiClient: apiClient}
}

// CreateCommandRegistry создает и настраивает реестр команд
func (f *CommandFactory) CreateCommandRegistry() *CommandRegistry {
	registry := NewCommandRegistry()

	// Регистрируем команды
	registry.RegisterCommand(NewRegisterCommand(f.apiClient))
	registry.RegisterCommand(NewLoginCommand(f.apiClient))
	// TODO: Добавить другие команды по мере их создания

	return registry
}
