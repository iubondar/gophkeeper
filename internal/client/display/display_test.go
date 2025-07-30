package display

import (
	"bytes"
	"fmt"
	"gophkeeper/internal/models"
	"os"
	"testing"
)

// captureOutput захватывает вывод fmt.Printf для тестирования
func captureOutput(f func()) string {
	// Сохраняем оригинальный stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Выполняем функцию
	f()

	// Восстанавливаем stdout
	w.Close()
	os.Stdout = old

	// Читаем захваченный вывод
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestDisplayTextSecret(t *testing.T) {
	display := NewDisplay()
	tests := []struct {
		name     string
		secret   *models.TextSecretData
		metadata string
		want     string
	}{
		{
			name: "отображение текстового секрета без метаданных",
			secret: &models.TextSecretData{
				Name: "Мой секрет",
				Text: "Секретный текст",
			},
			metadata: "",
			want:     "📝 Текстовый секрет: Мой секрет\n   Текст: Секретный текст\n\n",
		},
		{
			name: "отображение текстового секрета с метаданными",
			secret: &models.TextSecretData{
				Name: "Важный секрет",
				Text: "Очень важный текст",
			},
			metadata: "Создан 2024-01-01",
			want:     "📝 Текстовый секрет: Важный секрет\n   Текст: Очень важный текст\n   Метаданные: Создан 2024-01-01\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				display.DisplayTextSecret(tt.secret, tt.metadata)
			})

			if output != tt.want {
				t.Errorf("DisplayTextSecret() = %q, want %q", output, tt.want)
			}
		})
	}
}

func TestDisplayLoginPassword(t *testing.T) {
	display := NewDisplay()
	tests := []struct {
		name     string
		secret   *models.LoginPasswordData
		metadata string
		want     string
	}{
		{
			name: "отображение логин/пароль без URL и метаданных",
			secret: &models.LoginPasswordData{
				Name:     "GitHub",
				Login:    "user@example.com",
				Password: "password123",
				URL:      "",
			},
			metadata: "",
			want:     "🔐 Логин/Пароль: GitHub\n   Логин: user@example.com\n   Пароль: password123\n\n",
		},
		{
			name: "отображение логин/пароль с URL и метаданными",
			secret: &models.LoginPasswordData{
				Name:     "Google",
				Login:    "user@gmail.com",
				Password: "securepass",
				URL:      "https://google.com",
			},
			metadata: "Двухфакторная аутентификация включена",
			want:     "🔐 Логин/Пароль: Google\n   Логин: user@gmail.com\n   Пароль: securepass\n   URL: https://google.com\n   Метаданные: Двухфакторная аутентификация включена\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				display.DisplayLoginPassword(tt.secret, tt.metadata)
			})

			if output != tt.want {
				t.Errorf("DisplayLoginPassword() = %q, want %q", output, tt.want)
			}
		})
	}
}

func TestDisplayCardData(t *testing.T) {
	display := NewDisplay()
	tests := []struct {
		name     string
		secret   *models.CardData
		metadata string
		want     string
	}{
		{
			name: "отображение данных карты без метаданных",
			secret: &models.CardData{
				Name:   "Основная карта",
				Number: "1234 5678 9012 3456",
				Holder: "IVAN IVANOV",
				Expiry: "12/25",
				CVV:    "123",
			},
			metadata: "",
			want:     "💳 Данные карты: Основная карта\n   Номер: 1234 5678 9012 3456\n   Владелец: IVAN IVANOV\n   Срок действия: 12/25\n   CVV: 123\n\n",
		},
		{
			name: "отображение данных карты с метаданными",
			secret: &models.CardData{
				Name:   "Кредитная карта",
				Number: "9876 5432 1098 7654",
				Holder: "PETR PETROV",
				Expiry: "06/26",
				CVV:    "456",
			},
			metadata: "Лимит: 100000 руб",
			want:     "💳 Данные карты: Кредитная карта\n   Номер: 9876 5432 1098 7654\n   Владелец: PETR PETROV\n   Срок действия: 06/26\n   CVV: 456\n   Метаданные: Лимит: 100000 руб\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				display.DisplayCardData(tt.secret, tt.metadata)
			})

			if output != tt.want {
				t.Errorf("DisplayCardData() = %q, want %q", output, tt.want)
			}
		})
	}
}

