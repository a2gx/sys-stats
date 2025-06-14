package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/a2gx/sys-stats/internal/server"
)

type StartDaemonFlags struct {
	Background bool
	ConfigFile string
	Host       string
	Port       int
}

// StartDaemon запускает демон в фоновом режиме или в текущем контексте.
func (dm *DaemonManager) StartDaemon(flags StartDaemonFlags) error {
	if !flags.Background {
		// Запускаем процесс в текущем контексте
		return dm.runProcess()
	}

	// Что бы запустить процесс в фоне, нужно перенаправить
	// вывод логов в файл и отвязать процесс от терминала
	// PID (process ID) - идентификатор процесса сохраним в файл

	// Проверяем, не запущен ли уже демон
	if _, err := os.Stat(dm.pidFilePath); err == nil {
		return fmt.Errorf("daemon already running")
	}

	// Создаем лог-файл
	logf, err := os.OpenFile(dm.logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("ошибка открытия лог-файла: %w", err)
	}
	defer func() {
		if err := logf.Close(); err != nil {
			fmt.Printf("failed to close log file: %v\n", err)
		}
	}()

	// Создаем команду для запуска
	cmd := exec.Command(os.Args[0], "run",
		"--host", flags.Host,
		"--port", fmt.Sprintf("%d", flags.Port),
		"--config", flags.ConfigFile,
	)

	// Настраиваем перенаправление вывода
	cmd.Stdout = logf
	cmd.Stderr = logf

	// Настройка для отсоединения от родителя
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}

	// Запускаем процесс
	if err := dm.processManager.StartProcess(cmd); err != nil {
		return fmt.Errorf("error starting daemon: %w", err)
	}

	// Сохраняем PID
	pid := strconv.Itoa(cmd.Process.Pid)
	if err := os.WriteFile(dm.pidFilePath, []byte(pid), 0644); err != nil {
		return fmt.Errorf("error writing daemon PprocessID: %w", err)
	}

	fmt.Printf("Daemon started in background mode. PID: %s\n", pid)
	return nil
}

// runProcess запускает основной процесс демона
func (dm *DaemonManager) runProcess() error {
	// Настраиваем обработку сигналов
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// Запуск gRPC сервера
	srv := server.NewServer(dm.cfg)
	go func() {
		if err := srv.Start(dm.cfg.AddrGRPC); err != nil {
			fmt.Printf("gRPC server error: %v\n", err)
			stop <- syscall.SIGTERM
		}
	}()

	fmt.Println("Daemon successfully started")

	// Ожидаем сигнал остановки
	<-stop

	fmt.Println("Daemon successfully stopped")

	// Очистка ресурсов при завершении
	_ = os.Remove(dm.pidFilePath)
	_ = os.Remove(dm.logFilePath)

	return nil
}
