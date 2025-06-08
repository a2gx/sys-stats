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

// LogsDaemon отображает логи демона sys-stats.
func (dm *DaemonManager) LogsDaemon() error {
	// Проверяем, запущен ли демон
	if _, err := os.Stat(dm.pidFilePath); os.IsNotExist(err) {
		return fmt.Errorf("daemon not running: %w", err)
	}

	// Проверяем наличие файла логов
	if _, err := os.Stat(dm.logFilePath); os.IsNotExist(err) {
		return fmt.Errorf("loga file not found: %w", err)
	}

	// Открываем файл логов
	logFile, err := os.Open(dm.logFilePath)
	if err != nil {
		return fmt.Errorf("could not open loga file: %w", err)
	}
	defer func() {
		if err := logFile.Close(); err != nil {
			fmt.Printf("failed to close log file: %v\n", err)
		}
	}()

	fmt.Printf("Displaying logs from %s:\n", dm.logFilePath)
	fmt.Printf("Press Ctrl+C to stop displaying logs.\n")

	// Настраиваем канал для безопасной остановки процесса
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// Канал для завершения горутины мониторинга
	done := make(chan bool)

	currentPos, err := dm.displayCurrentLogs(logFile)
	if err != nil {
		return fmt.Errorf("error reading current logs: %w", err)
	}

	// Запускаем горутину для чтения и отслеживания новых логов
	go dm.monitorNewLogs(logFile, currentPos, stop, done)

	// Ожидаем сигнала остановки
	<-stop
	close(done)

	fmt.Println("\nStopping log display...")
	return nil
}

// displayCurrentLogs отображает текущее содержимое файла логов
func (dm *DaemonManager) displayCurrentLogs(file *os.File) (int64, error) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error for scanning log file: %w", err)
	}

	return file.Seek(0, io.SeekCurrent)
}

// monitorNewLogs отслеживает новые записи в файле логов и выводит их
func (dm *DaemonManager) monitorNewLogs(file *os.File, startPos int64, stop chan os.Signal, done chan bool) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	currentPos := startPos

	for {
		select {
		case <-ticker.C:
			if _, err := os.Stat(dm.logFilePath); os.IsNotExist(err) {
				fmt.Println("log file has been deleted, stopping monitoring.")
				stop <- syscall.SIGTERM
				return
			}

			f, err := file.Stat()
			if err != nil {
				fmt.Printf("failed reading log file stats: %v\n", err)
				continue
			}

			// Если размер файла увеличился, читаем новые строки
			if size := f.Size(); size > currentPos {
				// Перемещаем указатель на последнюю прочитанную позицию
				if _, err := file.Seek(currentPos, io.SeekStart); err != nil {
					fmt.Printf("failed to seek in log file: %v\n", err)
					continue
				}

				// Читаем новые строки
				newScanner := bufio.NewScanner(file)
				for newScanner.Scan() {
					fmt.Println(newScanner.Text())
				}

				currentPos, _ = file.Seek(0, io.SeekCurrent)
			}

		case <-done:
			return
		}
	}
}
