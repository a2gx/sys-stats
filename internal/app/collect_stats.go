package app

import (
	"log"
	"time"
)

func CollectStats(done <-chan bool) {
	for {
		select {
		case <-time.After(1 * time.Second):
			// Собираем статистику раз в секунду
			log.Println("Collect statistics ...")

		case <-done:
			return
		}
	}
}
