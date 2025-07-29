package shell

import (
	"testing"
)

func TestNewMenuManager(t *testing.T) {
	manager := NewMenuManager()

	if manager == nil {
		t.Fatal("NewMenuManager returned nil")
	}

	if manager.currentState != MenuStateMain {
		t.Errorf("Expected initial state to be %s, got %s", MenuStateMain, manager.currentState)
	}

	if manager.actionState != nil {
		t.Error("Expected actionState to be nil initially")
	}

	// Проверяем, что все состояния созданы
	expectedStates := []string{MenuStateMain, MenuStateAuthenticated, MenuStateDataType}
	for _, state := range expectedStates {
		if _, exists := manager.states[state]; !exists {
			t.Errorf("Expected state %s to exist", state)
		}
	}
}

func TestMenuManager_ShowMenu(t *testing.T) {
	manager := NewMenuManager()

	// Тестируем отображение главного меню
	manager.currentState = MenuStateMain
	manager.ShowMenu() // Не должно паниковать

	// Тестируем отображение меню авторизованного пользователя
	manager.currentState = MenuStateAuthenticated
	manager.ShowMenu() // Не должно паниковать

	// Тестируем отображение меню выбора типа данных
	manager.currentState = MenuStateDataType
	manager.ShowMenu() // Не должно паниковать
}

func TestMenuManager_GetCommandByID(t *testing.T) {
	manager := NewMenuManager()

	// Тестируем получение команды из главного меню
	manager.currentState = MenuStateMain

	command, exists := manager.GetCommandByID("1")
	if !exists {
		t.Error("Expected command to exist for ID '1'")
	}
	if command != CommandRegister {
		t.Errorf("Expected command %s, got %s", CommandRegister, command)
	}

	command, exists = manager.GetCommandByID("2")
	if !exists {
		t.Error("Expected command to exist for ID '2'")
	}
	if command != CommandLogin {
		t.Errorf("Expected command %s, got %s", CommandLogin, command)
	}

	// Тестируем несуществующий ID
	command, exists = manager.GetCommandByID("999")
	if exists {
		t.Error("Expected command to not exist for ID '999'")
	}
	if command != "" {
		t.Errorf("Expected empty command, got %s", command)
	}
}

func TestMenuManager_IsExitCommand(t *testing.T) {
	manager := NewMenuManager()

	if !manager.IsExitCommand(CommandExit) {
		t.Error("Expected CommandExit to be recognized as exit command")
	}

	if manager.IsExitCommand(CommandLogin) {
		t.Error("Expected CommandLogin to not be recognized as exit command")
	}
}

func TestMenuManager_IsBackCommand(t *testing.T) {
	manager := NewMenuManager()

	if !manager.IsBackCommand(CommandBack) {
		t.Error("Expected CommandBack to be recognized as back command")
	}

	if manager.IsBackCommand(CommandLogin) {
		t.Error("Expected CommandLogin to not be recognized as back command")
	}
}

func TestMenuManager_IsLogoutCommand(t *testing.T) {
	manager := NewMenuManager()

	if !manager.IsLogoutCommand(CommandLogout) {
		t.Error("Expected CommandLogout to be recognized as logout command")
	}

	if manager.IsLogoutCommand(CommandLogin) {
		t.Error("Expected CommandLogin to not be recognized as logout command")
	}
}

func TestMenuManager_SwitchToState(t *testing.T) {
	manager := NewMenuManager()

	// Тестируем переключение на существующее состояние
	if !manager.SwitchToState(MenuStateAuthenticated) {
		t.Error("Expected SwitchToState to return true for existing state")
	}
	if manager.currentState != MenuStateAuthenticated {
		t.Errorf("Expected currentState to be %s, got %s", MenuStateAuthenticated, manager.currentState)
	}

	// Тестируем переключение на несуществующее состояние
	if manager.SwitchToState("non_existent_state") {
		t.Error("Expected SwitchToState to return false for non-existent state")
	}
	if manager.currentState != MenuStateAuthenticated {
		t.Error("Expected currentState to remain unchanged")
	}
}

