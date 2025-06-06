package app

import (
	"testing"

	"github.com/a2gx/sys-stats/internal/stats"
	"github.com/stretchr/testify/assert"
)

func TestCalculateCPUUsage(t *testing.T) {
	history := []History{
		{CPUUsage: stats.CPUStat{User: 10, System: 20, Idle: 70}},
		{CPUUsage: stats.CPUStat{User: 30, System: 10, Idle: 60}},
	}
	count := len(history)
	result := calculateCPUUsage(history, count)

	assert.InDelta(t, 20.0, result.User, 0.001, "User CPU usage должен быть средним")
	assert.InDelta(t, 15.0, result.System, 0.001, "System CPU usage должен быть средним")
	assert.InDelta(t, 65.0, result.Idle, 0.001, "Idle CPU usage должен быть средним")
}

func TestCalculateLoadAverage(t *testing.T) {
	history := []History{
		{LoadAverage: 1.5},
		{LoadAverage: 2.5},
		{LoadAverage: 3.0},
	}
	count := len(history)
	result := calculateLoadAverage(history, count)

	assert.InDelta(t, 2.333, result, 0.001, "LoadAverage должен быть средним")
}

func TestCalculateDiskUsage(t *testing.T) {
	history := []History{
		{DiskUsage: stats.DiskUsage{TPS: 100, KBps: 200}},
		{DiskUsage: stats.DiskUsage{TPS: 300, KBps: 400}},
	}
	count := len(history)
	result := calculateDiskUsage(history, count)

	assert.InDelta(t, 200.0, result.TPS, 0.001, "TPS должен быть средним")
	assert.InDelta(t, 300.0, result.KBps, 0.001, "KBps должен быть средним")
}
