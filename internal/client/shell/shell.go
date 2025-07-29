// Package shell предоставляет интерактивный интерфейс для работы с GophKeeper.
// Включает управление меню, обработку команд и взаимодействие с пользователем
// через консольный интерфейс.
package shell

import (
	"context"
	"errors"
	"fmt"
	"gophkeeper/internal/client/cmd"
	"gophkeeper/internal/client/display"
	"gophkeeper/internal/client/input"
	"gophkeeper/internal/models"
)

// Shell представляет интерактивный интерфейс для работы с GophKeeper.
// Управляет меню, обработкой команд и взаимодействием с пользователем.
type Shell struct {
	commandRegistry CommandRegistry
	menuManager     *MenuManager
	inputHandler    InputHandler
	display         Display
	version         string
	buildTime       string
}

// NewShell создает новый экземпляр Shell с указанными зависимостями.
//
// Параметры:
//   - registry: реестр команд для выполнения операций
//   - inputHandler: обработчик пользовательского ввода
//   - display: интерфейс для отображения информации
//   - version: версия приложения
//   - buildTime: дата и время сборки
//
// Возвращает:
//   - *Shell: новый экземпляр интерактивного интерфейса
func NewShell(registry CommandRegistry, inputHandler InputHandler, display Display, version, buildTime string) *Shell {
	menuManager := NewMenuManager()
	menuManager.SetDisplay(display)

	return &Shell{
		commandRegistry: registry,
		menuManager:     menuManager,
		inputHandler:    inputHandler,
		display:         display,
		version:         version,
		buildTime:       buildTime,
	}
}

// isAuthError определяет, является ли ошибка ошибкой авторизации (401).
// Проверяет на специальную ошибку ErrUnauthorized.
//
// Параметры:
//   - err: ошибка для проверки
//
// Возвращает:
//   - bool: true если это ошибка авторизации
func isAuthError(err error) bool {
	return errors.Is(err, models.ErrUnauthorized)
}

// NewShellWithDefaults создает Shell с дефолтными зависимостями.
// Использует стандартные реализации InputHandler и Display.
//
// Параметры:
//   - registry: реестр команд для выполнения операций
//   - version: версия приложения
//   - buildTime: дата и время сборки
//
// Возвращает:
//   - *Shell: новый экземпляр интерактивного интерфейса
func NewShellWithDefaults(registry *cmd.CommandRegistry, version, buildTime string) *Shell {
	return NewShell(registry, input.NewInputHandler(), display.NewDisplay(), version, buildTime)
}

// Run запускает интерактивный интерфейс.
// Выполняет проверку доступности сервера и запускает основной цикл меню.
//
// Возвращает:
//   - error: ошибка в случае неудачи или при выходе из приложения
func (s *Shell) Run() error {
	s.display.Welcome()

	// Проверяем доступность сервера
	if err := s.checkServerHealth(); err != nil {
		return fmt.Errorf("ошибка подключения к серверу: %w", err)
	}

	s.display.ServerConnected()

	// Основной цикл меню
	for {
		s.menuManager.ShowMenu()
		choice := s.inputHandler.GetUserChoice()

		commandName, exists := s.menuManager.GetCommandByID(choice)
		if !exists {
			s.display.InvalidChoice()
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
		case commandName == CommandVersion:
			s.handleVersionCommand()
			continue
		case commandName == CommandGet || commandName == CommandDelete:
			s.handleGetDeleteCommand(commandName)
			continue
		case s.menuManager.IsActionCommand(commandName):
			s.handleActionCommand(commandName)
			continue
		case commandName == CommandUpdate:
			s.handleUpdateCommand()
			continue
		case s.menuManager.IsDataTypeCommand(commandName):
			s.handleDataTypeCommand(commandName)
			continue
		default:
			if result, err := s.executeCommand(commandName); err != nil {
				// Проверяем, является ли ошибка ошибкой авторизации
				if isAuthError(err) {
					s.display.AuthError()
					s.handleLogoutCommand()
				} else {
					s.display.ErrorMsg(err)
				}
			} else {
				s.handleCommandResult(commandName, result)
				if (commandName == CommandRegister || commandName == CommandLogin) && s.menuManager.GetCurrentState() == MenuStateMain {
					s.display.SwitchToUserMenuNotice()
					s.menuManager.SwitchToState(MenuStateAuthenticated)
				}
			}
		}
	}
}

