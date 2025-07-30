// Package os_utils предоставляет platform-agnostic функции для работы с системными директориями.
// Основная цель - обеспечить кроссплатформенную совместимость при работе с файловой системой.
package os_utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// NormalizeFilePath нормализует путь к файлу, выполняя необходимые преобразования.
// Включает расширение символа ~ в домашнюю директорию пользователя.
//
// Параметры:
//   - path: путь к файлу для нормализации
//
// Возвращает:
//   - string: нормализованный путь к файлу
//   - error: ошибка в случае неудачи
func NormalizeFilePath(path string) (string, error) {
	// Расширяем символ ~ в начале пути
	expandedPath, err := expandTilde(path)
	if err != nil {
		return "", err
	}

	// В будущем здесь могут быть добавлены другие преобразования
	// например, нормализация разделителей путей, обработка переменных окружения и т.д.

	return expandedPath, nil
}

// expandTilde расширяет символ ~ в начале пути на домашнюю директорию пользователя.
// Если путь не начинается с ~, возвращает путь без изменений.
//
// Параметры:
//   - path: путь для расширения
//
// Возвращает:
//   - string: расширенный путь
//   - error: ошибка в случае неудачи
func expandTilde(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	// Получаем домашнюю директорию пользователя
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Если путь равен "~", возвращаем домашнюю директорию
	if path == "~" {
		return homeDir, nil
	}

	// Если путь начинается с "~/", заменяем на домашнюю директорию
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:]), nil
	}

	// Для других случаев (например, "~username") возвращаем путь без изменений
	// так как это может быть ссылка на домашнюю директорию другого пользователя
	return path, nil
}

// GetDownloadsDir возвращает путь к стандартной папке загрузок пользователя.
// Поддерживает Windows, macOS и Linux.
// Создает папку загрузок, если она не существует.
//
// Возвращает:
//   - string: путь к папке загрузок
//   - error: ошибка в случае неудачи
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

// GenerateUniqueFilePath генерирует уникальный путь к файлу, добавляя номер к имени файла,
// если файл с таким именем уже существует.
// Например: file.txt -> file_1.txt -> file_2.txt
//
// Параметры:
//   - filePath: исходный путь к файлу
//
// Возвращает:
//   - string: уникальный путь к файлу
//   - error: ошибка в случае неудачи
func GenerateUniqueFilePath(filePath string) (string, error) {
	// Проверяем, существует ли файл с таким именем, и если да, добавляем номер
	counter := 1
	originalFilePath := filePath
	for {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			break
		}
		// Файл существует, добавляем номер
		ext := filepath.Ext(originalFilePath)
		nameWithoutExt := strings.TrimSuffix(originalFilePath, ext)
		filePath = strings.Join([]string{nameWithoutExt, "_", strconv.Itoa(counter), ext}, "")
		counter++
	}

	return filePath, nil
}