func TestDisplayFileData(t *testing.T) {
	display := NewDisplay()
	tests := []struct {
		name     string
		secret   *models.FileData
		metadata string
		want     string
	}{
		{
			name: "отображение файла с именем файла (не полный путь)",
			secret: &models.FileData{
				Name:     "Документ",
				FilePath: "document.pdf",
			},
			metadata: "",
			want:     "📁 Файл: Документ\n   Имя файла: document.pdf\n\n",
		},
		{
			name: "отображение файла с полным путем",
			secret: &models.FileData{
				Name:     "Важный документ",
				FilePath: "/home/user/documents/important.pdf",
			},
			metadata: "",
			want:     "📁 Файл: Важный документ\n   Путь: /home/user/documents/important.pdf\n\n",
		},
		{
			name: "отображение файла с Windows путем",
			secret: &models.FileData{
				Name:     "Windows файл",
				FilePath: "C:\\Users\\user\\Documents\\file.txt",
			},
			metadata: "",
			want:     "📁 Файл: Windows файл\n   Путь: C:\\Users\\user\\Documents\\file.txt\n\n",
		},
		{
			name: "отображение файла с метаданными",
			secret: &models.FileData{
				Name:     "Секретный файл",
				FilePath: "secret.txt",
			},
			metadata: "Размер: 1.5MB, Зашифрован",
			want:     "📁 Файл: Секретный файл\n   Имя файла: secret.txt\n   Метаданные: Размер: 1.5MB, Зашифрован\n\n",
		},
		{
			name: "отображение файла без пути",
			secret: &models.FileData{
				Name:     "Файл без пути",
				FilePath: "",
			},
			metadata: "",
			want:     "📁 Файл: Файл без пути\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				display.DisplayFileData(tt.secret, tt.metadata)
			})

			if output != tt.want {
				t.Errorf("DisplayFileData() = %q, want %q", output, tt.want)
			}
		})
	}
}

func TestDisplaySecretInfo(t *testing.T) {
	display := NewDisplay()
	tests := []struct {
		name       string
		secretName string
		secretType string
		metadata   string
		version    int
		want       string
	}{
		{
			name:       "отображение информации о секрете без метаданных",
			secretName: "Мой секрет",
			secretType: "text",
			metadata:   "",
			version:    1,
			want:       "=== Мой секрет ===\nТип секрета: text\nВерсия секрета: 1\n\n",
		},
		{
			name:       "отображение информации о секрете с метаданными",
			secretName: "Важный секрет",
			secretType: "login_password",
			metadata:   "Создан 2024-01-01, Обновлен 2024-01-15",
			version:    3,
			want:       "=== Важный секрет ===\nТип секрета: login_password\nВерсия секрета: 3\nМетаданные: Создан 2024-01-01, Обновлен 2024-01-15\n\n",
		},
		{
			name:       "отображение информации о файле",
			secretName: "Документ.pdf",
			secretType: "file",
			metadata:   "Размер: 2.5MB",
			version:    2,
			want:       "=== Документ.pdf ===\nТип секрета: file\nВерсия секрета: 2\nМетаданные: Размер: 2.5MB\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				display.DisplaySecretInfo(tt.secretName, tt.secretType, tt.metadata, tt.version)
			})

			if output != tt.want {
				t.Errorf("DisplaySecretInfo() = %q, want %q", output, tt.want)
			}
		})
	}
}

