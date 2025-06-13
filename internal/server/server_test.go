package server

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/a2gx/sys-stats/internal/config"
	"github.com/a2gx/sys-stats/proto/daemon"
	"github.com/stretchr/testify/require"
)

// Мок для стрима
type mockStream struct {
	daemon.DaemonStream_SysStatsStreamServer
	ctx    context.Context
	sent   []*daemon.SysStatsStreamResponse
	sendCh chan *daemon.SysStatsStreamResponse
}

func (m *mockStream) Context() context.Context {
	return m.ctx
}

func (m *mockStream) Send(resp *daemon.SysStatsStreamResponse) error {
	m.sent = append(m.sent, resp)
	m.sendCh <- resp
	return nil
}

func TestNewServer(t *testing.T) {
	cfg := &config.Config{}
	s := NewServer(cfg)

	require.NotNil(t, s)
	require.Equal(t, cfg, s.cfg)
	require.NotNil(t, s.subscribers)
}

func TestServer_Start(t *testing.T) {
	cfg := &config.Config{}
	s := NewServer(cfg)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := l.Addr().String()
	_ = l.Close()

	done := make(chan struct{})
	go func() {
		err := s.Start(addr)
		require.NoError(t, err)
		close(done)
	}()

	time.Sleep(200 * time.Millisecond)

	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	_ = conn.Close()
}

func TestServer_SysStatsStream(t *testing.T) {
	s := NewServer(&config.Config{
		LoadAverage: true,
		CPUUsage:    true,
		DiskUsage:   true,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream := &mockStream{
		ctx:    ctx,
		sendCh: make(chan *daemon.SysStatsStreamResponse, 1),
	}

	req := &daemon.SysStatsStreamRequest{
		LogInterval:  1,
		DataInterval: 1,
	}

	// Запускаем SysStatsStream в отдельной горутине
	go func() {
		_ = s.SysStatsStream(req, stream)
	}()

	// Ждем первое сообщение
	select {
	case <-stream.sendCh:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for message")
	}

	// Проверяем, что подписчик добавлен
	s.mu.Lock()
	require.Equal(t, 1, len(s.subscribers))
	s.mu.Unlock()

	// Завершаем контекст, чтобы удалить подписчика
	cancel()
	time.Sleep(100 * time.Millisecond)

	// Проверяем, что подписчик удалён
	s.mu.Lock()
	require.Equal(t, 0, len(s.subscribers))
	s.mu.Unlock()
}
