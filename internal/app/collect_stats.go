package app

import (
	"log"
	"time"
	
	"github.com/a2gx/sys-stats/internal/stats"
)

func CollectStats(done <-chan bool) {
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