func TestMenuManager_GoBack(t *testing.T) {
	manager := NewMenuManager()

	// Начинаем с главного меню
	manager.currentState = MenuStateMain

	// Пытаемся вернуться назад из главного меню (не должно работать)
	if manager.GoBack() {
		t.Error("Expected GoBack to return false when no parent state")
	}
	if manager.currentState != MenuStateMain {
		t.Error("Expected currentState to remain unchanged")
	}

	// Переходим в меню авторизованного пользователя
	manager.currentState = MenuStateAuthenticated

	// Возвращаемся назад
	if !manager.GoBack() {
		t.Error("Expected GoBack to return true when parent state exists")
	}
	if manager.currentState != MenuStateMain {
		t.Errorf("Expected currentState to be %s, got %s", MenuStateMain, manager.currentState)
	}
}

func TestMenuManager_GetCurrentState(t *testing.T) {
	manager := NewMenuManager()

	if manager.GetCurrentState() != MenuStateMain {
		t.Errorf("Expected GetCurrentState to return %s, got %s", MenuStateMain, manager.GetCurrentState())
	}

	manager.currentState = MenuStateAuthenticated
	if manager.GetCurrentState() != MenuStateAuthenticated {
		t.Errorf("Expected GetCurrentState to return %s, got %s", MenuStateAuthenticated, manager.GetCurrentState())
	}
}

func TestMenuManager_SetActionState(t *testing.T) {
	manager := NewMenuManager()

	// Устанавливаем состояние действия
	manager.SetActionState("upload", "text")

	if manager.actionState == nil {
		t.Error("Expected actionState to be set")
	}
	if manager.actionState.Action != "upload" {
		t.Errorf("Expected action to be 'upload', got %s", manager.actionState.Action)
	}
	if manager.actionState.Type != "text" {
		t.Errorf("Expected type to be 'text', got %s", manager.actionState.Type)
	}
}

func TestMenuManager_GetActionState(t *testing.T) {
	manager := NewMenuManager()

	// Изначально состояние действия должно быть nil
	if manager.GetActionState() != nil {
		t.Error("Expected GetActionState to return nil initially")
	}

	// Устанавливаем состояние действия
	manager.SetActionState("upload", "text")

	actionState := manager.GetActionState()
	if actionState == nil {
		t.Error("Expected GetActionState to return non-nil after setting")
		return
	}
	if actionState.Action != "upload" {
		t.Errorf("Expected action to be 'upload', got %s", actionState.Action)
	}
	if actionState.Type != "text" {
		t.Errorf("Expected type to be 'text', got %s", actionState.Type)
	}
}

func TestMenuManager_ClearActionState(t *testing.T) {
	manager := NewMenuManager()

	// Устанавливаем состояние действия
	manager.SetActionState("upload", "text")
	if manager.actionState == nil {
		t.Error("Expected actionState to be set")
	}

	// Очищаем состояние действия
	manager.ClearActionState()
	if manager.actionState != nil {
		t.Error("Expected actionState to be nil after clearing")
	}
}

func TestMenuManager_IsActionCommand(t *testing.T) {
	manager := NewMenuManager()

	if !manager.IsActionCommand(CommandUpload) {
		t.Error("Expected CommandUpload to be recognized as action command")
	}

	if manager.IsActionCommand(CommandLogin) {
		t.Error("Expected CommandLogin to not be recognized as action command")
	}
}

func TestMenuManager_IsDataTypeCommand(t *testing.T) {
	manager := NewMenuManager()

	// Тестируем все команды типа данных
	dataTypeCommands := []string{CommandText, CommandLoginPassword, CommandCard, CommandFile}
	for _, cmd := range dataTypeCommands {
		if !manager.IsDataTypeCommand(cmd) {
			t.Errorf("Expected %s to be recognized as data type command", cmd)
		}
	}

	// Тестируем команды, которые не являются командами типа данных
	nonDataTypeCommands := []string{CommandLogin, CommandRegister, CommandExit}
	for _, cmd := range nonDataTypeCommands {
		if manager.IsDataTypeCommand(cmd) {
			t.Errorf("Expected %s to not be recognized as data type command", cmd)
		}
	}
}
