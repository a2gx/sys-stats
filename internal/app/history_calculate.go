package app

import "github.com/a2gx/sys-stats/internal/stats"

func calculateCPUUsage(history []History, count int) stats.CPUStat {
	var d stats.CPUStat

	for _, entry := range history {
		d.User += entry.CPUUsage.User
		d.System += entry.CPUUsage.System
		d.Idle += entry.CPUUsage.Idle
	}

	return stats.CPUStat{
		User:   d.User / float64(count),
		System: d.System / float64(count),
		Idle:   d.Idle / float64(count),
	}
}

func calculateLoadAverage(history []History, count int) float64 {
	var sum float64

	for _, entry := range history {
		sum += entry.LoadAverage
	}

	return sum / float64(count)
}

func calculateDiskUsage(history []History, count int) stats.DiskUsage {
	var d stats.DiskUsage

	for _, entry := range history {
		d.TPS += entry.DiskUsage.TPS
		d.KBps += entry.DiskUsage.KBps
	}

	return stats.DiskUsage{
		TPS:  d.TPS / float64(count),
		KBps: d.KBps / float64(count),
	}
}
