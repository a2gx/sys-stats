package daemon

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
)

func (dm *DaemonManager) StopDaemon() error {
	cleanup := func() {
		// Очистка ресурсов если была критическая ошибка
		_ = os.Remove(dm.pidFilePath)
		_ = os.Remove(dm.logFilePath)
	}

	// Проверяем существование PID файла
	data, err := os.ReadFile(dm.pidFilePath)
	if err != nil {
		return fmt.Errorf("daemon not running: %w", err)
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		cleanup()
		return fmt.Errorf("invalid PID in file: %v", err)
	}

	// Ищем запущенный процесс по PID через абстракцию
	process, err := dm.processManager.FindProcess(pid)
	if err != nil {
		cleanup()
		return fmt.Errorf("cannot find process %d: %w", pid, err)
	}

	// Проверяем, существует ли процесс
	if err := process.Signal(syscall.Signal(0)); err != nil {
		cleanup()
		return fmt.Errorf("process with PID %d already finished: %w", pid, err)
	}

	// Отправляем сигнал для корректного завершения
	if err := process.Signal(syscall.SIGTERM); err != nil {
		cleanup()
		return fmt.Errorf("failed to send SIGTERM to process: %w", err)
	}

	fmt.Printf("Daemon with PID %d successfully stopped\n", pid)
	return nil
}
