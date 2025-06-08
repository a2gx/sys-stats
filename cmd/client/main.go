package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	daemon "github.com/a2gx/sys-stats/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var (
		host         string
		port         int
		logInterval  int32
		dataInterval int32
	)

	rootCmd := &cobra.Command{
		Use:   "daemon-client",
		Short: "GRPC client for daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			addrGRPC := fmt.Sprintf("%s:%d", host, port)
			return runClient(addrGRPC, logInterval, dataInterval)
		},
	}

	rootCmd.Flags().StringVar(&host, "host", "localhost", "Host fot gRPC server")
	rootCmd.Flags().IntVar(&port, "port", 50055, "Port fot gRPC server")
	rootCmd.Flags().Int32Var(&logInterval, "log-interval", 5, "Log interval in seconds")
	rootCmd.Flags().Int32Var(&dataInterval, "data-interval", 10, "Data interval in seconds")

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error executing command: %v", err)
		os.Exit(1)
	}
}

func runClient(addr string, logInterval, dataInterval int32) error {
	fmt.Printf("Starting daemon client (%s)...\n", addr)
	fmt.Printf("Log interval: %d seconds, Data interval: %d seconds\n", logInterval, dataInterval)

	// Создаем gRPC соединение с сервером
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server at %s: %w", addr, err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Printf("failed to close gRPC connection: %v\n", err)
		}
	}()

	// Настройки для корректного завершения работы клиента
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Канал для сигналов завершения
	cnStop := make(chan os.Signal, 1)
	signal.Notify(cnStop, syscall.SIGINT, syscall.SIGTERM)

	// Канал для ошибок потока
	cnErr := make(chan error, 1)

	// Создаем gRPC клиент
	client := daemon.NewDaemonStreamClient(conn)

	go func() {
		// Создаем поток для получения данных
		stream, err := client.SysStatsStream(ctx, &daemon.SysStatsStreamRequest{
			LogInterval:  logInterval,
			DataInterval: dataInterval,
		})
		if err != nil {
			cnErr <- fmt.Errorf("failed create stream: %w", err)
			return
		}

		fmt.Printf("Connected to gRPC server at %s\n", addr)
		fmt.Printf("Press Ctrl+C to stop the client.\n")
		fmt.Printf("----------------------------------------\n")

		var counter int

		for {
			resp, err := stream.Recv()

			if err == io.EOF {
				cnErr <- nil // спокойно завершаем поток
				return
			}
			if err != nil {
				cnErr <- fmt.Errorf("failed to receive data from stream: %w", err)
				return
			}

			counter++
			now := time.Now().Format("15:04:05")

			fmt.Printf("Received data at %s (count: %d):\n", now, counter)
			fmt.Printf("\tLoad Average: %.2f\n", resp.LoadAverage)
			fmt.Printf("\tCPU Usage:\n\t\tUser: %d%%,\n\t\tSystem: %d%%,\n\t\tIdle: %d%%\n",
				resp.CpuUsage.User, resp.CpuUsage.System, resp.CpuUsage.Idle)
			fmt.Printf("\tDisk Usage:\n\t\tTPS: %d,\n\t\tKBps: %d\n",
				resp.DiskUsage.Tps, resp.DiskUsage.KbPs)
			fmt.Printf("----------------------------------------\n")
		}
	}()

	// Ожидаем сигнала завершения или ошибки потока
	select {
	case <-cnStop:
		fmt.Println("Received termination signal, stopping client...")
		cancel()
		return nil
	case err := <-cnErr:
		if err != nil {
			return err
		}

		fmt.Println("Stream ended")
		return nil
	}
}
