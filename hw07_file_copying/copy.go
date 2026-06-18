package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync/atomic"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

const (
	bufferSize = 64 * 1024 // 64 KB
)

// copyData отвечает только за цикл чтения/записи.
func copyData(input *os.File, output *os.File, bytesToCopy int64, progressChan chan<- int64) (int64, error) {
	var copied int64
	buffer := make([]byte, bufferSize)

	for atomic.LoadInt64(&copied) < bytesToCopy {
		remaining := bytesToCopy - atomic.LoadInt64(&copied)
		chunkSize := bufferSize
		if remaining < int64(bufferSize) {
			chunkSize = int(remaining)
		}

		n, err := input.Read(buffer[:chunkSize])
		if n > 0 {
			if _, werr := output.Write(buffer[:n]); werr != nil {
				return copied, fmt.Errorf("ошибка записи: %w", werr)
			}
			// Атомарно увеличиваем счетчик
			newCopied := atomic.AddInt64(&copied, int64(n))
			// Отправляем новое значение в канал для обновления прогресс-бара
			progressChan <- newCopied
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			return copied, fmt.Errorf("ошибка чтения: %w", err)
		}
	}
	return atomic.LoadInt64(&copied), nil
}

// validateAndPrepare проверяет параметры и вычисляет bytesToCopy.
func validateAndPrepare(input *os.File, offset, limit int64) (int64, error) {
	info, err := input.Stat()
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrUnsupportedFile, err.Error())
	}
	if info.Size() == 0 {
		return 0, fmt.Errorf("файл имеет нулевой размер или не поддерживается (например, /dev/urandom)")
	}

	fileSize := info.Size()
	if offset >= fileSize {
		return 0, fmt.Errorf("%w: offset %d превышает размер файла %d", ErrOffsetExceedsFileSize, offset, fileSize)
	}

	bytesToCopy := limit
	if limit == 0 || limit > fileSize-offset {
		bytesToCopy = fileSize - offset
	}

	if _, err := input.Seek(offset, io.SeekStart); err != nil {
		return 0, fmt.Errorf("не удалось установить позицию в источнике: %w", err)
	}
	return bytesToCopy, nil
}

// openFiles инкапсулирует логику открытия файлов.
func openFiles(fromPath, toPath string) (*os.File, *os.File, error) {
	input, err := os.Open(fromPath)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", ErrUnsupportedFile, err.Error())
	}

	output, err := os.Create(toPath)
	if err != nil {
		input.Close() // Не забываем закрыть input в случае ошибки создания output
		return nil, nil, fmt.Errorf("%w: %s", ErrUnsupportedFile, err.Error())
	}
	return input, output, nil
}

// startProgressBar запускает горутину и возвращает канал done и указатель на счетчик.
func startProgressBar(bytesToCopy int64, progressChan <-chan int64) {
	go func() {
		for current := range progressChan {
			if bytesToCopy > 0 {
				percent := float64(current) / float64(bytesToCopy) * 100
				// \r возвращает курсор в начало строки для перезаписи
				fmt.Printf("\rПрогресс: %.1f%% (%d/%d байт)", percent, current, bytesToCopy)
			}
		}
		// Этот код выполнится после закрытия канала progressChan
		fmt.Println() // Переводим курсор на новую строку после завершения
	}()
}

// --- Главная функция: Оркестратор ---

func Copy(fromPath, toPath string, offset, limit int64) error {
	input, output, err := openFiles(fromPath, toPath)
	if err != nil {
		if input != nil {
			input.Close()
		}
		if output != nil {
			output.Close()
		}
		return err
	}
	defer input.Close()
	defer output.Close()

	bytesToCopy, err := validateAndPrepare(input, offset, limit)
	if err != nil {
		return err
	}

	// 1. Создаем канал для передачи прогресса
	progressChan := make(chan int64)

	// 2. Запускаем прогресс-бар, передав ему канал для чтения
	startProgressBar(bytesToCopy, progressChan)

	// 3. Запускаем копирование, передав ему канал для записи
	copiedCount, err := copyData(input, output, bytesToCopy, progressChan)

	// 4. ОБЯЗАТЕЛЬНО закрываем канал, чтобы горутина прогресс-бара завершила цикл for-range
	close(progressChan)

	if err != nil {
		return err
	}

	if copiedCount != bytesToCopy {
		return fmt.Errorf("скопировано %d байт из ожидаемых %d", copiedCount, bytesToCopy)
	}

	if err := output.Sync(); err != nil {
		return fmt.Errorf("ошибка синхронизации файла: %w", err)
	}
	return nil
}
