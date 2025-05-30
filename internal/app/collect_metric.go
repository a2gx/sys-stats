package app

import (
	"log"
	"sync"

	"github.com/a2gx/sys-stats/internal/config"
	"github.com/a2gx/sys-stats/internal/stats"
)

func collectMetric(cfg *config.Config) Metrics {
	entry := Metrics{}

	var wg sync.WaitGroup
	var mu sync.Mutex

	collects := []func(){
		func() {
			if !cfg.CPUUsage {
				return
			}

			stat, err := stats.GetCpuUsage()
			if err != nil {
				log.Printf("failed to collect `CPUUsage`: %v", err)
				return
			}

			mu.Lock()
			entry.CPUUsage = stat
			mu.Unlock()
		},
		func() {
			if !cfg.LoadAverage {
				return
			}

			loadAvg, err := stats.GetLoadAverage()
			if err != nil {
				log.Printf("failed to collect `LoadAverage`: %v", err)
				return
			}

			mu.Lock()
			entry.LoadAverage = loadAvg
			mu.Unlock()
		},
	}

	for _, collect := range collects {
		wg.Add(1)

		go func(collectFn func()) {
			defer wg.Done()
			collectFn()
		}(collect)
	}

	wg.Wait()

	return entry
}
