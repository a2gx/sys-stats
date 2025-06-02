package daemon

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func LogsDaemon() error {
	if _, err := os.Stat(PidFile); err != nil {
		return fmt.Errorf("daemon not running")
	}

	if _, err := os.Stat(LogFile); err != nil {
		return fmt.Errorf("log file not found: %v", err)
	}

	logf, err := os.Open(LogFile)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer func() {
		if err := logf.Close(); err != nil {
			fmt.Printf("failed to close log file: %v\n", err)
		}
	}()

	// Канал для безопасной остановки процесса
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// Канал для сообщения о завершении
	cnDone := make(chan bool)

	// Выводим текущее содержимое файла
	scanner := bufio.NewScanner(logf)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	// Запоминаем текущую позицию в файле
	currentPos, err := logf.Seek(0, io.SeekCurrent)
	if err != nil {
		return fmt.Errorf("failed to seek current file: %v", err)
	}

	// Запускаем горутину для чтения логов
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Если файл удалили, выходим из цикла и останавливаем чтение логов
				if _, err := os.Stat(LogFile); os.IsNotExist(err) {
					stop <- syscall.SIGTERM
					return
				}

				// Проверяем, есть ли новые данные
				fi, err := logf.Stat()
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "error getting file info: %v\n", err)
					continue
				}

				// Если размер файла увеличился, читаем новые данные
				if size := fi.Size(); size > currentPos {
					_, _ = logf.Seek(currentPos, io.SeekStart)

					newScanner := bufio.NewScanner(logf)
					for newScanner.Scan() {
						fmt.Println(newScanner.Text())
					}

					currentPos, _ = logf.Seek(0, io.SeekCurrent)
				}

			case <-cnDone:
				return
			}
		}
	}()

	// Ожидаем сигнал завершения
	<-stop
	close(cnDone)

	fmt.Println("Daemon successfully stopped")
	return nil
}
