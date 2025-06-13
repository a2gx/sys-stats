package server

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/a2gx/sys-stats/internal/config"
	"github.com/a2gx/sys-stats/internal/stats"
	"github.com/a2gx/sys-stats/proto/daemon"
	"google.golang.org/grpc"
)

type Server struct {
	daemon.UnimplementedDaemonStreamServer

	mu          sync.Mutex
	cfg         *config.Config
	subscribers map[*subscriber]struct{}
}

type subscriber struct {
	stream       daemon.DaemonStream_SysStatsStreamServer
	cancel       context.CancelFunc
	logInterval  int32
	dataInterval int32
}

// NewServer создает новый экземпляр gRPC сервера
func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg:         cfg,
		subscribers: make(map[*subscriber]struct{}),
	}
}

// Start запускает gRPC сервер
func (s *Server) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	daemon.RegisterDaemonStreamServer(grpcServer, s)

	fmt.Printf("gRpc server started on %s\n", addr)
	return grpcServer.Serve(lis)
}

func (s *Server) SysStatsStream(req *daemon.SysStatsStreamRequest, stream daemon.DaemonStream_SysStatsStreamServer) error {
	logInterval := req.GetLogInterval()
	dataInterval := req.GetDataInterval()

	if logInterval <= 0 {
		logInterval = 5
	}
	if dataInterval <= 0 {
		dataInterval = 15
	}

	ctx, cancel := context.WithCancel(stream.Context())
	sub := &subscriber{
		stream:       stream,
		cancel:       cancel,
		logInterval:  logInterval,
		dataInterval: dataInterval,
	}

	// добавляем подписчика
	s.mu.Lock()
	s.subscribers[sub] = struct{}{}
	s.mu.Unlock()

	fmt.Printf("new subscriber added: logInterval=%d, dataInterval=%d\n", logInterval, dataInterval)

	// запускаем горутину для отправки данных подписчику
	go s.handlerStream(ctx, sub)

	// ожидаем закрытия контекста подписчика
	<-ctx.Done()

	// удаляем подписчика при закрытии контекста
	s.mu.Lock()
	delete(s.subscribers, sub)
	s.mu.Unlock()

	fmt.Println("subscriber removed")

	return nil
}

func (s *Server) handlerStream(ctx context.Context, sub *subscriber) {
	collector := stats.NewCollector(s.cfg, stats.CollectorOptions{
		LogInterval:  int(sub.logInterval),
		DataInterval: int(sub.dataInterval),
	})

	// Буфер для истории
	var history []*stats.History

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	tickerLog := time.NewTicker(time.Duration(sub.logInterval) * time.Second)
	defer tickerLog.Stop()

	for {
		select {
		case <-ctx.Done():
			return // Завершаем горутину при закрытии контекста
		case <-ticker.C:
			entry := collector.HistoryCollect()
			history = append(history, entry)

			// Удаляем старые записи
			if len(history) > int(sub.dataInterval) {
				history = history[1:]
			}
		case <-tickerLog.C:
			// Если история не накопилась, пропускаем отправку клиенту
			if len(history) < int(sub.dataInterval) {
				continue
			}

			avg := collector.HistoryCalculate(history)

			resp := &daemon.SysStatsStreamResponse{
				LoadAverage: float32(avg.LoadAverage),

				CpuUsage: &daemon.CpuUsage{
					User:   float32(avg.CPUUsage.User),
					System: float32(avg.CPUUsage.System),
					Idle:   float32(avg.CPUUsage.Idle),
				},

				DiskUsage: &daemon.DiskUsage{
					Tps:  float32(avg.DiskUsage.TPS),
					KbPs: float32(avg.DiskUsage.KBps),
				},
			}

			if err := sub.stream.Send(resp); err != nil {
				fmt.Printf("failed to send response: %v\n", err)
				sub.cancel()
				return
			}
		}
	}
}
