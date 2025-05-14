package daemon

import (
	"fmt"
	"github.com/a2gx/sys-stats/internal/app"
	"os"
	"os/signal"
	"syscall"
)

func RunDaemon(logInterval, dataInterval int) {
	// Канал для перехвата сигналов для корректной остановки
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// Каналы для синхронизации
	cnDone := make(chan bool)

	fmt.Println("Daemon successfully started")

	// Запускаем горутины сбора и отправки статистики
	go app.CollectStats(cnDone)
	go app.OutputStats(logInterval, dataInterval, cnDone)

	// Ожидаем сигнал завершения
	<-stop
	fmt.Println("Daemon successfully stopped")

	// Останавливаем горутины
	close(cnDone)
}
