package daemon

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
)

func StopDaemon() error {
	defer func() {
		// Если возникла ошибка, удаляем PID и лог файлы
		_ = os.Remove(PidFile)
		_ = os.Remove(LogFile)
	}()

	// Проверяем существование PID файла
	data, err := os.ReadFile(PidFile)
	if err != nil {
		return fmt.Errorf("daemon not running or error reading PID file: %w", err)
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return fmt.Errorf("invalid PID in file: %v", err)
	}

	// Ищем запущенный процесс по PID
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("cannot find process %d: %w", pid, err)
	}

	// Проверяем, существует ли процесс
	// TODO Unix-специфично
	if err := process.Signal(syscall.Signal(0)); err != nil {
		return fmt.Errorf("process with PID %d already finished: %w", pid, err)
	}

	// Отправляем сигнал для корректного завершения
	if err := process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to send SIGTERM to process: %w", err)
	}

	fmt.Printf("Daemon with PID %d successfully stopped\n", pid)
	return nil
}
