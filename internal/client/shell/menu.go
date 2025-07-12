package shell

import (
	"fmt"
)

// MenuItem представляет элемент меню
type MenuItem struct {
	ID          string
	Title       string
	Description string
	Command     string
}

// MenuState представляет состояние меню
type MenuState struct {
	Title  string
	Items  []MenuItem
	Parent string // ID родительского состояния, пустая строка для корневого
}

// MenuManager управляет меню и навигацией
type MenuManager struct {
	states       map[string]MenuState
	currentState string
}

// NewMenuManager создает новый менеджер меню
func NewMenuManager() *MenuManager {
	manager := &MenuManager{
		states:       make(map[string]MenuState),
		currentState: "main",
	}

	// Главное меню
	manager.states["main"] = MenuState{
		Title: "Главное меню",
		Items: []MenuItem{
			{
				ID:          "1",
				Title:       "Регистрация",
				Description: "Создать новый аккаунт",
				Command:     "register",
			},
			{
				ID:          "2",
				Title:       "Вход",
				Description: "Войти в существующий аккаунт",
				Command:     "login",
			},
			{
				ID:          "3",
				Title:       "Выход",
				Description: "Завершить работу",
				Command:     "exit",
			},
		},
	}

	// Меню авторизованного пользователя
	manager.states["authenticated"] = MenuState{
		Title:  "Меню пользователя",
		Parent: "main",
		Items: []MenuItem{
			{
				ID:          "1",
				Title:       "Загрузить",
				Description: "Загрузить секрет",
				Command:     "upload",
			},
			{
				ID:          "2",
				Title:       "Обновить",
				Description: "Обновить секрет",
				Command:     "update",
			},
			{
				ID:          "3",
				Title:       "Получить",
				Description: "Получить секрет",
				Command:     "get",
			},
			{
				ID:          "4",
				Title:       "Удалить",
				Description: "Удалить секрет",
				Command:     "delete",
			},
			{
				ID:          "5",
				Title:       "Назад",
				Description: "Вернуться в главное меню",
				Command:     "back",
			},
		},
	}

	return manager
}

// ShowMenu отображает текущее меню
func (m *MenuManager) ShowMenu() {
	state, exists := m.states[m.currentState]
	if !exists {
		fmt.Println("Ошибка: состояние меню не найдено")
		return
	}

	fmt.Printf("=== %s ===\n", state.Title)
	for _, item := range state.Items {
		fmt.Printf("%s. %s - %s\n", item.ID, item.Title, item.Description)
	}
	fmt.Print("Выберите действие: ")
}

// GetCommandByID возвращает команду по ID пункта меню
func (m *MenuManager) GetCommandByID(id string) (string, bool) {
	state, exists := m.states[m.currentState]
	if !exists {
		return "", false
	}

	for _, item := range state.Items {
		if item.ID == id {
			return item.Command, true
		}
	}
	return "", false
}

// IsExitCommand проверяет, является ли команда командой выхода
func (m *MenuManager) IsExitCommand(command string) bool {
	return command == "exit"
}

// IsBackCommand проверяет, является ли команда командой возврата
func (m *MenuManager) IsBackCommand(command string) bool {
	return command == "back"
}

// SwitchToState переключает на другое состояние меню
func (m *MenuManager) SwitchToState(stateName string) bool {
	if _, exists := m.states[stateName]; exists {
		m.currentState = stateName
		return true
	}
	return false
}

// GetCurrentState возвращает текущее состояние
func (m *MenuManager) GetCurrentState() string {
	return m.currentState
}