// executeCommand выполняет команду с соответствующими аргументами.
// Получает данные от пользователя в зависимости от типа команды.
//
// Параметры:
//   - commandName: имя команды для выполнения
//
// Возвращает:
//   - any: результат выполнения команды
//   - error: ошибка в случае неудачи
func (s *Shell) executeCommand(commandName string) (any, error) {
	ctx := context.Background()

	var data any
	var err error

	switch commandName {
	case CommandRegister:
		credentials, err := s.inputHandler.GetUserCredentials()
		if err != nil {
			return nil, fmt.Errorf("ошибка получения данных пользователя: %w", err)
		}
		data = *credentials

	case CommandLogin:
		credentials, err := s.inputHandler.GetUserCredentials()
		if err != nil {
			return nil, fmt.Errorf("ошибка получения данных пользователя: %w", err)
		}
		data = *credentials

	case CommandGet, CommandDelete:
		secretName, err := s.inputHandler.GetSecretName()
		if err != nil {
			return nil, fmt.Errorf("ошибка получения названия секрета: %w", err)
		}
		data = secretName

	case CommandUpload:
		// Для upload нужен тип данных, который должен быть передан через контекст меню
		actionState := s.menuManager.GetActionState()
		if actionState == nil || actionState.Type == "" {
			return nil, fmt.Errorf("тип данных не выбран")
		}

		switch actionState.Type {
		case CommandText:
			data, err = s.inputHandler.GetTextData()
		case CommandLoginPassword:
			data, err = s.inputHandler.GetLoginPasswordData()
		case CommandCard:
			data, err = s.inputHandler.GetCardData()
		case CommandFile:
			data, err = s.inputHandler.GetFileData()
		default:
			return nil, fmt.Errorf("неизвестный тип данных: %s", actionState.Type)
		}

		if err != nil {
			return nil, fmt.Errorf("ошибка получения данных: %w", err)
		}

	default:
		return nil, fmt.Errorf("неизвестная команда: %s", commandName)
	}

	// Выполняем команду через реестр команд
	return s.commandRegistry.Execute(ctx, commandName, data)
}

// handleCommandResult обрабатывает результат выполнения команды.
// Отображает результат в зависимости от типа команды.
//
// Параметры:
//   - commandName: имя выполненной команды
//   - result: результат выполнения команды
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

// displaySecretResult отображает результат получения секрета.
// Показывает данные секрета в зависимости от его типа.
//
// Параметры:
//   - result: результат получения секрета
func (s *Shell) displaySecretResult(result *cmd.GetSecretResult) {
	switch result.Type {
	case models.SecretTypeText:
		if textSecret, ok := result.Data.(*models.TextSecretData); ok {
			s.display.DisplayTextSecret(textSecret, result.Metadata)
		}
	case models.SecretTypeLoginPassword:
		if loginPassword, ok := result.Data.(*models.LoginPasswordData); ok {
			s.display.DisplayLoginPassword(loginPassword, result.Metadata)
		}
	case models.SecretTypeCard:
		if cardData, ok := result.Data.(*models.CardData); ok {
			s.display.DisplayCardData(cardData, result.Metadata)
		}
	case models.SecretTypeFile:
		if fileData, ok := result.Data.(*models.FileData); ok {
			s.display.DisplayFileData(fileData, result.Metadata)
			// Для файлов дополнительно показываем сообщение о скачивании
			s.display.SuccessDownloadFile(fileData.FilePath)
		}
	}
}

// printSuccessMessage отображает сообщение об успешном выполнении команды.
//
// Параметры:
//   - commandName: имя выполненной команды
func (s *Shell) printSuccessMessage(commandName string) {
	switch commandName {
	case CommandRegister:
		s.display.SuccessRegistration()
	case CommandLogin:
		s.display.SuccessLogin()
	case CommandUpload:
		s.display.SuccessUpload()
	case CommandUpdate:
		s.display.SuccessUpdate()
	case CommandGet:
		s.display.SuccessGet()
	case CommandDelete:
		s.display.SuccessDelete()

	default:
		s.display.SuccessGeneric(commandName)
	}
}

// checkServerHealth проверяет доступность сервера.
// Выполняет команду health для проверки соединения.
//
// Возвращает:
//   - error: ошибка в случае недоступности сервера
func (s *Shell) checkServerHealth() error {
	_, err := s.commandRegistry.Execute(context.Background(), "health", nil)
	return err
}

// Добавляем приватные методы-обработчики

// handleExitCommand обрабатывает команду выхода из приложения.
//
// Возвращает:
//   - error: всегда nil
func (s *Shell) handleExitCommand() error {
	s.display.Goodbye()
	return nil
}

// handleLogoutCommand обрабатывает команду выхода из аккаунта.
// Переключает меню в основное состояние и очищает состояние действия.
func (s *Shell) handleLogoutCommand() {
	s.menuManager.SwitchToState(MenuStateMain)
	s.menuManager.ClearActionState()
	s.display.Logout()
}

// handleBackCommand обрабатывает команду возврата назад в меню.
// Показывает сообщение об ошибке, если возврат невозможен.
func (s *Shell) handleBackCommand() {
	if !s.menuManager.GoBack() {
		s.display.BackNotAllowed()
	}
}

// handleActionCommand обрабатывает команду действия (upload, show).
// Устанавливает состояние действия и переключает в меню выбора типа данных.
//
// Параметры:
//   - commandName: имя команды действия
func (s *Shell) handleActionCommand(commandName string) {
	s.menuManager.SetActionState(commandName, "")
	s.menuManager.SwitchToState(MenuStateDataType)
}

