package app

import (
	"log"
	"time"
)

func OutputStats(logInterval, dataInterval int, done <-chan bool) {
	// Сначала ждем dataInterval секунд, чтобы собрать статистику
	select {
	case <-time.After(time.Duration(dataInterval-logInterval) * time.Second):
	case <-done:
		return
	}

	// Отправляем статистику каждые logInterval секунд
	ticker := time.NewTicker(time.Duration(logInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Outputting statistics -->")

		case <-done:
			return
		}
	}
}
