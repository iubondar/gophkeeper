package shell

import (
	"context"
	"fmt"
	"gophkeeper/internal/client/cmd"
	"gophkeeper/internal/models"
)

// Shell представляет интерактивный интерфейс для работы с GophKeeper
type Shell struct {
	commandRegistry *cmd.CommandRegistry
	menuManager     *MenuManager
	inputHandler    *InputHandler
}

// NewShell создает новый экземпляр Shell
func NewShell(registry *cmd.CommandRegistry) *Shell {
	return &Shell{
		commandRegistry: registry,
		menuManager:     NewMenuManager(),
		inputHandler:    NewInputHandler(),
	}
}

// Run запускает интерактивный интерфейс
func (s *Shell) Run() error {
	welcome()

	// Проверяем доступность сервера
	if err := s.checkServerHealth(); err != nil {
		return fmt.Errorf("ошибка подключения к серверу: %w", err)
	}

	serverConnected()

	// Основной цикл меню
	for {
		s.menuManager.ShowMenu()
		choice := s.inputHandler.GetUserChoice()

		commandName, exists := s.menuManager.GetCommandByID(choice)
		if !exists {
			invalidChoice()
			continue
		}

		switch {
		case s.menuManager.IsExitCommand(commandName):
			return s.handleExitCommand()
		case s.menuManager.IsLogoutCommand(commandName):
			s.handleLogoutCommand()
			continue
		case s.menuManager.IsBackCommand(commandName):
			s.handleBackCommand()
			continue
		case s.menuManager.IsActionCommand(commandName):
			s.handleActionCommand(commandName)
			continue
		case s.menuManager.IsDataTypeCommand(commandName):
			s.handleDataTypeCommand(commandName)
			continue
		default:
			if result, err := s.executeCommand(commandName); err != nil {
				errorMsg(err)
			} else {
				s.handleCommandResult(commandName, result)
				if (commandName == CommandRegister || commandName == CommandLogin) && s.menuManager.GetCurrentState() == MenuStateMain {
					switchToUserMenuNotice()
					s.menuManager.SwitchToState(MenuStateAuthenticated)
				}
			}
		}
	}
}

// executeCommand выполняет команду с соответствующими аргументами
func (s *Shell) executeCommand(commandName string) (any, error) {
	ctx := context.Background()

	switch commandName {
	case CommandRegister:
		credentials, err := s.inputHandler.GetUserCredentials()
		if err != nil {
			return nil, fmt.Errorf("ошибка получения данных пользователя: %w", err)
		}
		return s.commandRegistry.Execute(ctx, commandName, *credentials)

	case CommandLogin:
		credentials, err := s.inputHandler.GetLoginCredentials()
		if err != nil {
			return nil, fmt.Errorf("ошибка получения данных пользователя: %w", err)
		}
		return s.commandRegistry.Execute(ctx, commandName, *credentials)

	default:
		return nil, fmt.Errorf("неизвестная команда: %s", commandName)
	}
}

// executeActionWithType выполняет действие с выбранным типом данных через команды
func (s *Shell) executeActionWithType(action, dataType string) (any, error) {
	ctx := context.Background()

	// Получаем данные в зависимости от типа
	var data any
	var err error

	switch dataType {
	case CommandText:
		data, err = s.inputHandler.GetTextData()
	case CommandLoginPassword:
		data, err = s.inputHandler.GetLoginPasswordData()
	case CommandCard:
		data, err = s.inputHandler.GetCardData()
	case CommandFile:
		data, err = s.inputHandler.GetFileData()
	default:
		return nil, fmt.Errorf("неизвестный тип данных: %s", dataType)
	}

	if err != nil {
		return nil, fmt.Errorf("ошибка получения данных: %w", err)
	}

	// Для операций get и delete нужен только название секрета
	if action == CommandGet || action == CommandDelete {
		secretName, err := s.inputHandler.GetSecretName()
		if err != nil {
			return nil, fmt.Errorf("ошибка получения названия секрета: %w", err)
		}
		data = secretName
	}

	// Выполняем команду через реестр команд
	return s.commandRegistry.Execute(ctx, action, data)
}

// handleCommandResult обрабатывает результат выполнения команды
func (s *Shell) handleCommandResult(commandName string, result any) {
	// Обрабатываем специальные случаи для команды get
	if commandName == CommandGet {
		if getResult, ok := result.(*cmd.GetSecretResult); ok {
			s.displaySecretResult(getResult)
			return
		}
	}

	// Для остальных команд просто показываем сообщение об успехе
	s.printSuccessMessage(commandName)
}

// displaySecretResult отображает результат получения секрета
func (s *Shell) displaySecretResult(result *cmd.GetSecretResult) {
	switch result.Type {
	case "text":
		if textSecret, ok := result.Data.(*models.TextSecretData); ok {
			DisplayTextSecret(textSecret, result.Metadata)
		}
	case "login_password":
		if loginPassword, ok := result.Data.(*models.LoginPasswordData); ok {
			DisplayLoginPassword(loginPassword, result.Metadata)
		}
	case "card":
		if cardData, ok := result.Data.(*models.CardData); ok {
			DisplayCardData(cardData, result.Metadata)
		}
	}
}

func (s *Shell) printSuccessMessage(commandName string) {
	switch commandName {
	case CommandRegister:
		successRegistration()
	case CommandLogin:
		successLogin()
	case CommandUpload:
		successUpload()
	case CommandUpdate:
		successUpdate()
	case CommandGet:
		successGet()
	case CommandDelete:
		successDelete()
	default:
		successGeneric(commandName)
	}
}

// checkServerHealth проверяет доступность сервера
func (s *Shell) checkServerHealth() error {
	_, err := s.commandRegistry.Execute(context.Background(), "health", nil)
	return err
}

// Добавляем приватные методы-обработчики
func (s *Shell) handleExitCommand() error {
	goodbye()
	return nil
}

func (s *Shell) handleLogoutCommand() {
	s.menuManager.SwitchToState(MenuStateMain)
	s.menuManager.ClearActionState()
	logout()
}

func (s *Shell) handleBackCommand() {
	if !s.menuManager.GoBack() {
		backNotAllowed()
	}
}

func (s *Shell) handleActionCommand(commandName string) {
	s.menuManager.SetActionState(commandName, "")
	s.menuManager.SwitchToState(MenuStateDataType)
}

func (s *Shell) handleDataTypeCommand(commandName string) {
	actionState := s.menuManager.GetActionState()
	if actionState != nil {
		actionState.Type = commandName
		if result, err := s.executeActionWithType(actionState.Action, actionState.Type); err != nil {
			errorMsg(err)
		} else {
			s.handleCommandResult(actionState.Action, result)
		}
		s.menuManager.ClearActionState()
		s.menuManager.SwitchToState(MenuStateAuthenticated)
	}
}
