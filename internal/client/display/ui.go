package display

import (
	"fmt"
)

// UI функции для меню и интерфейса
func MenuError(msg string) {
	fmt.Println(msg)
}

func MenuTitle(title string) {
	fmt.Printf("=== %s ===\n", title)
}

func MenuItem(id, title, desc string) {
	fmt.Printf("%s. %s - %s\n", id, title, desc)
}

func MenuChoice() {
	fmt.Print("Выберите действие: ")
}

func Welcome() {
	fmt.Println("=== GophKeeper CLI ===")
	fmt.Println("Подключение к серверу...")
}

func ServerConnected() {
	fmt.Println("✅ Успешно подключился к серверу!")
	fmt.Println()
}

func Goodbye() {
	fmt.Println("До свидания!")
}

func Logout() {
	fmt.Println("✅ Вы вышли из аккаунта")
	fmt.Println()
}

func BackNotAllowed() {
	fmt.Println("Нельзя вернуться назад")
}

func InvalidChoice() {
	fmt.Println("Неверный выбор. Попробуйте снова.")
}

func SuccessRegistration() {
	fmt.Println("✅ Регистрация успешна!")
}

func SuccessLogin() {
	fmt.Println("✅ Вход выполнен успешно!")
}

func SuccessUpload() {
	fmt.Println("✅ Секрет успешно загружен!")
}

func SuccessUpdate() {
	fmt.Println("✅ Секрет успешно обновлен!")
}

func SuccessGet() {
	fmt.Println("✅ Секрет успешно получен!")
}

func SuccessDelete() {
	fmt.Println("✅ Секрет успешно удален!")
}

func SuccessGeneric(commandName string) {
	fmt.Printf("✅ %s выполнена успешно!\n", commandName)
}

func ErrorMsg(err error) {
	fmt.Printf("Ошибка: %v\n", err)
}

func SwitchToUserMenuNotice() {
	fmt.Println("Переключение в меню пользователя...")
}
