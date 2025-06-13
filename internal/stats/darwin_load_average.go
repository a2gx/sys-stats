//go:build darwin

package stats

import (
	"bytes"
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

// getLoadAverageImpl возвращает среднюю загрузку системы за последние 1, 5 и 15 минут.
func getLoadAverageImpl() (float64, error) {
	cmd := exec.Command("sysctl", "-n", "vm.loadavg")
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	// Пример вывода: "{2.39 2.46 2.45}"
	output := string(bytes.TrimSpace(out))
	output = strings.Trim(output, "{}") // Удаляем фигурные скобки
	parts := strings.Fields(output)

	if len(parts) < 1 {
		return 0, errors.New("unexpected output format")
	}

	sum, ctn := 0.0, 0.0

	for _, part := range parts {
		loadAvg, err := strconv.ParseFloat(part, 64)
		if err != nil {
			continue
		}

		sum += loadAvg
		ctn += 1
	}

	return sum / ctn, nil
}