// Benchmark тесты для проверки производительности
func BenchmarkDisplayTextSecret(b *testing.B) {
	display := NewDisplay()
	secret := &models.TextSecretData{
		Name: "Тестовый секрет",
		Text: "Тестовый текст",
	}
	metadata := "Тестовые метаданные"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		output := captureOutput(func() {
			display.DisplayTextSecret(secret, metadata)
		})
		_ = output // Используем output чтобы избежать оптимизации
	}
}

func BenchmarkDisplayLoginPassword(b *testing.B) {
	display := NewDisplay()
	secret := &models.LoginPasswordData{
		Name:     "Тестовый аккаунт",
		Login:    "test@example.com",
		Password: "testpass",
		URL:      "https://example.com",
	}
	metadata := "Тестовые метаданные"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		output := captureOutput(func() {
			display.DisplayLoginPassword(secret, metadata)
		})
		_ = output
	}
}

func BenchmarkDisplayCardData(b *testing.B) {
	display := NewDisplay()
	secret := &models.CardData{
		Name:   "Тестовая карта",
		Number: "1234 5678 9012 3456",
		Holder: "TEST USER",
		Expiry: "12/25",
		CVV:    "123",
	}
	metadata := "Тестовые метаданные"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		output := captureOutput(func() {
			display.DisplayCardData(secret, metadata)
		})
		_ = output
	}
}

func BenchmarkDisplayFileData(b *testing.B) {
	display := NewDisplay()
	secret := &models.FileData{
		Name:     "Тестовый файл",
		FilePath: "test.pdf",
	}
	metadata := "Тестовые метаданные"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		output := captureOutput(func() {
			display.DisplayFileData(secret, metadata)
		})
		_ = output
	}
}

func BenchmarkDisplaySecretInfo(b *testing.B) {
	display := NewDisplay()
	secretName := "Тестовый секрет"
	secretType := "text"
	metadata := "Тестовые метаданные"
	version := 1

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		output := captureOutput(func() {
			display.DisplaySecretInfo(secretName, secretType, metadata, version)
		})
		_ = output
	}
}

// Тесты для методов структуры Display, которые не покрыты в ui_test.go
func TestDisplay_AuthError(t *testing.T) {
	display := NewDisplay()
	output := captureOutput(func() {
		display.AuthError()
	})
	expected := "❌ Ошибка авторизации. Необходимо войти в систему заново.\n"
	if output != expected {
		t.Errorf("AuthError() = %q, want %q", output, expected)
	}
}

func TestDisplay_SwitchToUserMenuNotice(t *testing.T) {
	display := NewDisplay()
	output := captureOutput(func() {
		display.SwitchToUserMenuNotice()
	})
	expected := "Переключение в меню пользователя...\n"
	if output != expected {
		t.Errorf("SwitchToUserMenuNotice() = %q, want %q", output, expected)
	}
}

// Тесты для остальных методов структуры Display
func TestDisplay_Welcome(t *testing.T) {
	display := NewDisplay()
	output := captureOutput(func() {
		display.Welcome()
	})
	expected := "=== GophKeeper CLI ===\nПодключение к серверу...\n"
	if output != expected {
		t.Errorf("Welcome() = %q, want %q", output, expected)
	}
}

func TestDisplay_ServerConnected(t *testing.T) {
	display := NewDisplay()
	output := captureOutput(func() {
		display.ServerConnected()
	})
	expected := "✅ Успешно подключился к серверу!\n\n"
	if output != expected {
		t.Errorf("ServerConnected() = %q, want %q", output, expected)
	}
}

func TestDisplay_InvalidChoice(t *testing.T) {
	display := NewDisplay()
	output := captureOutput(func() {
		display.InvalidChoice()
	})
	expected := "Неверный выбор. Попробуйте снова.\n"
	if output != expected {
		t.Errorf("InvalidChoice() = %q, want %q", output, expected)
	}
}

