package daemon

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/a2gx/sys-stats/internal/app"
	"github.com/a2gx/sys-stats/internal/config"
)

func RunProcess(cfg *config.Config, logInterval, dataInterval int) {
	// Устанавливаем формат логирования
	log.SetFlags(log.LstdFlags)

	// Канал для безопасной остановки процесса
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	cnDone := make(chan bool)

	fmt.Println("Daemon successfully started")

	// Запускаем горутины сбора и отправки статистики
	go app.CollectStats(cfg, cnDone)
	go app.OutputStats(logInterval, dataInterval, cnDone)

	// Останавливаем горутины
	<-stop
	close(cnDone)

	fmt.Println("Daemon successfully stopped")

	// Очищаем файлы перед выходом
	_ = os.Remove(PidFile)
	_ = os.Remove(LogFile)
}
