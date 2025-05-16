//go:build linux

package stats

import (
	"bytes"
	"errors"
	"os/exec"
	"strconv"
	"strings"

	"github.com/a2gx/sys-stats/pkg/utils"
)

func GetLoadAverage() (float64, error) {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return 0, 0, 0, err
	}

	errors

	// "0.44 0.52 0.59 1/123 4567"
	fields := strings.Fields(string(data))
	if len(fields) < 3 {
		return 0, 0, 0, errors.New("unexpected format of /proc/loadavg")
	}

	sum, ctn := 0.0, 0.0

	for i := 0; i < 3; i++ {
		loadAvg, err := strconv.ParseFloat(fields[i], 64)
		if err != nil {
			continue
		}

		sum += loadAvg
		ctn += 1
	}

	return utils.Round(sum/ctn, 2), nil
}
