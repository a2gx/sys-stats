//go:build darwin

package stats

// DiskUsage содержит статистику использования дисков
type DiskUsage struct {
	TPS  float64 // Передач в секунду
	KBps float64 // Килобайт в секунду (чтение + запись)
}

func getDiskUsageImpl() (DiskUsage, error) {
	return DiskUsage{}, nil
}
