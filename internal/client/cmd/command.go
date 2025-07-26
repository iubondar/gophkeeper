package cmd

import (
	"context"
	"fmt"
	"gophkeeper/internal/client/crypto"
)

// Command представляет универсальный интерфейс для всех команд
// args - это аргументы команды, которые передаются в метод Execute. Конкретная команда сама разбирает аргументы и приводит их к нужному типу.
// GetName - возвращает имя команды
// GetDescription - возвращает описание команды
type Command interface {
	Execute(ctx context.Context, args any) (any, error)
	GetName() string
	GetDescription() string
}

// CommandRegistry управляет реестром команд
type CommandRegistry struct {
	commands map[string]Command
}

// NewCommandRegistry создает новый реестр команд
func NewCommandRegistry(crypto *crypto.Crypto) *CommandRegistry {
	registry := CommandRegistry{
		commands: make(map[string]Command),
	}

	return &registry
}

// RegisterCommand регистрирует команду в реестре
func (r *CommandRegistry) RegisterCommand(cmd Command) {
	r.commands[cmd.GetName()] = cmd
}

// GetCommand возвращает команду по имени
func (r *CommandRegistry) GetCommand(name string) (Command, bool) {
	cmd, exists := r.commands[name]
	return cmd, exists
}

// Execute выполняет команду по имени
func (r *CommandRegistry) Execute(ctx context.Context, commandName string, args any) (any, error) {
	cmd, exists := r.GetCommand(commandName)
	if !exists {
		return nil, fmt.Errorf("команда '%s' не найдена", commandName)
	}

	return cmd.Execute(ctx, args)
}
