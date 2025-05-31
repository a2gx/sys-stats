//go:build darwin

package stats

func GetDiskUsage() (DiskStat, error) {
	return DiskStat{}, nil
}
