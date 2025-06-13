//go:build darwin

package stats

import (
	"bytes"
	"errors"
	"os/exec"
	"regexp"
	"strconv"
)

type CPUStat struct {
	User   float64
	System float64
	Idle   float64
}

// Пример строки: "CPU usage: 7.74% user, 8.91% sys, 83.33% idle"
var re = regexp.MustCompile(`([\d.]+)% user, ([\d.]+)% sys, ([\d.]+)% idle`)

func getCpuUsageImpl() (CPUStat, error) {
	cmd := exec.Command("top", "-l", "1", "-n", "0")
	out, err := cmd.Output()
	if err != nil {
		return CPUStat{}, err
	}

	lines := bytes.Split(out, []byte("\n"))
	var cpuLine string
	for _, line := range lines {
		if bytes.HasPrefix(line, []byte("CPU usage:")) {
			cpuLine = string(line)
			break
		}
	}

	if cpuLine == "" {
		return CPUStat{}, errors.New("CPU usage line not found")
	}

	matches := re.FindStringSubmatch(cpuLine)
	if len(matches) != 4 {
		return CPUStat{}, errors.New("failed to parse CPU usage line")
	}

	user, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return CPUStat{}, err
	}
	system, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return CPUStat{}, err
	}
	idle, err := strconv.ParseFloat(matches[3], 64)
	if err != nil {
		return CPUStat{}, err
	}

	return CPUStat{
		User:   user,
		System: system,
		Idle:   idle,
	}, nil
}
