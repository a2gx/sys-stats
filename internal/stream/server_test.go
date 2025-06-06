package stream

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/a2gx/sys-stats/internal/app"
	"github.com/a2gx/sys-stats/internal/stats"
	daemon "github.com/a2gx/sys-stats/proto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func startTestServer(t *testing.T) (*Server, string, func()) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	srv := NewServer(lis.Addr().String())
	grpcServer := grpc.NewServer()
	daemon.RegisterDaemonStreamServer(grpcServer, srv)

	go grpcServer.Serve(lis)

	return srv, lis.Addr().String(), func() {
		grpcServer.Stop()
		lis.Close()
	}
}

func TestSysStatsStream_Broadcast(t *testing.T) {
	srv, addr, cleanup := startTestServer(t)
	defer cleanup()

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	client := daemon.NewDaemonStreamClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	stream, err := client.SysStatsStream(ctx, &daemon.SysStatsStreamRequest{})
	require.NoError(t, err)

	// Отправляем данные через Broadcast
	go func() {
		time.Sleep(100 * time.Millisecond)
		srv.Broadcast(&app.History{
			LoadAverage: 1.23,
			CPUUsage:    stats.CPUStat{User: 10, System: 20, Idle: 70},
			DiskUsage:   stats.DiskUsage{TPS: 100, KBps: 2048},
		})
	}()

	resp, err := stream.Recv()
	require.NoError(t, err)
	assert.InDelta(t, 1.23, resp.LoadAverage, 0.01)
	assert.InDelta(t, 10, resp.CpuUsage.User, 0.01)
	assert.InDelta(t, 20, resp.CpuUsage.System, 0.01)
	assert.InDelta(t, 70, resp.CpuUsage.Idle, 0.01)
	assert.InDelta(t, 100, resp.DiskUsage.Tps, 0.01)
	assert.InDelta(t, 2048, resp.DiskUsage.KbPs, 0.01)
}
