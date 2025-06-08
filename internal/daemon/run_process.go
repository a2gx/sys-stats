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
	"github.com/a2gx/sys-stats/internal/stream"
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

	addr := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)
	streamServer := stream.NewServer(addr)

	// Запускаем сбор статистики
	collector := app.NewCollector(cfg, streamServer.Broadcast, app.Options{
		DataInterval: dataInterval,
		LogInterval:  logInterval,
	})
	collector.Start(ctx)

	if err := streamServer.Start(ctx); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}

	// Ожидаем сигнал остановки
	<-stop

	// Отменяем контекст, чтобы остановить горутины сбора статистики
	cancel()

	fmt.Println("Daemon successfully stopped")

	// Очищаем файлы перед выходом
	_ = os.Remove(PidFile)
	_ = os.Remove(LogFile)
}
