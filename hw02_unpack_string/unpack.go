package hw02unpackstring

import (
	"errors"
	"strconv"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	runes := []rune(input)
	var result []rune

	for i := 0; i < len(runes); i++ {
		currentChar := runes[i]
		// Если мы в начале строки или предыдущий символ был цифрой,
		// а текущий символ - тоже цифра, это нарушение формата.
		if (i == 0 || unicode.IsDigit(runes[i-1])) && unicode.IsDigit(currentChar) {
			return "", ErrInvalidString
		}
		// --- Проверка на наличие цифр после текущего символа ---
		if i+1 < len(runes) && unicode.IsDigit(runes[i+1]) {
			// --- Проверка на двузначное число (условие ошибки) ---
			if i+2 < len(runes) && unicode.IsDigit(runes[i+2]) {
				return "", ErrInvalidString
			}

			// Парсим количество повторений (это гарантированно одна цифра)
			count, err := strconv.Atoi(string(runes[i+1]))
			if err != nil {
				return "", ErrInvalidString
			}

			// Добавляем символ count раз (если count=0, цикл не выполнится)
			for j := 0; j < count; j++ {
				result = append(result, currentChar)
			}

			// Пропускаем следующую позицию, так как там была цифра
			i++
		} else {
			// Если цифр после символа нет, добавляем его один раз
			result = append(result, currentChar)
		}
	}

	return string(result), nil
}
