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

func TestGenerateUniqueFilePath(t *testing.T) {
	// Создаем временную директорию для тестов
	tempDir, err := os.MkdirTemp("", "test_unique_file")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name           string
		existingFiles  []string
		requestedPath  string
		expectedResult string
	}{
		{
			name:           "файл не существует",
			existingFiles:  []string{},
			requestedPath:  filepath.Join(tempDir, "test.txt"),
			expectedResult: filepath.Join(tempDir, "test.txt"),
		},
		{
			name:           "файл существует, добавляем номер",
			existingFiles:  []string{"test.txt"},
			requestedPath:  filepath.Join(tempDir, "test.txt"),
			expectedResult: filepath.Join(tempDir, "test_1.txt"),
		},
		{
			name:           "несколько файлов с похожими именами",
			existingFiles:  []string{"test.txt", "test_1.txt"},
			requestedPath:  filepath.Join(tempDir, "test.txt"),
			expectedResult: filepath.Join(tempDir, "test_2.txt"),
		},
		{
			name:           "файл без расширения",
			existingFiles:  []string{"testfile"},
			requestedPath:  filepath.Join(tempDir, "testfile"),
			expectedResult: filepath.Join(tempDir, "testfile_1"),
		},
		{
			name:           "файл с точкой в имени",
			existingFiles:  []string{"test.file.txt"},
			requestedPath:  filepath.Join(tempDir, "test.file.txt"),
			expectedResult: filepath.Join(tempDir, "test.file_1.txt"),
		},
		{
			name:           "файл с несколькими точками",
			existingFiles:  []string{"test.backup.txt"},
			requestedPath:  filepath.Join(tempDir, "test.backup.txt"),
			expectedResult: filepath.Join(tempDir, "test.backup_1.txt"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем существующие файлы
			for _, fileName := range tt.existingFiles {
				filePath := filepath.Join(tempDir, fileName)
				err := os.WriteFile(filePath, []byte("test content"), 0644)
				require.NoError(t, err)
			}

			// Вызываем функцию
			result, err := GenerateUniqueFilePath(tt.requestedPath)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedResult, result)

			// Проверяем, что файл с таким именем не существует
			_, err = os.Stat(result)
			assert.True(t, os.IsNotExist(err), "Файл с результатом не должен существовать")
		})
	}
}

func TestGenerateUniqueFilePath_WithRealFiles(t *testing.T) {
	// Создаем временную директорию
	tempDir, err := os.MkdirTemp("", "test_real_files")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Создаем несколько файлов
	files := []string{"document.pdf", "document_1.pdf", "document_2.pdf"}
	for _, fileName := range files {
		filePath := filepath.Join(tempDir, fileName)
		err := os.WriteFile(filePath, []byte("content"), 0644)
		require.NoError(t, err)
	}

	// Пытаемся создать файл с именем "document.pdf"
	requestedPath := filepath.Join(tempDir, "document.pdf")
	result, err := GenerateUniqueFilePath(requestedPath)
	require.NoError(t, err)

	// Ожидаем "document_3.pdf"
	expectedPath := filepath.Join(tempDir, "document_3.pdf")
	assert.Equal(t, expectedPath, result)

	// Проверяем, что файл с таким именем действительно не существует
	_, err = os.Stat(result)
	assert.True(t, os.IsNotExist(err))
}
