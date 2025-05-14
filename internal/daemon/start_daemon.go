package daemon

func StartDaemon(logInterval, dataInterval int) error {
	RunDaemon(logInterval, dataInterval)
	return nil
}
