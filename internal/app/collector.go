package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/a2gx/sys-stats/internal/config"
	"github.com/a2gx/sys-stats/internal/stats"
)

type StatsCollector struct {
	history []History
	mu      sync.Mutex
	cfg     *config.Config
	opt     Options
}

type Options struct {
	DataInterval int
	LogInterval  int
}

type History struct {
	CPUUsage    stats.CPUStat
	LoadAverage float64
	DiskUsage   stats.DiskUsage
}

func NewCollector(cfg *config.Config, opts Options) *StatsCollector {
	return &StatsCollector{
		history: make([]History, 0, opts.LogInterval),
		cfg:     cfg,
		opt:     opts,
	}
}

func (sc *StatsCollector) Start(ctx context.Context) {
	go sc.collectStats(ctx)
	go sc.outputStats(ctx)
}

func (sc *StatsCollector) collectStats(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sc.saveHistory()
		case <-ctx.Done():
			return
		}
	}
}

func (sc *StatsCollector) outputStats(ctx context.Context) {
	var dataInterval, logInterval = sc.opt.DataInterval, sc.opt.LogInterval
	var delay = dataInterval - logInterval

	// Сначала ждем dataInterval секунд, чтобы собрать статистику
	select {
	case <-time.After(time.Duration(delay) * time.Second):
	case <-ctx.Done():
		return
	}

	ticker := time.NewTicker(time.Duration(logInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			d := sc.readHistory()
			fmt.Printf("Output: %+v\n", d)
		case <-ctx.Done():
			return
		}
	}
}

func (sc *StatsCollector) saveHistory() {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Сбор статистики и добавление в историю
	sc.history = append(sc.history, historyCollect(sc.cfg))
}

func (sc *StatsCollector) readHistory() History {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	count := len(sc.history)
	if count == 0 {
		return History{} // возврат пустой статистики, если нет данных
	}

	result := History{
		CPUUsage:    calculateCPUUsage(sc.history, count),
		LoadAverage: calculateLoadAverage(sc.history, count),
		DiskUsage:   calculateDiskUsage(sc.history, count),
		// Можно добавить другие метрики, если нужно
	}

	// Очищаем историю
	sc.history = sc.history[:0]

	return result
}
