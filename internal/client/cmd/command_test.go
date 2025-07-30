package cmd

import (
	"context"
	"errors"
	"testing"

	"gophkeeper/internal/client/crypto"
)

// MockCommand представляет мок команды для тестирования
type MockCommand struct {
	name        string
	description string
	executeFunc func(ctx context.Context, args any) (any, error)
}

func (m *MockCommand) Execute(ctx context.Context, args any) (any, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, args)
	}
	return nil, nil
}

func (m *MockCommand) GetName() string {
	return m.name
}

func (m *MockCommand) GetDescription() string {
	return m.description
}

func TestNewCommandRegistry(t *testing.T) {
	crypto := crypto.NewCrypto()
	registry := NewCommandRegistry(crypto)

	if registry == nil {
		t.Fatal("NewCommandRegistry returned nil")
	}

	if registry.commands == nil {
		t.Fatal("commands map is nil")
	}
}

func TestCommandRegistry_RegisterCommand(t *testing.T) {
	crypto := crypto.NewCrypto()
	registry := NewCommandRegistry(crypto)

	mockCmd := &MockCommand{
		name:        "test",
		description: "Test command",
	}

	registry.RegisterCommand(mockCmd)

	if len(registry.commands) != 1 {
		t.Errorf("Expected 1 command, got %d", len(registry.commands))
	}

	cmd, exists := registry.commands["test"]
	if !exists {
		t.Fatal("Command not found in registry")
	}

	if cmd.GetName() != "test" {
		t.Errorf("Expected name 'test', got '%s'", cmd.GetName())
	}
}

func TestCommandRegistry_GetCommand(t *testing.T) {
	crypto := crypto.NewCrypto()
	registry := NewCommandRegistry(crypto)

	mockCmd := &MockCommand{
		name:        "test",
		description: "Test command",
	}

	registry.RegisterCommand(mockCmd)

	// Тест успешного получения команды
	cmd, exists := registry.GetCommand("test")
	if !exists {
		t.Fatal("Command should exist")
	}

	if cmd.GetName() != "test" {
		t.Errorf("Expected name 'test', got '%s'", cmd.GetName())
	}

	// Тест получения несуществующей команды
	cmd, exists = registry.GetCommand("nonexistent")
	if exists {
		t.Fatal("Command should not exist")
	}

	if cmd != nil {
		t.Errorf("Expected nil command, got %v", cmd)
	}
}

func TestCommandRegistry_Execute(t *testing.T) {
	crypto := crypto.NewCrypto()
	registry := NewCommandRegistry(crypto)

	ctx := context.Background()

	// Тест выполнения несуществующей команды
	_, err := registry.Execute(ctx, "nonexistent", "args")
	if err == nil {
		t.Fatal("Expected error for nonexistent command")
	}

	// Тест успешного выполнения команды
	mockCmd := &MockCommand{
		name:        "test",
		description: "Test command",
		executeFunc: func(ctx context.Context, args any) (any, error) {
			return "result", nil
		},
	}

	registry.RegisterCommand(mockCmd)

	result, err := registry.Execute(ctx, "test", "args")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != "result" {
		t.Errorf("Expected result 'result', got %v", result)
	}

	// Тест выполнения команды с ошибкой
	errorCmd := &MockCommand{
		name:        "error",
		description: "Error command",
		executeFunc: func(ctx context.Context, args any) (any, error) {
			return nil, errors.New("test error")
		},
	}

	registry.RegisterCommand(errorCmd)

	_, err = registry.Execute(ctx, "error", "args")
	if err == nil {
		t.Fatal("Expected error")
	}

	if err.Error() != "test error" {
		t.Errorf("Expected error 'test error', got '%s'", err.Error())
	}
}
