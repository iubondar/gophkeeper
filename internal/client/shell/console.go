package shell

import (
	"bufio"
	"fmt"
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
	input = strings.TrimSuffix(input, "\n")
	input = strings.TrimSuffix(input, "\r")
	return input, nil
}

func switchToUserMenuNotice() {
	fmt.Println("Переключение в меню пользователя...")
}