func TestDisplay_ErrorMsg(t *testing.T) {
	display := NewDisplay()
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "простая ошибка",
			err:      fmt.Errorf("test error"),
			expected: "Ошибка: test error\n",
		},
		{
			name:     "ошибка с деталями",
			err:      fmt.Errorf("connection failed: timeout"),
			expected: "Ошибка: connection failed: timeout\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				display.ErrorMsg(tt.err)
			})
			if output != tt.expected {
				t.Errorf("ErrorMsg() = %q, want %q", output, tt.expected)
			}
		})
	}
}

func TestDisplay_Goodbye(t *testing.T) {
	display := NewDisplay()
	output := captureOutput(func() {
		display.Goodbye()
	})
	expected := "До свидания!\n"
	if output != expected {
		t.Errorf("Goodbye() = %q, want %q", output, expected)
	}
}

func TestDisplay_Logout(t *testing.T) {
	display := NewDisplay()
	output := captureOutput(func() {
		display.Logout()
	})
	expected := "✅ Вы вышли из аккаунта\n\n"
	if output != expected {
		t.Errorf("Logout() = %q, want %q", output, expected)
	}
}

func TestDisplay_BackNotAllowed(t *testing.T) {
	display := NewDisplay()
	output := captureOutput(func() {
		display.BackNotAllowed()
	})
	expected := "Нельзя вернуться назад\n"
	if output != expected {
		t.Errorf("BackNotAllowed() = %q, want %q", output, expected)
	}
}

func TestDisplay_SuccessRegistration(t *testing.T) {
	display := NewDisplay()
	output := captureOutput(func() {
		display.SuccessRegistration()
	})
	expected := "✅ Регистрация успешна!\n"
	if output != expected {
		t.Errorf("SuccessRegistration() = %q, want %q", output, expected)
	}
}

func TestDisplay_SuccessLogin(t *testing.T) {
	display := NewDisplay()
	output := captureOutput(func() {
		display.SuccessLogin()
	})
	expected := "✅ Вход выполнен успешно!\n"
	if output != expected {
		t.Errorf("SuccessLogin() = %q, want %q", output, expected)
	}
}

func TestDisplay_SuccessUpload(t *testing.T) {
	display := NewDisplay()
	output := captureOutput(func() {
		display.SuccessUpload()
	})
	expected := "✅ Секрет успешно загружен!\n"
	if output != expected {
		t.Errorf("SuccessUpload() = %q, want %q", output, expected)
	}
}

func TestDisplay_SuccessUpdate(t *testing.T) {
	display := NewDisplay()
	output := captureOutput(func() {
		display.SuccessUpdate()
	})
	expected := "✅ Секрет успешно обновлен!\n"
	if output != expected {
		t.Errorf("SuccessUpdate() = %q, want %q", output, expected)
	}
}

func TestDisplay_SuccessGet(t *testing.T) {
	display := NewDisplay()
	output := captureOutput(func() {
		display.SuccessGet()
	})
	expected := "✅ Секрет успешно получен!\n"
	if output != expected {
		t.Errorf("SuccessGet() = %q, want %q", output, expected)
	}
}

func TestDisplay_SuccessDelete(t *testing.T) {
	display := NewDisplay()
	output := captureOutput(func() {
		display.SuccessDelete()
	})
	expected := "✅ Секрет успешно удален!\n"
	if output != expected {
		t.Errorf("SuccessDelete() = %q, want %q", output, expected)
	}
}

func TestDisplay_SuccessGeneric(t *testing.T) {
	display := NewDisplay()
	tests := []struct {
		name        string
		commandName string
		expected    string
	}{
		{
			name:        "общая команда",
			commandName: "upload",
			expected:    "✅ upload выполнена успешно!\n",
		},
		{
			name:        "команда с пробелами",
			commandName: "download file",
			expected:    "✅ download file выполнена успешно!\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				display.SuccessGeneric(tt.commandName)
			})
			if output != tt.expected {
				t.Errorf("SuccessGeneric() = %q, want %q", output, tt.expected)
			}
		})
	}
}

