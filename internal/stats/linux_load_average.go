//go:build linux

package stats

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getLoadAverageImpl() (float64, error) {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return 0, fmt.Errorf("failed to read /proc/loadavg: %w", err)
	}

	// "0.44 0.52 0.59 1/123 4567"
	fields := strings.Fields(string(data))
	if len(fields) < 3 {
		return 0, fmt.Errorf("unexpected format in /proc/loadavg")
	}

	var sum float64
	for i := 0; i < 3; i++ {
		value, err := strconv.ParseFloat(fields[i], 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse loadavg field %d: %w", i, err)
		}
		sum += value
	}

	average := sum / 3
	return average, nil
}
