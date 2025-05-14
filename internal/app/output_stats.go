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

	for {
		select {
		case <-time.After(time.Duration(logInterval) * time.Second):
			// Отправляем статистику каждые logInterval секунд
			log.Println("Outputting statistics -->")

		case <-done:
			return
		}
	}
}
