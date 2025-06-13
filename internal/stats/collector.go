package stats

import (
	"log"
	"sync"

	"github.com/a2gx/sys-stats/internal/config"
)

type Collector struct {
	history []*History
	cfg     *config.Config
	opt     CollectorOptions
}

type CollectorOptions struct {
	LogInterval  int
	DataInterval int
}

type History struct {
	LoadAverage float64
	CPUUsage    CPUStat
	DiskUsage   DiskUsage
}

func NewCollector(cfg *config.Config, opts CollectorOptions) *Collector {
	return &Collector{
		cfg: cfg,
		opt: opts,
	}
}

func (c *Collector) HistoryCollect() *History {
	entry := History{}

	var wg sync.WaitGroup
	var mu sync.Mutex

	collects := []func(){
		func() {
			if !c.cfg.LoadAverage {
				return
			}

			stat, err := GetLoadAverage()
			if err != nil {
				log.Printf("failed to collect `LoadAverage`: %v", err)
				return
			}

			mu.Lock()
			entry.LoadAverage = stat
			mu.Unlock()
		},
		func() {
			if !c.cfg.CPUUsage {
				return
			}

			stat, err := GetCpuUsage()
			if err != nil {
				log.Printf("failed to collect `CPUUsage`: %v", err)
				return
			}

			mu.Lock()
			entry.CPUUsage = stat
			mu.Unlock()
		},
		func() {
			if !c.cfg.DiskUsage {
				return
			}

			stat, err := GetDiskUsage()
			if err != nil {
				log.Printf("failed to collect `DiskUsage`: %v", err)
				return
			}

			mu.Lock()
			entry.DiskUsage = stat
			mu.Unlock()
		},
	}

	for _, collect := range collects {
		wg.Add(1)

		go func(collectFunc func()) {
			defer wg.Done()
			collectFunc()
		}(collect)
	}
	wg.Wait()

	return &entry
}

func (c *Collector) HistoryCalculate(history []*History) *History {
	if len(history) == 0 {
		return nil
	}

	result := &History{
		LoadAverage: 0,
		CPUUsage:    CPUStat{},
		DiskUsage:   DiskUsage{},
	}

	// Суммируем все метрики
	for _, entry := range history {
		result.LoadAverage += entry.LoadAverage

		result.CPUUsage.User += entry.CPUUsage.User
		result.CPUUsage.System += entry.CPUUsage.System
		result.CPUUsage.Idle += entry.CPUUsage.Idle

		result.DiskUsage.TPS += entry.DiskUsage.TPS
		result.DiskUsage.KBps += entry.DiskUsage.KBps
	}

	// Вычисляем средние значения
	count := float64(len(history))

	result.LoadAverage /= count

	result.CPUUsage.User /= count
	result.CPUUsage.System /= count
	result.CPUUsage.Idle /= count

	result.DiskUsage.TPS /= count
	result.DiskUsage.KBps /= count

	// Возвращаем результат
	return result
}
