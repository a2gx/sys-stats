package app

import (
	"log"
	"time"
)

func CollectStats(done <-chan bool) {
	// Собираем статистику раз в секунду
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Collect statistics ...")

		case <-done:
			return
		}
	}
}
