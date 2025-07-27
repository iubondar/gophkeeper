package os_utils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
