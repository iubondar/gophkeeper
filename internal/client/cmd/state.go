package cmd

import (
	"context"
	"fmt"
)

// StateCommand команда для переключения состояний меню
type StateCommand struct {
	state string
}

// NewStateCommand создает новую команду состояния
func NewStateCommand(state string) *StateCommand {
	return &StateCommand{state: state}
}

// Execute переключает состояние меню
func (c *StateCommand) Execute(ctx context.Context, args any) error {
	// В реальной реализации здесь была бы логика переключения состояния
	// Пока просто выводим сообщение
	fmt.Printf("Переключение в состояние: %s\n", c.state)
	return nil
}

// GetName возвращает имя команды
func (c *StateCommand) GetName() string {
	return c.state
}

// GetDescription возвращает описание команды
func (c *StateCommand) GetDescription() string {
	return fmt.Sprintf("Переключение в состояние %s", c.state)
}
