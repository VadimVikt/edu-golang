package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	result := make(Environment)
	listFile, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range listFile {
		if strings.Contains(file.Name(), "=") {
			result[file.Name()] = EnvValue{"", false}
			continue
		}
		envName, err := prepareEnv(dir + "/" + file.Name())
		if err != nil || envName == "" {
			result[file.Name()] = EnvValue{"", false}
			continue
		}
		result[file.Name()] = EnvValue{envName, true}
	}
	return result, nil
}

func prepareEnv(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		// Если файл пустой или произошла ошибка чтения
		if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) {
			return "", fmt.Errorf("ошибка при чтении первого байта из %s: %w", filePath, err)
		}
		// Файл пуст, возвращаем пустую строку без ошибки
		return "", nil
	}
	value := scanner.Text()
	value = strings.TrimRight(value, " \t")
	// 2. Заменяем терминальные нули на перевод строки
	value = processWithLogging(value)
	return value, nil
}
func processWithLogging(s string) string {
	result := make([]rune, 0, len(s))
	for _, char := range s {
		if char == '\x00' {
			char = '\n'
		}
		result = append(result, char)
	}
	return string(result)
}
