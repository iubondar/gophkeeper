package os_utils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeFilePath(t *testing.T) {
	// Получаем домашнюю директорию для сравнения
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "путь без тильды",
			input:    "/path/to/file",
			expected: "/path/to/file",
			wantErr:  false,
		},
		{
			name:     "путь с тильдой в середине",
			input:    "/path/~file",
			expected: "/path/~file",
			wantErr:  false,
		},
		{
			name:     "только тильда",
			input:    "~",
			expected: homeDir,
			wantErr:  false,
		},
		{
			name:     "тильда с путем",
			input:    "~/Documents/file.txt",
			expected: filepath.Join(homeDir, "Documents", "file.txt"),
			wantErr:  false,
		},
		{
			name:     "тильда с относительным путем",
			input:    "~/file.txt",
			expected: filepath.Join(homeDir, "file.txt"),
			wantErr:  false,
		},
		{
			name:     "тильда с именем пользователя (не поддерживается)",
			input:    "~username/file.txt",
			expected: "~username/file.txt",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NormalizeFilePath(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestNormalizeFilePath_CrossPlatform(t *testing.T) {
	// Тестируем на разных платформах
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	// Тестируем с тильдой
	normalizedPath, err := NormalizeFilePath("~/test.txt")
	require.NoError(t, err)
	expectedPath := filepath.Join(homeDir, "test.txt")
	assert.Equal(t, expectedPath, normalizedPath)

	// Проверяем, что путь использует правильные разделители для ОС
	assert.True(t, filepath.IsAbs(normalizedPath))
}

func TestGetDownloadsDir(t *testing.T) {
	// Тестируем получение папки загрузок
	downloadsDir, err := GetDownloadsDir()
	require.NoError(t, err)
	assert.NotEmpty(t, downloadsDir)

	// Проверяем, что папка существует
	stat, err := os.Stat(downloadsDir)
	require.NoError(t, err)
	assert.True(t, stat.IsDir())

	// Проверяем, что путь абсолютный
	assert.True(t, filepath.IsAbs(downloadsDir))
}

func TestGetDownloadsDir_PlatformSpecific(t *testing.T) {
	// Тестируем platform-specific логику
	downloadsDir, err := GetDownloadsDir()
	require.NoError(t, err)

	switch runtime.GOOS {
	case "windows":
		// На Windows проверяем, что путь содержит Downloads или OneDrive
		assert.True(t,
			filepath.Base(downloadsDir) == "Downloads" ||
				filepath.Base(filepath.Dir(downloadsDir)) == "OneDrive",
			"Windows downloads dir should be in Downloads or OneDrive/Downloads")

	case "darwin":
		// На macOS проверяем, что путь заканчивается на Downloads
		assert.Equal(t, "Downloads", filepath.Base(downloadsDir))

	case "linux":
		// На Linux проверяем, что путь заканчивается на Downloads или соответствует XDG
		base := filepath.Base(downloadsDir)
		assert.True(t,
			base == "Downloads" ||
				base == "downloads" ||
				os.Getenv("XDG_DOWNLOAD_DIR") != "",
			"Linux downloads dir should be Downloads or XDG_DOWNLOAD_DIR")
	}
}

func TestGetDownloadsDir_CreatesDirectory(t *testing.T) {
	// Создаем временную папку для теста
	tempDir, err := os.MkdirTemp("", "test_downloads")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Сохраняем оригинальную домашнюю папку
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Устанавливаем временную папку как домашнюю
	os.Setenv("HOME", tempDir)

	// Удаляем папку Downloads если она существует
	downloadsPath := filepath.Join(tempDir, "Downloads")
	os.RemoveAll(downloadsPath)

	// Вызываем функцию
	downloadsDir, err := GetDownloadsDir()
	require.NoError(t, err)

	// Проверяем, что папка была создана
	stat, err := os.Stat(downloadsDir)
	require.NoError(t, err)
	assert.True(t, stat.IsDir())
}
