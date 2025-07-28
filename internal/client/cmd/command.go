// Package cmd предоставляет систему команд для клиента GophKeeper.
// Включает интерфейс Command, реестр команд и реализации конкретных команд
// для работы с секретами, файлами и аутентификацией.
package cmd

import (
	"context"
	"fmt"
	"gophkeeper/internal/client/crypto"
)

// Command представляет универсальный интерфейс для всех команд.
// Каждая команда должна реализовывать этот интерфейс для интеграции
// с системой команд клиента.
type Command interface {
	// Execute выполняет команду с заданными аргументами.
	// args - это аргументы команды, которые передаются в метод Execute.
	// Конкретная команда сама разбирает аргументы и приводит их к нужному типу.
	//
	// Параметры:
	//   - ctx: контекст выполнения
	//   - args: аргументы команды (тип зависит от конкретной команды)
	//
	// Возвращает:
	//   - any: результат выполнения команды
	//   - error: ошибка в случае неудачи
	Execute(ctx context.Context, args any) (any, error)

	// GetName возвращает уникальное имя команды.
	//
	// Возвращает:
	//   - string: имя команды
	GetName() string

	// GetDescription возвращает описание команды для отображения пользователю.
	//
	// Возвращает:
	//   - string: описание команды
	GetDescription() string
}

// CommandRegistry управляет реестром команд.
// Позволяет регистрировать, получать и выполнять команды по имени.
type CommandRegistry struct {
	commands map[string]Command
}

// NewCommandRegistry создает новый реестр команд.
//
// Параметры:
//   - crypto: экземпляр криптографического модуля
//
// Возвращает:
//   - *CommandRegistry: новый реестр команд
func NewCommandRegistry(crypto *crypto.Crypto) *CommandRegistry {
	registry := CommandRegistry{
		commands: make(map[string]Command),
	}

	return &registry
}

// RegisterCommand регистрирует команду в реестре.
// Если команда с таким именем уже существует, она будет перезаписана.
//
// Параметры:
//   - cmd: команда для регистрации
func (r *CommandRegistry) RegisterCommand(cmd Command) {
	r.commands[cmd.GetName()] = cmd
}

// GetCommand возвращает команду по имени.
//
// Параметры:
//   - name: имя команды
//
// Возвращает:
//   - Command: найденная команда или nil
//   - bool: true если команда найдена, false в противном случае
func (r *CommandRegistry) GetCommand(name string) (Command, bool) {
	cmd, exists := r.commands[name]
	return cmd, exists
}

// Execute выполняет команду по имени с заданными аргументами.
//
// Параметры:
//   - ctx: контекст выполнения
//   - commandName: имя команды для выполнения
//   - args: аргументы команды
//
// Возвращает:
//   - any: результат выполнения команды
//   - error: ошибка в случае неудачи или если команда не найдена
func (r *CommandRegistry) Execute(ctx context.Context, commandName string, args any) (any, error) {
	cmd, exists := r.GetCommand(commandName)
	if !exists {
		return nil, fmt.Errorf("команда '%s' не найдена", commandName)
	}

	return cmd.Execute(ctx, args)
}
