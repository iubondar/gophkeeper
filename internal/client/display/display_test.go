package display

import (
	"bytes"
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
				DisplayTextSecret(tt.secret, tt.metadata)
			})

			if output != tt.want {
				t.Errorf("DisplayTextSecret() = %q, want %q", output, tt.want)
			}
		})
	}
}

func TestDisplayLoginPassword(t *testing.T) {
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
				DisplayLoginPassword(tt.secret, tt.metadata)
			})

			if output != tt.want {
				t.Errorf("DisplayLoginPassword() = %q, want %q", output, tt.want)
			}
		})
	}
}

func TestDisplayCardData(t *testing.T) {
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
				DisplayCardData(tt.secret, tt.metadata)
			})

			if output != tt.want {
				t.Errorf("DisplayCardData() = %q, want %q", output, tt.want)
			}
		})
	}
}

func TestDisplayFileData(t *testing.T) {
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
				DisplayFileData(tt.secret, tt.metadata)
			})

			if output != tt.want {
				t.Errorf("DisplayFileData() = %q, want %q", output, tt.want)
			}
		})
	}
}

func TestDisplaySecretInfo(t *testing.T) {
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
				DisplaySecretInfo(tt.secretName, tt.secretType, tt.metadata, tt.version)
			})

			if output != tt.want {
				t.Errorf("DisplaySecretInfo() = %q, want %q", output, tt.want)
			}
		})
	}
}

// Benchmark тесты для проверки производительности
func BenchmarkDisplayTextSecret(b *testing.B) {
	secret := &models.TextSecretData{
		Name: "Тестовый секрет",
		Text: "Тестовый текст",
	}
	metadata := "Тестовые метаданные"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		output := captureOutput(func() {
			DisplayTextSecret(secret, metadata)
		})
		_ = output // Используем output чтобы избежать оптимизации
	}
}

func BenchmarkDisplayLoginPassword(b *testing.B) {
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
			DisplayLoginPassword(secret, metadata)
		})
		_ = output
	}
}

func BenchmarkDisplayCardData(b *testing.B) {
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
			DisplayCardData(secret, metadata)
		})
		_ = output
	}
}

func BenchmarkDisplayFileData(b *testing.B) {
	secret := &models.FileData{
		Name:     "Тестовый файл",
		FilePath: "test.pdf",
	}
	metadata := "Тестовые метаданные"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		output := captureOutput(func() {
			DisplayFileData(secret, metadata)
		})
		_ = output
	}
}

func BenchmarkDisplaySecretInfo(b *testing.B) {
	secretName := "Тестовый секрет"
	secretType := "text"
	metadata := "Тестовые метаданные"
	version := 1

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		output := captureOutput(func() {
			DisplaySecretInfo(secretName, secretType, metadata, version)
		})
		_ = output
	}
}
