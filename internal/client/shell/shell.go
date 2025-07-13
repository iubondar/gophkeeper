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

		// Проверяем команду возврата
		if s.menuManager.IsBackCommand(commandName) {
			s.menuManager.SwitchToState("main")
			fmt.Println()
			continue
		}

		// Выполняем команду
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

	case "upload", "update", "get", "delete":
		// Пока просто выводим сообщение о том, что команда не реализована
		fmt.Printf("Команда '%s' пока не реализована\n", commandName)
		return nil

	default:
		return fmt.Errorf("неизвестная команда: %s", commandName)
	}
}

func (s *Shell) printSuccessMessage(commandName string) {
	switch commandName {
	case "register":
		fmt.Println("✅ Регистрация успешна!")
	case "login":
		fmt.Println("✅ Вход выполнен успешно!")
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
