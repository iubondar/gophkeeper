package shell

import (
	"bufio"
	"fmt"
	"gophkeeper/internal/models"
	"os"
	"strings"
)

func promptLogin() {
	fmt.Print("Введите логин: ")
}

func promptPassword() {
	fmt.Print("Введите пароль: ")
}

func promptSecretName() {
	fmt.Print("Введите название секрета: ")
}

func promptTextSecret() {
	fmt.Print("Введите текст секрета: ")
}

func promptCardNumber() {
	fmt.Print("Введите номер карты: ")
}

func promptCardHolder() {
	fmt.Print("Введите имя владельца: ")
}

func promptCardExpiry() {
	fmt.Print("Введите срок действия (MM/YY): ")
}

func promptCardCVV() {
	fmt.Print("Введите CVV: ")
}

func promptFilePath() {
	fmt.Print("Введите путь к файлу: ")
}

func promptURL() {
	fmt.Print("Введите URL (опционально): ")
}

func promptMetadata() {
	fmt.Print("Введите метаданные (опционально): ")
}

func menuError(msg string) {
	fmt.Println(msg)
}

func menuTitle(title string) {
	fmt.Printf("=== %s ===\n", title)
}

func menuItem(id, title, desc string) {
	fmt.Printf("%s. %s - %s\n", id, title, desc)
}

func menuChoice() {
	fmt.Print("Выберите действие: ")
}

func welcome() {
	fmt.Println("=== GophKeeper CLI ===")
	fmt.Println("Подключение к серверу...")
}

func serverConnected() {
	fmt.Println("✅ Успешно подключился к серверу!")
	fmt.Println()
}

func goodbye() {
	fmt.Println("До свидания!")
}

func logout() {
	fmt.Println("✅ Вы вышли из аккаунта")
	fmt.Println()
}

func backNotAllowed() {
	fmt.Println("Нельзя вернуться назад")
}

func invalidChoice() {
	fmt.Println("Неверный выбор. Попробуйте снова.")
}

func successRegistration() {
	fmt.Println("✅ Регистрация успешна!")
}

func successLogin() {
	fmt.Println("✅ Вход выполнен успешно!")
}

func successUpload() {
	fmt.Println("✅ Секрет успешно загружен!")
}

func successUpdate() {
	fmt.Println("✅ Секрет успешно обновлен!")
}

func successGet() {
	fmt.Println("✅ Секрет успешно получен!")
}

func successDelete() {
	fmt.Println("✅ Секрет успешно удален!")
}

func successGeneric(commandName string) {
	fmt.Printf("✅ %s выполнена успешно!\n", commandName)
}

func errorMsg(err error) {
	fmt.Printf("Ошибка: %v\n", err)
}

func readLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)
	input = strings.TrimSuffix(input, "\n")
	input = strings.TrimSuffix(input, "\r")
	return input, nil
}

func switchToUserMenuNotice() {
	fmt.Println("Переключение в меню пользователя...")
}

// DisplayTextSecret отображает текстовый секрет
func DisplayTextSecret(secret *models.TextSecretData, metadata string) {
	fmt.Printf("📝 Текстовый секрет: %s\n", secret.Name)
	fmt.Printf("   Текст: %s\n", secret.Text)
	if metadata != "" {
		fmt.Printf("   Метаданные: %s\n", metadata)
	}
	fmt.Println()
}

// DisplayLoginPassword отображает секрет логин/пароль
func DisplayLoginPassword(secret *models.LoginPasswordData, metadata string) {
	fmt.Printf("🔐 Логин/Пароль: %s\n", secret.Name)
	fmt.Printf("   Логин: %s\n", secret.Login)
	fmt.Printf("   Пароль: %s\n", secret.Password)
	if secret.URL != "" {
		fmt.Printf("   URL: %s\n", secret.URL)
	}
	if metadata != "" {
		fmt.Printf("   Метаданные: %s\n", metadata)
	}
	fmt.Println()
}

// DisplayCardData отображает данные банковской карты
func DisplayCardData(secret *models.CardData, metadata string) {
	fmt.Printf("💳 Банковская карта: %s\n", secret.Name)
	fmt.Printf("   Номер: %s\n", secret.Number)
	fmt.Printf("   Держатель: %s\n", secret.Holder)
	fmt.Printf("   Срок действия: %s\n", secret.Expiry)
	fmt.Printf("   CVV: %s\n", secret.CVV)
	if metadata != "" {
		fmt.Printf("   Метаданные: %s\n", metadata)
	}
	fmt.Println()
}
