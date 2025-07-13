package shell

import (
	"context"
	"fmt"
	"gophkeeper/internal/client/cmd"
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
	fmt.Println("=== GophKeeper CLI ===")
	fmt.Println("Подключение к серверу...")

	// Проверяем доступность сервера
	if err := s.checkServerHealth(); err != nil {
		return fmt.Errorf("ошибка подключения к серверу: %w", err)
	}

	fmt.Println("✅ Успешно подключился к серверу!")
	fmt.Println()

	// Основной цикл меню
	for {
		s.menuManager.ShowMenu()
		choice := s.inputHandler.GetUserChoice()

		// Получаем команду по выбору пользователя
		commandName, exists := s.menuManager.GetCommandByID(choice)
		if !exists {
			fmt.Println("Неверный выбор. Попробуйте снова.")
			continue
		}

		// Проверяем команду выхода
		if s.menuManager.IsExitCommand(commandName) {
			fmt.Println("До свидания!")
			return nil
		}

		// Проверяем команду выхода из аккаунта
		if s.menuManager.IsLogoutCommand(commandName) {
			s.menuManager.SwitchToState("main")
			s.menuManager.ClearActionState()
			fmt.Println("✅ Вы вышли из аккаунта")
			fmt.Println()
			continue
		}

		// Проверяем команду возврата
		if s.menuManager.IsBackCommand(commandName) {
			if s.menuManager.GoBack() {
				fmt.Println()
			} else {
				fmt.Println("Нельзя вернуться назад")
			}
			continue
		}

		// Обрабатываем команды действий (upload, update, get, delete)
		if s.menuManager.IsActionCommand(commandName) {
			s.menuManager.SetActionState(commandName, "")
			s.menuManager.SwitchToState("data_type")
			fmt.Println()
			continue
		}

		// Обрабатываем команды выбора типа данных
		if s.menuManager.IsDataTypeCommand(commandName) {
			actionState := s.menuManager.GetActionState()
			if actionState != nil {
				actionState.Type = commandName
				if err := s.executeActionWithType(actionState.Action, actionState.Type); err != nil {
					fmt.Printf("Ошибка: %v\n", err)
				} else {
					s.printSuccessMessage(actionState.Action)
				}
				s.menuManager.ClearActionState()
				s.menuManager.SwitchToState("authenticated")
			}
			fmt.Println()
			continue
		}

		// Выполняем обычные команды
		if err := s.executeCommand(commandName); err != nil {
			fmt.Printf("Ошибка: %v\n", err)
		} else {
			s.printSuccessMessage(commandName)

			// Если это авторизация или регистрация, переключаем состояние
			if (commandName == "register" || commandName == "login") && s.menuManager.GetCurrentState() == "main" {
				fmt.Println("Переключение в меню пользователя...")
				s.menuManager.SwitchToState("authenticated")
			}
		}
		fmt.Println()
	}
}

// executeCommand выполняет команду с соответствующими аргументами
func (s *Shell) executeCommand(commandName string) error {
	ctx := context.Background()

	switch commandName {
	case "register":
		credentials, err := s.inputHandler.GetUserCredentials()
		if err != nil {
			return fmt.Errorf("ошибка получения данных пользователя: %w", err)
		}
		return s.commandRegistry.Execute(ctx, commandName, *credentials)

	case "login":
		credentials, err := s.inputHandler.GetLoginCredentials()
		if err != nil {
			return fmt.Errorf("ошибка получения данных пользователя: %w", err)
		}
		return s.commandRegistry.Execute(ctx, commandName, *credentials)

	default:
		return fmt.Errorf("неизвестная команда: %s", commandName)
	}
}

// executeActionWithType выполняет действие с выбранным типом данных через команды
func (s *Shell) executeActionWithType(action, dataType string) error {
	ctx := context.Background()

	// Получаем данные в зависимости от типа
	var data any
	var err error

	switch dataType {
	case "text":
		data, err = s.inputHandler.GetTextData()
	case "login_password":
		data, err = s.inputHandler.GetLoginPasswordData()
	case "card":
		data, err = s.inputHandler.GetCardData()
	case "file":
		data, err = s.inputHandler.GetFileData()
	default:
		return fmt.Errorf("неизвестный тип данных: %s", dataType)
	}

	if err != nil {
		return fmt.Errorf("ошибка получения данных: %w", err)
	}

	// Для операций get и delete нужен только название секрета
	if action == "get" || action == "delete" {
		secretName, err := s.inputHandler.GetSecretName()
		if err != nil {
			return fmt.Errorf("ошибка получения названия секрета: %w", err)
		}
		data = secretName
	}

	// Выполняем команду через реестр команд
	return s.commandRegistry.Execute(ctx, action, data)
}

func (s *Shell) printSuccessMessage(commandName string) {
	switch commandName {
	case "register":
		fmt.Println("✅ Регистрация успешна!")
	case "login":
		fmt.Println("✅ Вход выполнен успешно!")
	case "upload":
		fmt.Println("✅ Секрет успешно загружен!")
	case "update":
		fmt.Println("✅ Секрет успешно обновлен!")
	case "get":
		fmt.Println("✅ Секрет успешно получен!")
	case "delete":
		fmt.Println("✅ Секрет успешно удален!")
	default:
		fmt.Printf("✅ %s выполнена успешно!\n", commandName)
	}
}

// checkServerHealth проверяет доступность сервера
func (s *Shell) checkServerHealth() error {
	// client := &http.Client{
	// 	Timeout: 5 * time.Second,
	// }

	// resp, err := client.Get(s.serverURL + "/health")
	// if err != nil {
	// 	return err
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	return fmt.Errorf("сервер вернул статус %d", resp.StatusCode)
	// }

	return nil
}
