// Package os_utils предоставляет platform-agnostic функции для работы с системными директориями.
// Основная цель - обеспечить кроссплатформенную совместимость при работе с файловой системой.
package os_utils

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetDownloadsDir возвращает путь к стандартной папке загрузок пользователя
// Поддерживает Windows, macOS и Linux
func GetDownloadsDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	var downloadsDir string

	switch runtime.GOOS {
	case "windows":
		// На Windows папка загрузок может быть в разных местах
		// Сначала проверяем стандартное место
		downloadsDir = filepath.Join(homeDir, "Downloads")

		// Если стандартная папка не существует, проверяем альтернативные места
		if _, err := os.Stat(downloadsDir); os.IsNotExist(err) {
			// Проверяем OneDrive Downloads (если используется)
			oneDriveDownloads := filepath.Join(homeDir, "OneDrive", "Downloads")
			if _, err := os.Stat(oneDriveDownloads); err == nil {
				downloadsDir = oneDriveDownloads
			}
		}

	case "darwin":
		// На macOS папка загрузок обычно в ~/Downloads
		downloadsDir = filepath.Join(homeDir, "Downloads")

	case "linux":
		// На Linux проверяем переменную окружения XDG_DOWNLOAD_DIR
		if xdgDir := os.Getenv("XDG_DOWNLOAD_DIR"); xdgDir != "" {
			downloadsDir = xdgDir
		} else {
			// По умолчанию используем ~/Downloads
			downloadsDir = filepath.Join(homeDir, "Downloads")
		}

	default:
		// Для других ОС используем стандартный подход
		downloadsDir = filepath.Join(homeDir, "Downloads")
	}

	// Создаем папку, если она не существует
	if _, err := os.Stat(downloadsDir); os.IsNotExist(err) {
		err = os.MkdirAll(downloadsDir, 0755)
		if err != nil {
			return "", err
		}
	}

	return downloadsDir, nil
}
