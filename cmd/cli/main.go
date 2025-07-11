package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const serverURL = "http://localhost:8080"

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type SaltResponse struct {
	Salt string `json:"salt"`
}

func main() {
	fmt.Println("=== GophKeeper CLI ===")
	fmt.Println("Подключение к серверу...")

	// Проверяем доступность сервера
	if err := checkServerHealth(); err != nil {
		fmt.Printf("Ошибка подключения к серверу: %v\n", err)
		fmt.Println("Убедитесь, что сервер запущен на порту 8080")
		return
	}

	fmt.Println("✅ Успешно подключился к серверу!")
	fmt.Println()

	// Основной цикл меню
	for {
		showMainMenu()
		choice := getUserChoice()

		switch choice {
		case "1":
			handleRegistration()
		case "2":
			handleLogin()
		case "3":
			fmt.Println("До свидания!")
			return
		default:
			fmt.Println("Неверный выбор. Попробуйте снова.")
		}
		fmt.Println()
	}
}

func checkServerHealth() error {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(serverURL + "/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("сервер вернул статус %d", resp.StatusCode)
	}

	return nil
}

func showMainMenu() {
	fmt.Println("=== Главное меню ===")
	fmt.Println("1. Регистрация")
	fmt.Println("2. Вход")
	fmt.Println("3. Выход")
	fmt.Print("Выберите действие (1-3): ")
}

func getUserChoice() string {
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	choice = strings.TrimSuffix(choice, "\r") // Убираем carriage return для Windows
	return choice
}

func handleRegistration() {
	fmt.Println("\n=== Регистрация ===")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите логин: ")
	login, _ := reader.ReadString('\n')
	login = strings.TrimSpace(login)

	fmt.Print("Введите пароль: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	// Создаем запрос для регистрации
	requestBody := map[string]string{
		"login":         login,
		"password_hash": "test-hash", // В реальной реализации здесь будет хэш
		"salt":          "test-salt", // В реальной реализации здесь будет соль
	}

	jsonData, _ := json.Marshal(requestBody)

	resp, err := http.Post(serverURL+"/api/register", "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Printf("Ошибка при регистрации: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Ответ сервера: %s\n", string(body))
}

func handleLogin() {
	fmt.Println("\n=== Вход ===")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите логин: ")
	login, _ := reader.ReadString('\n')
	login = strings.TrimSpace(login)

	fmt.Print("Введите пароль: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	// Сначала получаем соль
	saltRequest := map[string]string{"login": login}
	jsonData, _ := json.Marshal(saltRequest)

	resp, err := http.Post(serverURL+"/api/login", "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Printf("Ошибка при получении соли: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var saltResp SaltResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &saltResp)

	fmt.Printf("Получена соль: %s\n", saltResp.Salt)

	// Теперь аутентифицируемся
	authRequest := map[string]string{
		"login":         login,
		"password_hash": "test-hash", // В реальной реализации здесь будет хэш
	}

	jsonData, _ = json.Marshal(authRequest)

	resp, err = http.Post(serverURL+"/api/authenticate", "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Printf("Ошибка при аутентификации: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var authResp AuthResponse
	body, _ = io.ReadAll(resp.Body)
	json.Unmarshal(body, &authResp)

	fmt.Printf("Успешная аутентификация!\n")
	fmt.Printf("Access Token: %s\n", authResp.AccessToken)
	fmt.Printf("Refresh Token: %s\n", authResp.RefreshToken)
	fmt.Printf("Время жизни токена: %d секунд\n", authResp.ExpiresIn)
}
