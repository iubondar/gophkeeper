package display

import (
	"fmt"
	"testing"
)

// Тесты для UI функций
func TestMenuError(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want string
	}{
		{
			name: "отображение ошибки меню",
			msg:  "Произошла ошибка",
			want: "Произошла ошибка\n",
		},
		{
			name: "отображение пустой ошибки",
			msg:  "",
			want: "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				MenuError(tt.msg)
			})

			if output != tt.want {
				t.Errorf("MenuError() = %q, want %q", output, tt.want)
			}
		})
	}
}

func TestMenuTitle(t *testing.T) {
	tests := []struct {
		name  string
		title string
		want  string
	}{
		{
			name:  "отображение заголовка меню",
			title: "Главное меню",
			want:  "=== Главное меню ===\n",
		},
		{
			name:  "отображение пустого заголовка",
			title: "",
			want:  "===  ===\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				MenuTitle(tt.title)
			})

			if output != tt.want {
				t.Errorf("MenuTitle() = %q, want %q", output, tt.want)
			}
		})
	}
}

func TestMenuItem(t *testing.T) {
	tests := []struct {
		name  string
		id    string
		title string
		desc  string
		want  string
	}{
		{
			name:  "отображение пункта меню",
			id:    "1",
			title: "Войти",
			desc:  "Войти в систему",
			want:  "1. Войти - Войти в систему\n",
		},
		{
			name:  "отображение пункта меню с пустым описанием",
			id:    "2",
			title: "Выход",
			desc:  "",
			want:  "2. Выход - \n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				MenuItem(tt.id, tt.title, tt.desc)
			})

			if output != tt.want {
				t.Errorf("MenuItem() = %q, want %q", output, tt.want)
			}
		})
	}
}

func TestMenuChoice(t *testing.T) {
	output := captureOutput(func() {
		MenuChoice()
	})

	want := "Выберите действие: "
	if output != want {
		t.Errorf("MenuChoice() = %q, want %q", output, want)
	}
}

func TestWelcome(t *testing.T) {
	output := captureOutput(func() {
		Welcome()
	})

	want := "=== GophKeeper CLI ===\nПодключение к серверу...\n"
	if output != want {
		t.Errorf("Welcome() = %q, want %q", output, want)
	}
}

func TestServerConnected(t *testing.T) {
	output := captureOutput(func() {
		ServerConnected()
	})

	want := "✅ Успешно подключился к серверу!\n\n"
	if output != want {
		t.Errorf("ServerConnected() = %q, want %q", output, want)
	}
}

func TestGoodbye(t *testing.T) {
	output := captureOutput(func() {
		Goodbye()
	})

	want := "До свидания!\n"
	if output != want {
		t.Errorf("Goodbye() = %q, want %q", output, want)
	}
}

func TestLogout(t *testing.T) {
	output := captureOutput(func() {
		Logout()
	})

	want := "✅ Вы вышли из аккаунта\n\n"
	if output != want {
		t.Errorf("Logout() = %q, want %q", output, want)
	}
}

func TestBackNotAllowed(t *testing.T) {
	output := captureOutput(func() {
		BackNotAllowed()
	})

	want := "Нельзя вернуться назад\n"
	if output != want {
		t.Errorf("BackNotAllowed() = %q, want %q", output, want)
	}
}

func TestInvalidChoice(t *testing.T) {
	output := captureOutput(func() {
		InvalidChoice()
	})

	want := "Неверный выбор. Попробуйте снова.\n"
	if output != want {
		t.Errorf("InvalidChoice() = %q, want %q", output, want)
	}
}

func TestSuccessRegistration(t *testing.T) {
	output := captureOutput(func() {
		SuccessRegistration()
	})

	want := "✅ Регистрация успешна!\n"
	if output != want {
		t.Errorf("SuccessRegistration() = %q, want %q", output, want)
	}
}

func TestSuccessLogin(t *testing.T) {
	output := captureOutput(func() {
		SuccessLogin()
	})

	want := "✅ Вход выполнен успешно!\n"
	if output != want {
		t.Errorf("SuccessLogin() = %q, want %q", output, want)
	}
}

