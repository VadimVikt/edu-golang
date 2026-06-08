package main

import (
	"bufio"
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
	// Place your code here
	result := make(Environment)
	listFile, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range listFile {
		fmt.Printf("List of file - %s\n", file.Name())
		if strings.Contains(file.Name(), "=") {
			fmt.Printf("Имя файла %s содержит заперщенный символ = \n", file.Name())
			result[file.Name()] = EnvValue{"", false}
			continue
		}
		envName, err := prepareEnv(dir + "/" + file.Name())
		if err != nil || envName == "" {
			fmt.Printf("Файл %s пуст или произошла ошибка чтения переменную установить невозможно \n", file.Name())
			result[file.Name()] = EnvValue{"", false}
			continue
		}
		if value, exist := os.LookupEnv(envName); exist == true {
			os.Unsetenv(value)
		}
		//fmt.Println(envName)
		result[file.Name()] = EnvValue{envName, true}
	}
	fmt.Println("Результат работы - ", result)
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
		if err := scanner.Err(); err != nil && err != io.EOF {
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
	var result []rune
	for i, char := range s {
		if char == '\x00' {
			fmt.Printf("Найден нулевой символ на позиции %d\n", i)
			// Замена на \n
			char = '\n'
		}
		result = append(result, char)
	}
	return string(result)
}
