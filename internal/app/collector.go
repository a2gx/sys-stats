package app

import (
	"context"
	"github.com/a2gx/sys-stats/internal/config"
	"sync"
	"time"
)

type StatsCollector struct {
	history []map[string]interface{}
	mu      sync.Mutex
	cfg     *config.Config
	opt     Options
}

type Options struct {
	DataInterval int
	LogInterval  int
}

func NewCollector(cfg *config.Config, opts Options) *StatsCollector {
	return &StatsCollector{
		history: make([]map[string]interface{}, opts.LogInterval),
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
			// TODO: Collect statistics here
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
