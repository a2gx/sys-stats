package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/a2gx/sys-stats/internal/config"
	"github.com/a2gx/sys-stats/proto"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
)

func main() {
	configFile := pflag.StringP("config", "c", "/configs/config.yaml", "Path to configuration file")
	pflag.Parse()

	cfg, err := config.NewConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	addr := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)

	fmt.Printf("GRPC client starting: %s\n", addr)

	var conn *grpc.ClientConn
	var attempts int

	for {
		conn, err = grpc.NewClient(addr)
		if err == nil {
			break // Connection successful
		}

		attempts++

		fmt.Printf("Failed to connection, attempt: %d\n", attempts)
		time.Sleep(3 * time.Second)

		if attempts >= 10 {
			log.Fatalf("Failed to connection, error: %v\n", err)
		}
	}
	defer func() {
		if conn == nil {
			return
		}

		if err := conn.Close(); err != nil {
			fmt.Printf("Failed to close gRPC connection: %v", err)
		}
	}()

	client := daemon.NewDaemonStreamClient(conn)
	stream, err := client.SysStatsStream(context.Background(), &daemon.SysStatsStreamRequest{})
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("Stream ended")
			break
		}

		fmt.Printf("Load Average: \t%.2f\n", resp.LoadAverage)
		fmt.Printf("Cpu Usage: \t%+v\n", resp.CpuUsage)
		fmt.Printf("Disk Usage: \t%+v\n", resp.DiskUsage)
	}
}
