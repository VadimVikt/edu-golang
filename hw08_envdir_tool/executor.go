package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 127
	}
	// 1. Копируем текущее окружение
	newEnv := os.Environ()
	// 2. Удаляем существующие переменные, которые мы хотим переопределить
	keysToRemove := make(map[string]bool)
	for k, v := range env {
		if v.NeedRemove {
			keysToRemove[k] = true
		}
	}
	filteredEnv := make([]string, 0, len(newEnv))
	for _, kv := range newEnv {
		parts := strings.SplitN(kv, "=", 2)
		key := parts[0]
		if !keysToRemove[key] {
			filteredEnv = append(filteredEnv, kv)
		}
	}
	// 3. Добавляем новые переменные из env
	for k, v := range env {
		filteredEnv = append(filteredEnv, k+"="+v.Value)
	}
	// Подготовка команды
	execCmd := exec.Command(cmd[0], cmd[1:]...) //nolint
	execCmd.Env = filteredEnv
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	// Запуск команды
	err := execCmd.Run()
	// Обработка кода выхода
	var exitCode int
	//nolint
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			// Программа завершилась с ненулевым кодом
			if status, ok := exitError.Sys().(interface{ ExitStatus() int }); ok {
				exitCode = status.ExitStatus()
			} else {
				exitCode = 1 // Неизвестная ошибка
			}
		}
	} else {
		// Успешное завершение
		exitCode = 0
	}
	return exitCode
}
