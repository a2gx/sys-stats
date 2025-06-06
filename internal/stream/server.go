package stream

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"

	"github.com/a2gx/sys-stats/internal/app"
	daemon "github.com/a2gx/sys-stats/proto"
)

type Server struct {
	daemon.UnimplementedDaemonStreamServer

	addr        string
	mu          sync.RWMutex
	subscribers map[subscriber]struct{}
}

type subscriber chan *daemon.SysStatsStreamResponse

func NewServer(addr string) *Server {
	return &Server{
		addr:        addr,
		subscribers: make(map[subscriber]struct{}),
	}
}

func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	daemon.RegisterDaemonStreamServer(grpcServer, s)

	go func() {
		<-ctx.Done() // Ждем завершения контекста
		log.Println("grpc server is stopping...")
		grpcServer.GracefulStop() // Останавливаем сервер
	}()

	log.Println("grpc server started...")
	return grpcServer.Serve(lis)
}

// SysStatsStream реализует потоковую передачу данных
// Пока соединение открыто, клиент будет получать обновления
func (s *Server) SysStatsStream(_ *daemon.SysStatsStreamRequest, stream daemon.DaemonStream_SysStatsStreamServer) error {
	ch := make(subscriber, 100) // буферизированный канал для подписчика

	// Добавляем нового подписчика при подключении
	s.mu.Lock()
	s.subscribers[ch] = struct{}{}
	s.mu.Unlock()

	// Удаляем подписчика при отключении
	defer func() {
		s.mu.Lock()
		delete(s.subscribers, ch)
		close(ch) // Закрываем канал
		s.mu.Unlock()
	}()

	for {
		select {
		case <-stream.Context().Done():
			return nil // Завершаем поток при закрытии контекста
		case msg := <-ch:
			if err := stream.Send(msg); err != nil {
				return err // Ошибка при отправке сообщения
			}
		}
	}
}

// Broadcast рассылает данные всем подписчикам
func (s *Server) Broadcast(d *app.History) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	msg := &daemon.SysStatsStreamResponse{
		LoadAverage: float32(d.LoadAverage),
		CpuUsage: &daemon.CpuUsage{
			User:   float32(d.CPUUsage.User),
			System: float32(d.CPUUsage.System),
			Idle:   float32(d.CPUUsage.Idle),
		},
		DiskUsage: &daemon.DiskUsage{
			Tps:  float32(d.DiskUsage.TPS),
			KbPs: float32(d.DiskUsage.KBps),
		},
	}

	for ch := range s.subscribers {
		select {
		case ch <- msg:
		default: // non-blocking send
		}
	}
}
