package daemon

import (
	"context"
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

	// Создаем контекст с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())

	fmt.Println("Daemon successfully started")

	// Запускаем сбор статистики
	collector := app.NewCollector(cfg, app.Options{
		DataInterval: dataInterval,
		LogInterval:  logInterval,
	})
	collector.Start(ctx)

	// Ожидаем сигнал остановки
	<-stop

	// Отменяем контекст, чтобы остановить горутины сбора статистики
	cancel()

	fmt.Println("Daemon successfully stopped")

	// Очищаем файлы перед выходом
	_ = os.Remove(PidFile)
	_ = os.Remove(LogFile)
}
