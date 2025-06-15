//go:build linux

package stats

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// DiskUsage содержит статистику использования дисков
type DiskUsage struct {
	TPS  float64 // Передач в секунду
	KBps float64 // Килобайт в секунду (чтение + запись)
}

type rawDiskUsage struct {
	readCompleted  uint64
	writeCompleted uint64
	sectorsRead    uint64
	sectorsWritten uint64
}

// readDiskStats читает данные из /proc/diskstats для всех дисков
func readDiskStats() (map[string]rawDiskUsage, error) {
	data, err := os.ReadFile("/proc/diskstats")
	if err != nil {
		return nil, fmt.Errorf("failed to read /proc/diskstats: %w", err)
	}

	diskStats := make(map[string]rawDiskUsage)

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 14 {
			continue // Пропускаем некорректные строки
		}

		// Отфильтровываем только физические диски
		deviceName := fields[2]
		if !isPhysicalDisk(deviceName) {
			continue
		}

		readCompleted, err := strconv.ParseUint(fields[3], 10, 64)
		if err != nil {
			continue
		}

		sectorsRead, err := strconv.ParseUint(fields[5], 10, 64)
		if err != nil {
			continue
		}

		writeCompleted, err := strconv.ParseUint(fields[7], 10, 64)
		if err != nil {
			continue
		}

		sectorsWritten, err := strconv.ParseUint(fields[9], 10, 64)
		if err != nil {
			continue
		}

		diskStats[deviceName] = rawDiskUsage{
			readCompleted:  readCompleted,
			writeCompleted: writeCompleted,
			sectorsRead:    sectorsRead,
			sectorsWritten: sectorsWritten,
		}
	}

	if len(diskStats) == 0 {
		return nil, fmt.Errorf("no disk stats found in /proc/diskstats")
	}

	return diskStats, nil
}

// isPhysicalDisk определяет, является ли устройство физическим диском
func isPhysicalDisk(name string) bool {
	// Типичные префиксы физических дисков в Linux
	prefixes := []string{"sd", "hd", "vd", "nvme", "xvd"}

	for _, prefix := range prefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}

	return false
}

// totalDiskUsage вычисляет статистику использования дисков
func totalDiskUsage(start, end map[string]rawDiskUsage, durationSec float64) DiskUsage {
	var totalReadCompleted, totalWriteCompleted uint64
	var totalSectorsRead, totalSectorsWritten uint64

	// Суммируем статистику по всем дискам
	for deviceName, endStat := range end {
		if startStat, ok := start[deviceName]; ok {
			totalReadCompleted += endStat.readCompleted - startStat.readCompleted
			totalWriteCompleted += endStat.writeCompleted - startStat.writeCompleted
			totalSectorsRead += endStat.sectorsRead - startStat.sectorsRead
			totalSectorsWritten += endStat.sectorsWritten - startStat.sectorsWritten
		}
	}

	// Вычисляем TPS (передач в секунду)
	tps := float64(totalReadCompleted+totalWriteCompleted) / durationSec

	// Вычисляем KBps (KB/s, сектор = 512 байт)
	kbps := float64(totalSectorsRead+totalSectorsWritten) * 512.0 / 1024.0 / durationSec

	return DiskUsage{
		TPS:  tps,
		KBps: kbps,
	}
}

// getDiskUsageImpl возвращает статистику использования дисков
func getDiskUsageImpl() (DiskUsage, error) {
	start, err := readDiskStats()
	if err != nil {
		return DiskUsage{}, err
	}

	// Ждем 1 секунду для получения дельты
	time.Sleep(1 * time.Second)

	end, err := readDiskStats()
	if err != nil {
		return DiskUsage{}, err
	}

	return totalDiskUsage(start, end, 1.0), nil
}