func TestDisplay_SuccessDownloadFile(t *testing.T) {
	display := NewDisplay()
	tests := []struct {
		name     string
		filePath string
		expected string
	}{
		{
			name:     "простой путь",
			filePath: "document.pdf",
			expected: "✅ Файл успешно скачан и сохранен: document.pdf\n",
		},
		{
			name:     "полный путь",
			filePath: "/home/user/downloads/document.pdf",
			expected: "✅ Файл успешно скачан и сохранен: /home/user/downloads/document.pdf\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				display.SuccessDownloadFile(tt.filePath)
			})
			if output != tt.expected {
				t.Errorf("SuccessDownloadFile() = %q, want %q", output, tt.expected)
			}
		})
	}
}

func TestDisplay_DisplayVersion(t *testing.T) {
	display := NewDisplay()
	tests := []struct {
		name      string
		version   string
		buildTime string
		expected  string
	}{
		{
			name:      "версия с датой сборки",
			version:   "1.0.0",
			buildTime: "2024-01-01 12:00:00",
			expected:  "=== Версия GophKeeper CLI ===\nВерсия: 1.0.0\nДата сборки: 2024-01-01 12:00:00\n\n",
		},
		{
			name:      "версия без даты сборки",
			version:   "2.1.0",
			buildTime: "",
			expected:  "=== Версия GophKeeper CLI ===\nВерсия: 2.1.0\nДата сборки: \n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				display.DisplayVersion(tt.version, tt.buildTime)
			})
			if output != tt.expected {
				t.Errorf("DisplayVersion() = %q, want %q", output, tt.expected)
			}
		})
	}
}

func TestDisplay_MenuTitle(t *testing.T) {
	display := NewDisplay()
	tests := []struct {
		name     string
		title    string
		expected string
	}{
		{
			name:     "простой заголовок",
			title:    "Главное меню",
			expected: "=== Главное меню ===\n",
		},
		{
			name:     "заголовок с символами",
			title:    "Меню пользователя & Настройки",
			expected: "=== Меню пользователя & Настройки ===\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				display.MenuTitle(tt.title)
			})
			if output != tt.expected {
				t.Errorf("MenuTitle() = %q, want %q", output, tt.expected)
			}
		})
	}
}

func TestDisplay_MenuItem(t *testing.T) {
	display := NewDisplay()
	tests := []struct {
		name        string
		id          string
		title       string
		description string
		expected    string
	}{
		{
			name:        "простой элемент меню",
			id:          "1",
			title:       "Загрузить секрет",
			description: "Загрузить новый секрет на сервер",
			expected:    "1. Загрузить секрет - Загрузить новый секрет на сервер\n",
		},
		{
			name:        "элемент с цифрами",
			id:          "2",
			title:       "Получить секрет",
			description: "Скачать секрет с сервера",
			expected:    "2. Получить секрет - Скачать секрет с сервера\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				display.MenuItem(tt.id, tt.title, tt.description)
			})
			if output != tt.expected {
				t.Errorf("MenuItem() = %q, want %q", output, tt.expected)
			}
		})
	}
}

func TestDisplay_MenuChoice(t *testing.T) {
	display := NewDisplay()
	output := captureOutput(func() {
		display.MenuChoice()
	})
	expected := "Выберите действие: "
	if output != expected {
		t.Errorf("MenuChoice() = %q, want %q", output, expected)
	}
}

func TestDisplay_MenuError(t *testing.T) {
	display := NewDisplay()
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "простая ошибка",
			message:  "Неверный выбор",
			expected: "Неверный выбор\n",
		},
		{
			name:     "ошибка с деталями",
			message:  "Ошибка соединения: timeout",
			expected: "Ошибка соединения: timeout\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				display.MenuError(tt.message)
			})
			if output != tt.expected {
				t.Errorf("MenuError() = %q, want %q", output, tt.expected)
			}
		})
	}
}
