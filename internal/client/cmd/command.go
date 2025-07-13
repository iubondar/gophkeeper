package cmd

import (
	"context"
	"fmt"
	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/models"
)

// Command представляет универсальный интерфейс для всех команд
// args - это аргументы команды, которые передаются в метод Execute. Конкретная команда сама разбирает аргументы и приводит их к нужному типу.
// GetName - возвращает имя команды
// GetDescription - возвращает описание команды
type Command interface {
	Execute(ctx context.Context, args any) error
	GetName() string
	GetDescription() string
}

// GophKeeperClient это интерфейс для работы с API сервера, который предоставляет все API методы
type GophKeeperClient interface {
	Register(ctx context.Context, in models.RegisterIn) error
	Login(ctx context.Context, in models.LoginIn) (salt string, err error)
	Authenticate(ctx context.Context, in models.AuthenticateIn) error
	UploadSecret(ctx context.Context, secret models.SecretData) error
	UpdateSecret(ctx context.Context, secret models.SecretData) error
	GetSecret(ctx context.Context, secretName string) (*models.SecretData, error)
	DeleteSecret(ctx context.Context, secretName string) error
}

// CommandRegistry управляет реестром команд
type CommandRegistry struct {
	commands map[string]Command
}

// NewCommandRegistry создает новый реестр команд
func NewCommandRegistry(apiClient GophKeeperClient, crypto *crypto.Crypto) *CommandRegistry {
	registry := CommandRegistry{
		commands: make(map[string]Command),
	}

	registry.registerCommand(NewRegisterCommand(apiClient, crypto))
	registry.registerCommand(NewLoginCommand(apiClient, crypto))
	registry.registerCommand(NewUploadCommand(apiClient, crypto))
	registry.registerCommand(NewUpdateCommand(apiClient, crypto))
	registry.registerCommand(NewGetCommand(apiClient, crypto))
	registry.registerCommand(NewDeleteCommand(apiClient, crypto))

	return &registry
}

// RegisterCommand регистрирует команду в реестре
func (r *CommandRegistry) registerCommand(cmd Command) {
	r.commands[cmd.GetName()] = cmd
}

// GetCommand возвращает команду по имени
func (r *CommandRegistry) GetCommand(name string) (Command, bool) {
	cmd, exists := r.commands[name]
	return cmd, exists
}

// Execute выполняет команду по имени
func (r *CommandRegistry) Execute(ctx context.Context, commandName string, args any) error {
	cmd, exists := r.GetCommand(commandName)
	if !exists {
		return fmt.Errorf("команда '%s' не найдена", commandName)
	}

	return cmd.Execute(ctx, args)
}
