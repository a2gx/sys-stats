package daemon

import (
	"fmt"
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

	// Канал для безопасной остановки процесса
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// Канал для сообщения о завершении
	cnDone := make(chan bool)

	// Запускаем горутину для чтения логов
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// TODO как то надо выводить логи

			case <-cnDone:
				return
			}
		}
	}()

	// Ожидаем сигнал завершения
	<-stop
	close(cnDone)

	fmt.Println("\nDaemon successfully stopped")
	return nil
}
