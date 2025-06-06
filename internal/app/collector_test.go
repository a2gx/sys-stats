package app

import (
	"context"
	"testing"
	"time"

	"github.com/a2gx/sys-stats/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestStatsCollector_SaveAndReadHistory(t *testing.T) {
	cfg := &config.Config{}
	opts := Options{DataInterval: 2, LogInterval: 1}
	c := NewCollector(cfg, nil, opts)

	// Добавляем несколько записей
	c.saveHistory()
	c.saveHistory()
	c.saveHistory()

	h := c.readHistory()
	// Проверяем, что возвращается History и история очищается
	assert.IsType(t, History{}, h)
	assert.Equal(t, 0, len(c.history))
}

func TestStatsCollector_StartAndStop(t *testing.T) {
	cfg := &config.Config{}
	opts := Options{DataInterval: 2, LogInterval: 1}
	c := NewCollector(cfg, nil, opts)

	ctx, cancel := context.WithCancel(context.Background())
	c.Start(ctx)
	time.Sleep(120 * time.Millisecond)
	cancel()
	// Проверяем, что горутины завершаются без паники
}
