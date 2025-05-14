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
	go readLogsFile(logf, currentPos, stop, cnDone)

	// Ожидаем сигнал завершения
	<-stop
	close(cnDone)

	fmt.Println("\nDaemon successfully stopped")
	return nil
}

func readLogsFile(logf *os.File, currentPos int64, stop chan os.Signal, done <-chan bool) {
	for {
		select {
		case <-time.After(100 * time.Millisecond):
			// Проверяем, существует ли файл
			if _, err := os.Stat(LogFile); os.IsNotExist(err) {
				fmt.Println("Log file deleted. Stopping daemon...")
				stop <- syscall.SIGTERM // Отправляем сигнал в канал stop
				return
			}

			// Проверяем, есть ли новые данные
			fi, err := logf.Stat()
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "error getting file info: %v\n", err)
				continue
			}

			if size := fi.Size(); size > currentPos {
				// Если размер файла увеличился, читаем новые данные
				logf.Seek(currentPos, io.SeekStart)

				newScanner := bufio.NewScanner(logf)
				for newScanner.Scan() {
					fmt.Println(newScanner.Text())
				}

				currentPos, _ = logf.Seek(0, io.SeekCurrent)
			}

		case <-done:
			return
		}
	}
}
