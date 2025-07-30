// Package display предоставляет функции для отображения информации пользователю.
// Включает методы для отображения различных типов секретов, меню, сообщений
// об успехе и ошибках в консольном интерфейсе.
package display

import (
	"fmt"
	"gophkeeper/internal/models"
	"strings"
)

// Display представляет интерфейс для отображения информации пользователю.
// Содержит методы для отображения различных типов секретов и системных сообщений.
type Display struct{}

// NewDisplay создает новый экземпляр Display.
//
// Возвращает:
//   - *Display: новый экземпляр отображения
func NewDisplay() *Display {
	return &Display{}
}

// DisplayTextSecret отображает текстовый секрет в консоли.
//
// Параметры:
//   - secret: данные текстового секрета
//   - metadata: метаданные секрета
func (d *Display) DisplayTextSecret(secret *models.TextSecretData, metadata string) {
	fmt.Printf("📝 Текстовый секрет: %s\n", secret.Name)
	fmt.Printf("   Текст: %s\n", secret.Text)
	if metadata != "" {
		fmt.Printf("   Метаданные: %s\n", metadata)
	}
	fmt.Println()
}

// DisplayLoginPassword отображает секрет логин/пароль в консоли.
//
// Параметры:
//   - secret: данные логина и пароля
//   - metadata: метаданные секрета
func (d *Display) DisplayLoginPassword(secret *models.LoginPasswordData, metadata string) {
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

// DisplayCardData отображает данные банковской карты в консоли.
//
// Параметры:
//   - secret: данные банковской карты
//   - metadata: метаданные секрета
func (d *Display) DisplayCardData(secret *models.CardData, metadata string) {
	fmt.Printf("💳 Данные карты: %s\n", secret.Name)
	fmt.Printf("   Номер: %s\n", secret.Number)
	fmt.Printf("   Владелец: %s\n", secret.Holder)
	fmt.Printf("   Срок действия: %s\n", secret.Expiry)
	fmt.Printf("   CVV: %s\n", secret.CVV)
	if metadata != "" {
		fmt.Printf("   Метаданные: %s\n", metadata)
	}
	fmt.Println()
}

// DisplayFileData отображает данные файла в консоли.
//
// Параметры:
//   - secret: данные файла
//   - metadata: метаданные секрета
func (d *Display) DisplayFileData(secret *models.FileData, metadata string) {
	fmt.Printf("📁 Файл: %s\n", secret.Name)

	// Если FilePath содержит оригинальное имя файла (не полный путь), показываем его
	if secret.FilePath != "" && !strings.Contains(secret.FilePath, "/") && !strings.Contains(secret.FilePath, "\\") {
		fmt.Printf("   Имя файла: %s\n", secret.FilePath)
	} else if secret.FilePath != "" {
		fmt.Printf("   Путь: %s\n", secret.FilePath)
	}

	if metadata != "" {
		fmt.Printf("   Метаданные: %s\n", metadata)
	}
	fmt.Println()
}

// DisplaySecretInfo отображает общую информацию о секрете.
//
// Параметры:
//   - secretName: имя секрета
//   - secretType: тип секрета
//   - metadata: метаданные секрета
//   - version: версия секрета
func (d *Display) DisplaySecretInfo(secretName, secretType, metadata string, version int) {
	fmt.Printf("=== %s ===\n", secretName)
	fmt.Printf("Тип секрета: %s\n", secretType)
	fmt.Printf("Версия секрета: %d\n", version)
	if metadata != "" {
		fmt.Printf("Метаданные: %s\n", metadata)
	}
	fmt.Println()
}

// Welcome отображает приветственное сообщение при запуске приложения.
func (d *Display) Welcome() {
	fmt.Println("=== GophKeeper CLI ===")
	fmt.Println("Подключение к серверу...")
}

// ServerConnected отображает сообщение об успешном подключении к серверу.
func (d *Display) ServerConnected() {
	fmt.Println("✅ Успешно подключился к серверу!")
	fmt.Println()
}

// InvalidChoice отображает сообщение о неверном выборе в меню.
func (d *Display) InvalidChoice() {
	fmt.Println("Неверный выбор. Попробуйте снова.")
}

// ErrorMsg отображает сообщение об ошибке.
//
// Параметры:
//   - err: ошибка для отображения
func (d *Display) ErrorMsg(err error) {
	fmt.Printf("Ошибка: %v\n", err)
}

// AuthError отображает сообщение об ошибке авторизации.
func (d *Display) AuthError() {
	fmt.Println("❌ Ошибка авторизации. Необходимо войти в систему заново.")
}

// SwitchToUserMenuNotice отображает уведомление о переключении в меню пользователя.
func (d *Display) SwitchToUserMenuNotice() {
	fmt.Println("Переключение в меню пользователя...")
}

// Goodbye отображает прощальное сообщение при выходе из приложения.
func (d *Display) Goodbye() {
	fmt.Println("До свидания!")
}

// Logout отображает сообщение об успешном выходе из аккаунта.
func (d *Display) Logout() {
	fmt.Println("✅ Вы вышли из аккаунта")
	fmt.Println()
}

// BackNotAllowed отображает сообщение о невозможности вернуться назад.
func (d *Display) BackNotAllowed() {
	fmt.Println("Нельзя вернуться назад")
}

// SuccessRegistration отображает сообщение об успешной регистрации.
func (d *Display) SuccessRegistration() {
	fmt.Println("✅ Регистрация успешна!")
}

// SuccessLogin отображает сообщение об успешном входе в систему.
func (d *Display) SuccessLogin() {
	fmt.Println("✅ Вход выполнен успешно!")
}

// SuccessUpload отображает сообщение об успешной загрузке секрета.
func (d *Display) SuccessUpload() {
	fmt.Println("✅ Секрет успешно загружен!")
}

// SuccessUpdate отображает сообщение об успешном обновлении секрета.
func (d *Display) SuccessUpdate() {
	fmt.Println("✅ Секрет успешно обновлен!")
}

// SuccessGet отображает сообщение об успешном получении секрета.
func (d *Display) SuccessGet() {
	fmt.Println("✅ Секрет успешно получен!")
}

// SuccessDelete отображает сообщение об успешном удалении секрета.
func (d *Display) SuccessDelete() {
	fmt.Println("✅ Секрет успешно удален!")
}

// SuccessGeneric отображает общее сообщение об успешном выполнении команды.
//
// Параметры:
//   - commandName: название выполненной команды
func (d *Display) SuccessGeneric(commandName string) {
	fmt.Printf("✅ %s выполнена успешно!\n", commandName)
}

// SuccessDownloadFile отображает сообщение об успешном скачивании файла.
//
// Параметры:
//   - filePath: путь к скачанному файлу
func (d *Display) SuccessDownloadFile(filePath string) {
	fmt.Printf("✅ Файл успешно скачан и сохранен: %s\n", filePath)
}

// DisplayVersion отображает информацию о версии приложения.
//
// Параметры:
//   - version: версия приложения
//   - buildTime: дата и время сборки
func (d *Display) DisplayVersion(version, buildTime string) {
	fmt.Println("=== Версия GophKeeper CLI ===")
	fmt.Printf("Версия: %s\n", version)
	fmt.Printf("Дата сборки: %s\n", buildTime)
	fmt.Println()
}

// MenuTitle отображает заголовок меню.
//
// Параметры:
//   - title: заголовок меню
func (d *Display) MenuTitle(title string) {
	fmt.Printf("=== %s ===\n", title)
}

// MenuItem отображает элемент меню с номером, названием и описанием.
//
// Параметры:
//   - id: номер пункта меню
//   - title: название пункта меню
//   - description: описание пункта меню
func (d *Display) MenuItem(id, title, description string) {
	fmt.Printf("%s. %s - %s\n", id, title, description)
}

// MenuChoice отображает приглашение к выбору пункта меню.
func (d *Display) MenuChoice() {
	fmt.Print("Выберите действие: ")
}

// MenuError отображает ошибку в меню.
//
// Параметры:
//   - message: сообщение об ошибке
func (d *Display) MenuError(message string) {
	fmt.Println(message)
}
