package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Парсим аргументы командной строки
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go-envdir /path/to/env/dir command [arg...]")
		os.Exit(2)
	}
	dirPath := args[0]
	cmdArgs := args[1:]

	// Считываем переменные окружения из директории
	env, err := ReadDir(dirPath)
	if err != nil {
		log.Fatalf("Failed to read directory %s: %v\n", dirPath, err)
	}
	// Запускаем команду с новым окружением
	exitCode := RunCmd(cmdArgs, env)
	// Завершаем работу утилиты с тем же кодом, что и запущенная программа
	os.Exit(exitCode)
}
