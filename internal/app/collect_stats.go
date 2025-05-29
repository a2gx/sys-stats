package app

import (
	"fmt"
	"github.com/a2gx/sys-stats/internal/config"
	"log"
	"time"

	"github.com/a2gx/sys-stats/internal/stats"
)

func CollectStats(cfg *config.Config, done <-chan bool) {
	fmt.Printf("Collecting statistics with configuration: %+v\n", cfg)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Собираем статистику раз в секунду
			res, err := stats.GetLoadAverage()
			if err != nil {
				log.Printf("error: %v", err)
			} else {
				log.Printf("result: %f", res)
			}

		case <-done:
			return
		}
	}
}
