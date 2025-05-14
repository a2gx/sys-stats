package daemon

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/a2gx/sys-stats/internal/app"
)

func RunProcess(logInterval, dataInterval int) {
	// Устанавливаем формат логирования
	log.SetFlags(log.LstdFlags)
	log.SetPrefix("daemon: ")

	// Канал для безопасной остановки процесса
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	cnDone := make(chan bool)

	fmt.Println("Daemon successfully started")

	// Запускаем горутины сбора и отправки статистики
	go app.CollectStats(cnDone)
	go app.OutputStats(logInterval, dataInterval, cnDone)

	// Останавливаем горутины
	<-stop
	close(cnDone)

	fmt.Println("Daemon successfully stopped")

	// Удаляем PID файл и лог файл, если они существуют
	_ = os.Remove(PidFile)
	_ = os.Remove(LogFile)
}
