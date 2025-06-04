package main

import (
	"fmt"
	"log"

	"github.com/a2gx/sys-stats/internal/config"
	"github.com/spf13/pflag"
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
}
