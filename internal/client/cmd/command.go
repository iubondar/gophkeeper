package cmd

import (
	"context"
	"fmt"
	"sync"
)

// Command представляет универсальный интерфейс для всех команд
type Command interface {
	Execute(ctx context.Context, args interface{}) error
	GetName() string
	GetDescription() string
}

// CommandRegistry управляет реестром команд
type CommandRegistry struct {
	commands map[string]Command
	mu       sync.RWMutex
}

// NewCommandRegistry создает новый реестр команд
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[string]Command),
	}
}

// RegisterCommand регистрирует команду в реестре
func (r *CommandRegistry) RegisterCommand(cmd Command) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.commands[cmd.GetName()] = cmd
}

// GetCommand возвращает команду по имени
func (r *CommandRegistry) GetCommand(name string) (Command, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cmd, exists := r.commands[name]
	return cmd, exists
}

// ListCommands возвращает список всех зарегистрированных команд
func (r *CommandRegistry) ListCommands() []Command {
	r.mu.RLock()
	defer r.mu.RUnlock()

	commands := make([]Command, 0, len(r.commands))
	for _, cmd := range r.commands {
		commands = append(commands, cmd)
	}
	return commands
}

// CommandExecutor выполняет команды
type CommandExecutor struct {
	registry *CommandRegistry
}

// NewCommandExecutor создает новый исполнитель команд
func NewCommandExecutor(registry *CommandRegistry) *CommandExecutor {
	return &CommandExecutor{registry: registry}
}

// Execute выполняет команду по имени
func (e *CommandExecutor) Execute(ctx context.Context, commandName string, args interface{}) error {
	cmd, exists := e.registry.GetCommand(commandName)
	if !exists {
		return fmt.Errorf("команда '%s' не найдена", commandName)
	}

	return cmd.Execute(ctx, args)
}
