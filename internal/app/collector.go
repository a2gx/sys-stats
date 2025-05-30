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
	history []Metrics
	mu      sync.Mutex
	cfg     *config.Config
	opt     Options
}

type Options struct {
	DataInterval int
	LogInterval  int
}

type Metrics struct {
	CPUUsage    stats.CPUStat
	LoadAverage float64
}

func NewCollector(cfg *config.Config, opts Options) *StatsCollector {
	return &StatsCollector{
		history: make([]Metrics, 0, opts.LogInterval),
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
		// TODO: Output statistics here
		case <-ctx.Done():
			return
		}
	}
}

func (sc *StatsCollector) saveHistory() {
	entry := collectMetric(sc.cfg)

	sc.mu.Lock()
	fmt.Printf("Collected stats: %+v\n", entry)
	sc.history = append(sc.history, entry)
	sc.mu.Unlock()
}
