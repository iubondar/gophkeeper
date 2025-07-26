package validators

import (
	"fmt"
	"strconv"

	"gophkeeper/internal/models"

	"github.com/asaskevich/govalidator"
)

// ValidateUserCredentials проверяет данные пользователя
func ValidateUserCredentials(credentials *models.UserCredentials) (bool, error) {
	if credentials == nil {
		return false, fmt.Errorf("данные пользователя не могут быть пустыми")
	}
	if ok, err := govalidator.ValidateStruct(credentials); !ok {
		if err != nil {
			return false, err
		}
		return false, fmt.Errorf("некорректные данные пользователя")
	}
	return true, nil
}

// ValidateTextSecretData проверяет данные текстового секрета
func ValidateTextSecretData(data *models.TextSecretData) (bool, error) {
	if data == nil {
		return false, fmt.Errorf("данные секрета не могут быть пустыми")
	}
	if ok, err := govalidator.ValidateStruct(data); !ok {
		if err != nil {
			return false, err
		}
		return false, fmt.Errorf("некорректные данные текстового секрета")
	}
	return true, nil
}

// ValidateLoginPasswordData проверяет данные логина и пароля
func ValidateLoginPasswordData(data *models.LoginPasswordData) (bool, error) {
	if data == nil {
		return false, fmt.Errorf("данные секрета не могут быть пустыми")
	}
	if ok, err := govalidator.ValidateStruct(data); !ok {
		if err != nil {
			return false, err
		}
		return false, fmt.Errorf("некорректные данные логина и пароля")
	}
	return true, nil
}

// isValidExpiry проверяет формат срока действия карты MM/YY и диапазон месяца
func isValidExpiry(expiry string) bool {
	if len(expiry) != 5 || expiry[2] != '/' {
		return false
	}
	month := expiry[:2]
	year := expiry[3:]
	if !govalidator.IsNumeric(month) || !govalidator.IsNumeric(year) {
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

// ValidateCardData проверяет данные банковской карты
func ValidateCardData(data *models.CardData) (bool, error) {
	if data == nil {
		return false, fmt.Errorf("данные секрета не могут быть пустыми")
	}
	if ok, err := govalidator.ValidateStruct(data); !ok {
		if err != nil {
			return false, err
		}
		return false, fmt.Errorf("некорректные данные карты")
	}
	if !govalidator.IsCreditCard(data.Number) {
		return false, fmt.Errorf("номер карты некорректен")
	}
	if !isValidExpiry(data.Expiry) {
		return false, fmt.Errorf("срок действия должен быть в формате MM/YY и месяц от 01 до 12")
	}
	if (len(data.CVV) != 3 && len(data.CVV) != 4) || !govalidator.IsNumeric(data.CVV) {
		return false, fmt.Errorf("CVV должен состоять из 3 или 4 цифр")
	}
	return true, nil
}

// ValidateFileData проверяет данные файла
func ValidateFileData(data *models.FileData) (bool, error) {
	if data == nil {
		return false, fmt.Errorf("данные секрета не могут быть пустыми")
	}
	if ok, err := govalidator.ValidateStruct(data); !ok {
		if err != nil {
			return false, err
		}
		return false, fmt.Errorf("некорректные данные файла")
	}
	return true, nil
}

// ValidateSecretName проверяет название секрета
func ValidateSecretName(name string) (bool, error) {
	if govalidator.IsNull(name) {
		return false, fmt.Errorf("название секрета не может быть пустым")
	}

	if !govalidator.IsAlphanumeric(name) {
		return false, fmt.Errorf("название секрета должно содержать только буквы и цифры")
	}

	if len(name) < 3 {
		return false, fmt.Errorf("название секрета должно содержать не менее 3 символов")
	}

	if len(name) > 255 {
		return false, fmt.Errorf("название секрета не может быть длиннее 255 символов")
	}
	return true, nil
}
