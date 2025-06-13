//go:build linux

package stats

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type CPUStat struct {
	User   float64
	System float64
	Idle   float64
}

type rawCPU struct {
	user    uint64
	nice    uint64
	system  uint64
	idle    uint64
	iowait  uint64
	irq     uint64
	softirq uint64
}

func readCPUStat() (rawCPU, error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return rawCPU{}, fmt.Errorf("failed to read /proc/stat: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)
			if len(fields) < 8 {
				return rawCPU{}, fmt.Errorf("unexpected format in /proc/stat")
			}
			values := make([]uint64, 8)
			for i := 0; i < 8; i++ {
				val, err := strconv.ParseUint(fields[i+1], 10, 64)
				if err != nil {
					return rawCPU{}, fmt.Errorf("invalid value in /proc/stat: %w", err)
				}
				values[i] = val
			}
			return rawCPU{
				user:    values[0],
				nice:    values[1],
				system:  values[2],
				idle:    values[3],
				iowait:  values[4],
				irq:     values[5],
				softirq: values[6],
			}, nil
		}
	}
	return rawCPU{}, fmt.Errorf("cpu line not found in /proc/stat")
}

func calculateCPUStat(start, end rawCPU) CPUStat {
	startTotal := start.user + start.nice + start.system + start.idle + start.iowait + start.irq + start.softirq
	endTotal := end.user + end.nice + end.system + end.idle + end.iowait + end.irq + end.softirq

	totalDelta := float64(endTotal - startTotal)
	if totalDelta == 0 {
		return CPUStat{}
	}

	userDelta := float64((end.user + end.nice) - (start.user + start.nice))
	systemDelta := float64((end.system + end.irq + end.softirq) - (start.system + start.irq + start.softirq))
	idleDelta := float64((end.idle + end.iowait) - (start.idle + start.iowait))

	return CPUStat{
		User:   100 * userDelta / totalDelta,
		System: 100 * systemDelta / totalDelta,
		Idle:   100 * idleDelta / totalDelta,
	}
}

func getCpuUsageImpl() (CPUStat, error) {
	start, err := readCPUStat()
	if err != nil {
		return CPUStat{}, err
	}

	time.Sleep(1 * time.Second)

	end, err := readCPUStat()
	if err != nil {
		return CPUStat{}, err
	}

	return calculateCPUStat(start, end), nil
}
