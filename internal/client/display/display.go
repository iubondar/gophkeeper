package display

import (
	"fmt"
	"gophkeeper/internal/models"
)

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
func DisplayFileData(secret *models.FileData, metadata string) {
	fmt.Printf("📁 Файл: %s\n", secret.Name)
	fmt.Printf("   Путь: %s\n", secret.FilePath)
	if metadata != "" {
		fmt.Printf("   Метаданные: %s\n", metadata)
	}
	fmt.Println()
}

// DisplaySecretInfo отображает информацию о секрете (версия, название, метаданные)
func DisplaySecretInfo(secretName, secretType, metadata string, version int) {
	fmt.Printf("=== %s ===\n", secretName)
	fmt.Printf("Тип секрета: %s\n", secretType)
	fmt.Printf("Версия секрета: %d\n", version)
	if metadata != "" {
		fmt.Printf("Метаданные: %s\n", metadata)
	}
	fmt.Println()
}
