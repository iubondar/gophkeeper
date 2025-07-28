package display

import (
	"fmt"
	"gophkeeper/internal/models"
	"strings"
)

// Display представляет интерфейс для отображения информации
type Display struct{}

// NewDisplay создает новый экземпляр Display
func NewDisplay() *Display {
	return &Display{}
}

// DisplayTextSecret отображает текстовый секрет
func (d *Display) DisplayTextSecret(secret *models.TextSecretData, metadata string) {
	fmt.Printf("📝 Текстовый секрет: %s\n", secret.Name)
	fmt.Printf("   Текст: %s\n", secret.Text)
	if metadata != "" {
		fmt.Printf("   Метаданные: %s\n", metadata)
	}
	fmt.Println()
}

// DisplayLoginPassword отображает секрет логин/пароль
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

// DisplayCardData отображает данные банковской карты
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

// DisplayFileData отображает данные файла
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

// DisplaySecretInfo отображает информацию о секрете (версия, название, метаданные)
func (d *Display) DisplaySecretInfo(secretName, secretType, metadata string, version int) {
	fmt.Printf("=== %s ===\n", secretName)
	fmt.Printf("Тип секрета: %s\n", secretType)
	fmt.Printf("Версия секрета: %d\n", version)
	if metadata != "" {
		fmt.Printf("Метаданные: %s\n", metadata)
	}
	fmt.Println()
}

// Welcome отображает приветствие
func (d *Display) Welcome() {
	fmt.Println("=== GophKeeper CLI ===")
	fmt.Println("Подключение к серверу...")
}

// ServerConnected отображает сообщение об успешном подключении
func (d *Display) ServerConnected() {
	fmt.Println("✅ Успешно подключился к серверу!")
	fmt.Println()
}

// InvalidChoice отображает сообщение о неверном выборе
func (d *Display) InvalidChoice() {
	fmt.Println("Неверный выбор. Попробуйте снова.")
}

// ErrorMsg отображает сообщение об ошибке
func (d *Display) ErrorMsg(err error) {
	fmt.Printf("Ошибка: %v\n", err)
}

// SwitchToUserMenuNotice отображает уведомление о переключении в меню пользователя
func (d *Display) SwitchToUserMenuNotice() {
	fmt.Println("Переключение в меню пользователя...")
}

// Goodbye отображает прощание
func (d *Display) Goodbye() {
	fmt.Println("До свидания!")
}

// Logout отображает сообщение о выходе из аккаунта
func (d *Display) Logout() {
	fmt.Println("✅ Вы вышли из аккаунта")
	fmt.Println()
}

// BackNotAllowed отображает сообщение о невозможности вернуться назад
func (d *Display) BackNotAllowed() {
	fmt.Println("Нельзя вернуться назад")
}

// SuccessRegistration отображает сообщение об успешной регистрации
func (d *Display) SuccessRegistration() {
	fmt.Println("✅ Регистрация успешна!")
}

// SuccessLogin отображает сообщение об успешном входе
func (d *Display) SuccessLogin() {
	fmt.Println("✅ Вход выполнен успешно!")
}

// SuccessUpload отображает сообщение об успешной загрузке
func (d *Display) SuccessUpload() {
	fmt.Println("✅ Секрет успешно загружен!")
}

// SuccessUpdate отображает сообщение об успешном обновлении
func (d *Display) SuccessUpdate() {
	fmt.Println("✅ Секрет успешно обновлен!")
}

// SuccessGet отображает сообщение об успешном получении
func (d *Display) SuccessGet() {
	fmt.Println("✅ Секрет успешно получен!")
}

// SuccessDelete отображает сообщение об успешном удалении
func (d *Display) SuccessDelete() {
	fmt.Println("✅ Секрет успешно удален!")
}

// SuccessGeneric отображает общее сообщение об успехе
func (d *Display) SuccessGeneric(commandName string) {
	fmt.Printf("✅ %s выполнена успешно!\n", commandName)
}

// SuccessDownloadFile отображает сообщение об успешном скачивании файла
func (d *Display) SuccessDownloadFile(filePath string) {
	fmt.Printf("✅ Файл успешно скачан и сохранен: %s\n", filePath)
}

// DisplayVersion отображает информацию о версии
func (d *Display) DisplayVersion(version, buildTime string) {
	fmt.Println("=== Версия GophKeeper CLI ===")
	fmt.Printf("Версия: %s\n", version)
	fmt.Printf("Дата сборки: %s\n", buildTime)
	fmt.Println()
}

// MenuTitle отображает заголовок меню
func (d *Display) MenuTitle(title string) {
	fmt.Printf("=== %s ===\n", title)
}

// MenuItem отображает элемент меню
func (d *Display) MenuItem(id, title, description string) {
	fmt.Printf("%s. %s - %s\n", id, title, description)
}

// MenuChoice отображает приглашение к выбору
func (d *Display) MenuChoice() {
	fmt.Print("Выберите действие: ")
}

// MenuError отображает ошибку меню
func (d *Display) MenuError(message string) {
	fmt.Println(message)
}
