//go:build linux

package stats

import (
	"bytes"
	"errors"
	"os/exec"
	"regexp"
	"strconv"
	
	"github.com/a2gx/sys-stats/pkg/utils"
)

type CPUStat struct {
	User   float64
	System float64
	Idle   float64
}

// Пример строки: "Cpu(s):  5.3%us,  2.1%sy,  0.0%ni, 90.9%id, ..."
var re = regexp.MustCompile(`([\d.]+)%us, *([\d.]+)%sy,.*?([\d.]+)%id`)

// GetCpuUsage извлекает %user, %system и %idle из команды `top` на Linux
func GetCpuUsage() (CPUStat, error) {
	cmd := exec.Command("top", "-bn1") // b - batch mode, n1 - один запуск
	out, err := cmd.Output()
	if err != nil {
		return CPUStat{}, err
	}

	lines := bytes.Split(out, []byte("\n"))
	var cpuLine string
	for _, line := range lines {
		if bytes.Contains(line, []byte("Cpu(s):")) {
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
		User:   utils.Round(user, 2),
		System: utils.Round(system, 2),
		Idle:   utils.Round(idle, 2),
	}, nil
}
