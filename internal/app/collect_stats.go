package app

import (
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
			log.Println("Collect statistics ...")

		case <-done:
			return
		}
	}
}
