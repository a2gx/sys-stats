package stats

import (
	"testing"

	"github.com/a2gx/sys-stats/internal/config"
	"github.com/stretchr/testify/require"
)

// Сохраняем оригинальные функции для восстановления после теста
var (
	origGetLoadAverage = GetLoadAverage
	origGetCpuUsage    = GetCpuUsage
	origGetDiskUsage   = GetDiskUsage
)

func mockStats() {
	GetLoadAverage = func() (float64, error) {
		return 1.23, nil
	}
	GetCpuUsage = func() (CPUStat, error) {
		return CPUStat{User: 10, System: 5, Idle: 85}, nil
	}
	GetDiskUsage = func() (DiskUsage, error) {
		return DiskUsage{TPS: 100, KBps: 2048}, nil
	}
}

func restoreStats() {
	GetLoadAverage = origGetLoadAverage
	GetCpuUsage = origGetCpuUsage
	GetDiskUsage = origGetDiskUsage
}

func TestNewCollector(t *testing.T) {
	cfg := &config.Config{
		LoadAverage: true,
		CPUUsage:    true,
		DiskUsage:   false,
	}
	opts := CollectorOptions{
		LogInterval:  10,
		DataInterval: 5,
	}

	c := NewCollector(cfg, opts)

	require.NotNil(t, c)
	require.Equal(t, cfg, c.cfg)
	require.Equal(t, opts, c.opt)
	require.Nil(t, c.history)
}

func TestCollector_HistoryCollect(t *testing.T) {
	mockStats()
	defer restoreStats()

	cfg := &config.Config{
		LoadAverage: true,
		CPUUsage:    true,
		DiskUsage:   true,
	}
	c := NewCollector(cfg, CollectorOptions{})

	h := c.HistoryCollect()

	require.NotNil(t, h)
	require.Equal(t, 1.23, h.LoadAverage)
	require.Equal(t, CPUStat{User: 10, System: 5, Idle: 85}, h.CPUUsage)
	require.Equal(t, DiskUsage{TPS: 100, KBps: 2048}, h.DiskUsage)
}

func TestCollector_HistoryCalculate(t *testing.T) {
	c := NewCollector(&config.Config{}, CollectorOptions{})

	history := []*History{
		{
			LoadAverage: 1.0,
			CPUUsage:    CPUStat{User: 10, System: 5, Idle: 85},
			DiskUsage:   DiskUsage{TPS: 100, KBps: 2048},
		},
		{
			LoadAverage: 3.0,
			CPUUsage:    CPUStat{User: 30, System: 15, Idle: 55},
			DiskUsage:   DiskUsage{TPS: 200, KBps: 4096},
		},
	}

	result := c.HistoryCalculate(history)
	require.NotNil(t, result)
	require.InDelta(t, 2.0, result.LoadAverage, 0.0001)
	require.InDelta(t, 20.0, result.CPUUsage.User, 0.0001)
	require.InDelta(t, 10.0, result.CPUUsage.System, 0.0001)
	require.InDelta(t, 70.0, result.CPUUsage.Idle, 0.0001)
	require.InDelta(t, 150.0, result.DiskUsage.TPS, 0.0001)
	require.InDelta(t, 3072.0, result.DiskUsage.KBps, 0.0001)
}