// handleDataTypeCommand обрабатывает выбор типа данных.
// Выполняет команду действия с выбранным типом данных.
//
// Параметры:
//   - commandName: выбранный тип данных
func (s *Shell) handleDataTypeCommand(commandName string) {
	actionState := s.menuManager.GetActionState()
	if actionState != nil {
		actionState.Type = commandName
		if result, err := s.executeCommand(actionState.Action); err != nil {
			// Проверяем, является ли ошибка ошибкой авторизации
			if isAuthError(err) {
				s.display.AuthError()
				s.handleLogoutCommand()
			} else {
				s.display.ErrorMsg(err)
			}
		} else {
			s.handleCommandResult(actionState.Action, result)
		}
		s.menuManager.ClearActionState()
		s.menuManager.SwitchToState(MenuStateAuthenticated)
	}
}

// handleGetDeleteCommand обрабатывает команды получения и удаления секретов.
//
// Параметры:
//   - commandName: имя команды (get или delete)
func (s *Shell) handleGetDeleteCommand(commandName string) {
	if result, err := s.executeCommand(commandName); err != nil {
		// Проверяем, является ли ошибка ошибкой авторизации
		if isAuthError(err) {
			s.display.AuthError()
			s.handleLogoutCommand()
		} else {
			s.display.ErrorMsg(err)
		}
	} else {
		s.handleCommandResult(commandName, result)
	}
}

// handleVersionCommand обрабатывает команду отображения версии.
// Получает информацию о версии и отображает ее пользователю.
func (s *Shell) handleVersionCommand() {
	if result, err := s.commandRegistry.Execute(context.Background(), CommandVersion, nil); err != nil {
		// Проверяем, является ли ошибка ошибкой авторизации
		if isAuthError(err) {
			s.display.AuthError()
			s.handleLogoutCommand()
		} else {
			s.display.ErrorMsg(err)
		}
	} else {
		if versionResult, ok := result.(*cmd.VersionResult); ok {
			s.display.DisplayVersion(versionResult.Version, versionResult.BuildTime)
		}
	}
}

// handleUpdateCommand обрабатывает команду обновления секрета.
// Получает текущие данные секрета, запрашивает новые данные и выполняет обновление.
func (s *Shell) handleUpdateCommand() {
	// Получаем название секрета через input
	secretName, err := s.inputHandler.GetSecretName()
	if err != nil {
		s.display.ErrorMsg(fmt.Errorf("ошибка получения названия секрета: %w", err))
		return
	}

	// Получаем данные секрета через команду show
	showResult, err := s.commandRegistry.Execute(context.Background(), CommandShow, secretName)
	if err != nil {
		// Проверяем, является ли ошибка ошибкой авторизации
		if isAuthError(err) {
			s.display.AuthError()
			s.handleLogoutCommand()
			return
		} else {
			s.display.ErrorMsg(fmt.Errorf("ошибка при получении данных секрета: %w", err))
			return
		}
	}

	showData, ok := showResult.(*cmd.ShowSecretResult)
	if !ok {
		s.display.ErrorMsg(fmt.Errorf("неверный формат данных секрета"))
		return
	}

	// Отображаем информацию о секрете через display
	s.display.DisplaySecretInfo(secretName, showData.Type, showData.Metadata, showData.Version)

	// Запрашиваем ввод новых данных через input
	s.inputHandler.PromptEnterNewData()

	// Получаем обновленные данные в зависимости от типа секрета
	var updatedData any
	switch showData.Type {
	case models.SecretTypeText:
		updatedData, err = s.inputHandler.GetUpdatedTextData(secretName)
	case models.SecretTypeLoginPassword:
		updatedData, err = s.inputHandler.GetUpdatedLoginPasswordData(secretName)
	case models.SecretTypeCard:
		updatedData, err = s.inputHandler.GetUpdatedCardData(secretName)
	case models.SecretTypeFile:
		updatedData, err = s.inputHandler.GetUpdatedFileData(secretName)
	default:
		s.display.ErrorMsg(fmt.Errorf("неподдерживаемый тип секрета: %s", showData.Type))
		return
	}

	if err != nil {
		s.display.ErrorMsg(fmt.Errorf("ошибка при получении обновленных данных: %w", err))
		return
	}

	// Создаем данные для обновления
	updateData := &cmd.UpdateData{
		SecretName: secretName,
		Version:    showData.Version,
		Type:       showData.Type,
		Data:       updatedData,
	}

	// Выполняем обновление
	if result, err := s.commandRegistry.Execute(context.Background(), CommandUpdate, updateData); err != nil {
		// Проверяем, является ли ошибка ошибкой авторизации
		if isAuthError(err) {
			s.display.AuthError()
			s.handleLogoutCommand()
		} else {
			s.display.ErrorMsg(err)
		}
	} else {
		s.handleCommandResult(CommandUpdate, result)
	}
}
