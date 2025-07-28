package shell

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

// ActionState представляет состояние выбранного действия
type ActionState struct {
	Action string // upload, update, get, delete
	Type   string // text, login_password, card, file
}

// MenuManager управляет меню и навигацией
type MenuManager struct {
	states       map[string]MenuState
	currentState string
	actionState  *ActionState // Сохраняем выбранное действие и тип
	display      Display      // Инжектированный Display интерфейс
}

// NewMenuManager создает новый менеджер меню
func NewMenuManager() *MenuManager {
	manager := &MenuManager{
		states:       make(map[string]MenuState),
		currentState: MenuStateMain,
		actionState:  nil,
		display:      nil, // Будет установлен позже
	}

	// Главное меню
	manager.states[MenuStateMain] = MenuState{
		Title: "Главное меню",
		Items: []MenuItem{
			{
				ID:          "1",
				Title:       "Регистрация",
				Description: "Создать новый аккаунт",
				Command:     CommandRegister,
			},
			{
				ID:          "2",
				Title:       "Вход",
				Description: "Войти в существующий аккаунт",
				Command:     CommandLogin,
			},
			{
				ID:          "3",
				Title:       "Версия",
				Description: "Показать версию и дату сборки",
				Command:     CommandVersion,
			},
			{
				ID:          "4",
				Title:       "Выход",
				Description: "Завершить работу",
				Command:     CommandExit,
			},
		},
	}

	// Меню авторизованного пользователя
	manager.states[MenuStateAuthenticated] = MenuState{
		Title:  "Меню пользователя",
		Parent: MenuStateMain,
		Items: []MenuItem{
			{
				ID:          "1",
				Title:       "Загрузить",
				Description: "Загрузить секрет",
				Command:     CommandUpload,
			},
			{
				ID:          "2",
				Title:       "Обновить",
				Description: "Обновить секрет",
				Command:     CommandUpdate,
			},
			{
				ID:          "3",
				Title:       "Получить",
				Description: "Получить секрет",
				Command:     CommandGet,
			},
			{
				ID:          "4",
				Title:       "Удалить",
				Description: "Удалить секрет",
				Command:     CommandDelete,
			},
			{
				ID:          "5",
				Title:       "Выйти из аккаунта",
				Description: "Выйти из аккаунта",
				Command:     CommandLogout,
			},
		},
	}

	// Меню выбора типа данных
	manager.states[MenuStateDataType] = MenuState{
		Title:  "Выберите тип данных",
		Parent: MenuStateAuthenticated,
		Items: []MenuItem{
			{
				ID:          "1",
				Title:       "Текст",
				Description: "Простой текстовый секрет",
				Command:     CommandText,
			},
			{
				ID:          "2",
				Title:       "Логин-Пароль",
				Description: "Данные для входа в систему",
				Command:     CommandLoginPassword,
			},
			{
				ID:          "3",
				Title:       "Данные карты",
				Description: "Информация о банковской карте",
				Command:     CommandCard,
			},
			{
				ID:          "4",
				Title:       "Файл",
				Description: "Зашифрованный файл",
				Command:     CommandFile,
			},
			{
				ID:          "5",
				Title:       "Назад",
				Description: "Вернуться к предыдущему меню",
				Command:     CommandBack,
			},
		},
	}

	return manager
}

// SetDisplay устанавливает Display интерфейс для MenuManager
func (m *MenuManager) SetDisplay(display Display) {
	m.display = display
}

// ShowMenu отображает текущее меню
func (m *MenuManager) ShowMenu() {
	state, exists := m.states[m.currentState]
	if !exists {
		if m.display != nil {
			m.display.MenuError("Ошибка: состояние меню не найдено")
		}
		return
	}

	if m.display != nil {
		m.display.MenuTitle(state.Title)
		for _, item := range state.Items {
			m.display.MenuItem(item.ID, item.Title, item.Description)
		}
		m.display.MenuChoice()
	}
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
	return command == CommandExit
}

// IsBackCommand проверяет, является ли команда командой возврата
func (m *MenuManager) IsBackCommand(command string) bool {
	return command == CommandBack
}

// IsLogoutCommand проверяет, является ли команда командой выхода из аккаунта
func (m *MenuManager) IsLogoutCommand(command string) bool {
	return command == CommandLogout
}

// SwitchToState переключает на другое состояние меню
func (m *MenuManager) SwitchToState(stateName string) bool {
	if _, exists := m.states[stateName]; exists {
		m.currentState = stateName
		return true
	}
	return false
}

// GoBack переходит к предыдущему состоянию меню
func (m *MenuManager) GoBack() bool {
	state, exists := m.states[m.currentState]
	if !exists || state.Parent == "" {
		return false
	}
	m.currentState = state.Parent
	return true
}

// GetCurrentState возвращает текущее состояние
func (m *MenuManager) GetCurrentState() string {
	return m.currentState
}

// SetActionState устанавливает состояние выбранного действия
func (m *MenuManager) SetActionState(action, dataType string) {
	m.actionState = &ActionState{
		Action: action,
		Type:   dataType,
	}
}

// GetActionState возвращает текущее состояние действия
func (m *MenuManager) GetActionState() *ActionState {
	return m.actionState
}

// ClearActionState очищает состояние действия
func (m *MenuManager) ClearActionState() {
	m.actionState = nil
}

// IsActionCommand проверяет, является ли команда командой действия (upload)
func (m *MenuManager) IsActionCommand(command string) bool {
	return command == CommandUpload
}

// IsDataTypeCommand проверяет, является ли команда командой выбора типа данных
func (m *MenuManager) IsDataTypeCommand(command string) bool {
	return command == CommandText || command == CommandLoginPassword || command == CommandCard || command == CommandFile
}
