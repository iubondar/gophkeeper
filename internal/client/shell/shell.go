package shell

import (
	"bufio"
	"context"
	"fmt"
	"gophkeeper/internal/client/cmd"
	"gophkeeper/internal/models"
	"os"
	"strings"
)

// Shell представляет интерактивный интерфейс для работы с GophKeeper
type Shell struct {
	commandProducer cmd.CommandProducer
}

// NewShell создает новый экземпляр Shell
func NewShell(producer cmd.CommandProducer) *Shell {
	return &Shell{commandProducer: producer}
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
		s.showMainMenu()
		choice := s.getUserChoice()

		switch choice {
		case "1":
			fmt.Println("\n=== Регистрация ===")

			// Получаем данные от пользователя
			credentials, err := s.getUserCredentials()
			if err != nil {
				return fmt.Errorf("ошибка получения данных пользователя: %w", err)
			}
			registerHandler := s.commandProducer.Register()
			err = (&registerHandler).Run(context.Background(), *credentials)
			if err != nil {
				fmt.Printf("Ошибка: %v\n", err)
			}
		case "2":
			s.handleLogin()
		case "3":
			fmt.Println("До свидания!")
			return nil
		default:
			fmt.Println("Неверный выбор. Попробуйте снова.")
		}
		fmt.Println()
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

// showMainMenu показывает главное меню
func (s *Shell) showMainMenu() {
	fmt.Println("=== Главное меню ===")
	fmt.Println("1. Регистрация")
	fmt.Println("2. Вход")
	fmt.Println("3. Выход")
	fmt.Print("Выберите действие (1-3): ")
}

// getUserChoice получает выбор пользователя
func (s *Shell) getUserChoice() string {
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	choice = strings.TrimSuffix(choice, "\r") // Убираем carriage return для Windows
	return choice
}

// getUserCredentials получает логин и пароль от пользователя
func (s *Shell) getUserCredentials() (credentials *models.UserCredentials, err error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите логин: ")
	login, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	login = strings.TrimSpace(login)

	fmt.Print("Введите пароль: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	password = strings.TrimSpace(password)

	// Валидация входных данных
	if login == "" {
		return nil, fmt.Errorf("логин не может быть пустым")
	}
	if password == "" {
		return nil, fmt.Errorf("пароль не может быть пустым")
	}

	return &models.UserCredentials{Login: login, Password: password}, nil
}

// handleLogin обрабатывает вход пользователя
func (s *Shell) handleLogin() {

}
