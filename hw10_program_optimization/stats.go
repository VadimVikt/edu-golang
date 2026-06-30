package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	jsoniter "github.com/bytedance/sonic" //nolint
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	domainLower := strings.ToLower(domain)
	patternStr := "\\." + regexp.QuoteMeta(domainLower)
	domainRegexp, err := regexp.Compile(patternStr)
	if err != nil {
		return nil, fmt.Errorf("invalid domain pattern: %w", err)
	}
	// --- НАЧАЛО ОПТИМИЗАЦИИ ---
	const chunkSize = 5 * 1024 * 1024 // Читаем блоками по 5 МБ

	scanner := bufio.NewScanner(r)

	// Настраиваем большой начальный буфер и его максимальный рост.
	buf := make([]byte, 0, chunkSize/2)
	scanner.Buffer(buf, chunkSize)

	var user User // Структура для переиспользования

	for scanner.Scan() {
		line := scanner.Bytes()

		// Быстрая проверка на валидность JSON строки
		if len(line) == 0 || line[0] != '{' || line[len(line)-1] != '}' {
			continue
		}
		// 2. Используем sonic.Unmarsha вместо стандартного
		// Sonic выполняет эту операцию значительно быстрее.
		if err := jsoniter.Unmarshal(line, &user); err != nil {
			continue
		}

		email := user.Email
		atIndex := strings.Index(email, "@")
		if atIndex == -1 {
			continue
		}
		domainPart := strings.ToLower(email[atIndex+1:])
		if domainRegexp.MatchString(domainPart) {
			result[domainPart]++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	return result, nil
}
