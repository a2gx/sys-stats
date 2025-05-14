package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

const (
	PidFile = "/tmp/daemon.pid"
	LogFile = "/tmp/daemon.log"
)

func StartDaemon(detect bool, logInterval, dataInterval int) error {
	if !detect {
		RunDaemon(logInterval, dataInterval)
		return nil
	}

	// Создаем и запускаем фоновый процесс
	// Нужно запустить эту же команду отдельным процессом и освободить терминал
	// PID (process ID) - идентификатор процесса сохраним в файл

	// Проверка: если PID-файл существует, значит демон уже запущен
	if _, err := os.Stat(PidFile); err == nil {
		return fmt.Errorf("daemon already running")
	}

	// Подготавливаем вывод в лог-файл
	logf, err := os.OpenFile(LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("error opening log file: %w", err)
	}

	// Создаем новый процесс
	cmd := exec.Command(os.Args[0], "run",
		"--log-interval", fmt.Sprintf("%d", logInterval),
		"--data-interval", fmt.Sprintf("%d", dataInterval))

	// Перенаправляем stdout и stderr в лог-файл
	// Теперь все что выводится в stdout и stderr будет записываться в лог-файл
	// Это касается только вывода через пакет log
	cmd.Stdout = logf
	cmd.Stderr = logf

	// Отвязываем новый процесс от родителя
	// Позволяет демону работать в фоновом режиме, стандартная практика для Unix демонов.
	// TODO на Windows это может вызвать ошибку. Условная компиляция или отключить в рантайме ??
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}

	// Запускаем процесс
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting daemon: %w", err)
	}

	// Сохраняем PID в файл
	pid := strconv.Itoa(cmd.Process.Pid)
	if err := os.WriteFile(PidFile, []byte(pid), 0644); err != nil {
		return fmt.Errorf("error writing daemon pid: %w", err)
	}

	fmt.Printf("Daemon started in background mode. PID: %s\n", pid)
	return nil
}
