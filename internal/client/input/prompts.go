package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Базовые промпты для ввода данных
func promptLogin() {
	fmt.Print("Введите логин: ")
}

func promptPassword() {
	fmt.Print("Введите пароль: ")
}

func promptSecretName() {
	fmt.Print("Введите название секрета: ")
}

func promptTextSecret() {
	fmt.Print("Введите текст секрета: ")
}

func promptCardNumber() {
	fmt.Print("Введите номер карты: ")
}

func promptCardHolder() {
	fmt.Print("Введите имя владельца: ")
}

func promptCardExpiry() {
	fmt.Print("Введите срок действия (MM/YY): ")
}

func promptCardCVV() {
	fmt.Print("Введите CVV: ")
}

func promptFilePath() {
	fmt.Print("Введите путь к файлу: ")
}

func promptURL() {
	fmt.Print("Введите URL (опционально): ")
}

func promptMetadata() {
	fmt.Print("Введите метаданные (опционально): ")
}

// Промпты для обновления данных
func promptUpdateTextSecret() {
	fmt.Print("Введите новый текст секрета: ")
}

func promptUpdateLogin() {
	fmt.Print("Введите новый логин: ")
}

func promptUpdatePassword() {
	fmt.Print("Введите новый пароль: ")
}

func promptUpdateURL() {
	fmt.Print("Введите новый URL (опционально): ")
}

func promptUpdateCardNumber() {
	fmt.Print("Введите новый номер карты: ")
}

func promptUpdateCardHolder() {
	fmt.Print("Введите новое имя владельца: ")
}

func promptUpdateCardExpiry() {
	fmt.Print("Введите новый срок действия (MM/YY): ")
}

func promptUpdateCardCVV() {
	fmt.Print("Введите новый CVV: ")
}

func promptUpdateFilePath() {
	fmt.Print("Введите новый путь к файлу: ")
}

func promptEnterNewData() {
	fmt.Println("Введите новые данные:")
}

func readLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)
	input = strings.TrimSuffix(input, "\n")
	input = strings.TrimSuffix(input, "\r")
	return input, nil
}
