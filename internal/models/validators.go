package models

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// ValidateUserCredentials проверяет данные пользователя
func ValidateUserCredentials(credentials *UserCredentials) (bool, error) {
	if credentials == nil {
		return false, fmt.Errorf("данные пользователя не могут быть пустыми")
	}
	if credentials.Login == "" {
		return false, fmt.Errorf("логин не может быть пустым")
	}
	if credentials.Password == "" {
		return false, fmt.Errorf("пароль не может быть пустым")
	}
	return true, nil
}

// ValidateTextSecretData проверяет данные текстового секрета
func ValidateTextSecretData(data *TextSecretData) (bool, error) {
	if data == nil {
		return false, fmt.Errorf("данные секрета не могут быть пустыми")
	}
	if data.Name == "" {
		return false, fmt.Errorf("название секрета не может быть пустым")
	}
	if data.Text == "" {
		return false, fmt.Errorf("текст секрета не может быть пустым")
	}
	return true, nil
}

// ValidateLoginPasswordData проверяет данные логина и пароля
func ValidateLoginPasswordData(data *LoginPasswordData) (bool, error) {
	if data == nil {
		return false, fmt.Errorf("данные секрета не могут быть пустыми")
	}
	if data.Name == "" {
		return false, fmt.Errorf("название секрета не может быть пустым")
	}
	if data.Login == "" {
		return false, fmt.Errorf("логин не может быть пустым")
	}
	if data.Password == "" {
		return false, fmt.Errorf("пароль не может быть пустым")
	}
	return true, nil
}

// isDigitsOnly проверяет, что строка состоит только из цифр
func isDigitsOnly(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// isValidExpiry проверяет формат срока действия карты MM/YY и диапазон месяца
func isValidExpiry(expiry string) bool {
	if len(expiry) != 5 || expiry[2] != '/' {
		return false
	}
	month := expiry[:2]
	year := expiry[3:]
	if !isDigitsOnly(month) || !isDigitsOnly(year) {
		return false
	}
	m, err := strconv.Atoi(month)
	if err != nil || m < 1 || m > 12 {
		return false
	}
	if len(year) != 2 {
		return false
	}
	return true
}

// isValidEmail выполняет базовую проверку email (наличие @ и точки после @)
func isValidEmail(email string) bool {
	at := strings.Index(email, "@")
	if at < 1 || at == len(email)-1 {
		return false
	}
	dot := strings.LastIndex(email, ".")
	if dot < at+2 || dot == len(email)-1 {
		return false
	}
	return true
}

// ValidateCardData проверяет данные банковской карты
func ValidateCardData(data *CardData) (bool, error) {
	if data == nil {
		return false, fmt.Errorf("данные секрета не могут быть пустыми")
	}
	if data.Name == "" {
		return false, fmt.Errorf("название секрета не может быть пустым")
	}
	if data.Number == "" {
		return false, fmt.Errorf("номер карты не может быть пустым")
	}
	if len(data.Number) != 16 || !isDigitsOnly(data.Number) {
		return false, fmt.Errorf("номер карты должен состоять из 16 цифр")
	}
	if !ValidateLuhn(data.Number) {
		return false, fmt.Errorf("номер карты не прошёл проверку по алгоритму Луна")
	}
	if data.Holder == "" {
		return false, fmt.Errorf("имя владельца не может быть пустым")
	}
	if data.Expiry == "" {
		return false, fmt.Errorf("срок действия не может быть пустым")
	}
	if !isValidExpiry(data.Expiry) {
		return false, fmt.Errorf("срок действия должен быть в формате MM/YY и месяц от 01 до 12")
	}
	if data.CVV == "" {
		return false, fmt.Errorf("CVV не может быть пустым")
	}
	if (len(data.CVV) != 3 && len(data.CVV) != 4) || !isDigitsOnly(data.CVV) {
		return false, fmt.Errorf("CVV должен состоять из 3 или 4 цифр")
	}
	return true, nil
}

// ValidateFileData проверяет данные файла
func ValidateFileData(data *FileData) (bool, error) {
	if data == nil {
		return false, fmt.Errorf("данные секрета не могут быть пустыми")
	}
	if data.Name == "" {
		return false, fmt.Errorf("название секрета не может быть пустым")
	}
	if data.FilePath == "" {
		return false, fmt.Errorf("путь к файлу не может быть пустым")
	}
	return true, nil
}

// ValidateSecretName проверяет название секрета
func ValidateSecretName(name string) (bool, error) {
	if name == "" {
		return false, fmt.Errorf("название секрета не может быть пустым")
	}
	return true, nil
}

// ValidateLuhn проверяет, соответствует ли строка алгоритму Луна
// https://ru.wikipedia.org/wiki/%D0%90%D0%BB%D0%B3%D0%BE%D1%80%D0%B8%D1%82%D0%BC_%D0%9B%D1%83%D0%BD%D0%B0
func ValidateLuhn(input string) bool {
	// Проверяем, что строка состоит только из цифр и не пуста
	if len(input) == 0 {
		return false
	}

	for _, r := range input {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	if len(input) == 1 {
		// Однозначные числа считаются валидными только если они равны 0
		return input[0] == '0'
	}

	sum := 0
	isSecond := false // Флаг для обработки каждой второй цифры

	// Итерируемся по строке справа налево
	for i := len(input) - 1; i >= 0; i-- {
		// Вычитая код символа '0' из кода текущего символа, мы получаем числовое значение цифры.
		digit := int(input[i] - '0')

		if isSecond {
			digit *= 2
			if digit > 9 {
				digit = (digit / 10) + (digit % 10)
			}
		}

		sum += digit
		isSecond = !isSecond
	}

	return sum%10 == 0
}