func TestSuccessUpload(t *testing.T) {
	output := captureOutput(func() {
		SuccessUpload()
	})

	want := "✅ Секрет успешно загружен!\n"
	if output != want {
		t.Errorf("SuccessUpload() = %q, want %q", output, want)
	}
}

func TestSuccessUpdate(t *testing.T) {
	output := captureOutput(func() {
		SuccessUpdate()
	})

	want := "✅ Секрет успешно обновлен!\n"
	if output != want {
		t.Errorf("SuccessUpdate() = %q, want %q", output, want)
	}
}

func TestSuccessGet(t *testing.T) {
	output := captureOutput(func() {
		SuccessGet()
	})

	want := "✅ Секрет успешно получен!\n"
	if output != want {
		t.Errorf("SuccessGet() = %q, want %q", output, want)
	}
}

func TestSuccessDelete(t *testing.T) {
	output := captureOutput(func() {
		SuccessDelete()
	})

	want := "✅ Секрет успешно удален!\n"
	if output != want {
		t.Errorf("SuccessDelete() = %q, want %q", output, want)
	}
}

func TestSuccessDownloadFile(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     string
	}{
		{
			name:     "отображение успешной загрузки файла",
			filePath: "/path/to/file.pdf",
			want:     "✅ Файл успешно скачан и сохранен: /path/to/file.pdf\n",
		},
		{
			name:     "отображение успешной загрузки файла с пустым путем",
			filePath: "",
			want:     "✅ Файл успешно скачан и сохранен: \n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				SuccessDownloadFile(tt.filePath)
			})

			if output != tt.want {
				t.Errorf("SuccessDownloadFile() = %q, want %q", output, tt.want)
			}
		})
	}
}

func TestSuccessGeneric(t *testing.T) {
	tests := []struct {
		name        string
		commandName string
		want        string
	}{
		{
			name:        "отображение успешного выполнения команды",
			commandName: "upload",
			want:        "✅ upload выполнена успешно!\n",
		},
		{
			name:        "отображение успешного выполнения пустой команды",
			commandName: "",
			want:        "✅  выполнена успешно!\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				SuccessGeneric(tt.commandName)
			})

			if output != tt.want {
				t.Errorf("SuccessGeneric() = %q, want %q", output, tt.want)
			}
		})
	}
}

func TestErrorMsg(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "отображение ошибки",
			err:  fmt.Errorf("сетевая ошибка"),
			want: "Ошибка: сетевая ошибка\n",
		},
		{
			name: "отображение nil ошибки",
			err:  nil,
			want: "Ошибка: <nil>\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				ErrorMsg(tt.err)
			})

			if output != tt.want {
				t.Errorf("ErrorMsg() = %q, want %q", output, tt.want)
			}
		})
	}
}

func TestSwitchToUserMenuNotice(t *testing.T) {
	output := captureOutput(func() {
		SwitchToUserMenuNotice()
	})

	want := "Переключение в меню пользователя...\n"
	if output != want {
		t.Errorf("SwitchToUserMenuNotice() = %q, want %q", output, want)
	}
}

func TestDisplayVersion(t *testing.T) {
	tests := []struct {
		name      string
		version   string
		buildTime string
		want      string
	}{
		{
			name:      "отображение версии",
			version:   "1.0.0",
			buildTime: "2024-01-01T12:00:00Z",
			want:      "=== Версия GophKeeper CLI ===\nВерсия: 1.0.0\nДата сборки: 2024-01-01T12:00:00Z\n\n",
		},
		{
			name:      "отображение версии с пустыми значениями",
			version:   "",
			buildTime: "",
			want:      "=== Версия GophKeeper CLI ===\nВерсия: \nДата сборки: \n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				DisplayVersion(tt.version, tt.buildTime)
			})

			if output != tt.want {
				t.Errorf("DisplayVersion() = %q, want %q", output, tt.want)
			}
		})
	}
}
