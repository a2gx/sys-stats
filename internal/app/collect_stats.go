package app

import (
	"github.com/a2gx/sys-stats/internal/stats"
	"log"
	"time"
)

func CollectStats(done <-chan bool) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Собираем статистику раз в секунду
			res, err := stats.GetCpuUsage()
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
