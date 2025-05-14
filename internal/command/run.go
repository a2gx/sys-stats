package command

import (
	"github.com/a2gx/sys-stats/internal/daemon"
	"github.com/spf13/cobra"
)

const (
	DefaultLogInterval  = 5
	DefaultDataInterval = 15
)

func NewCommandRun() *cobra.Command {
	var logInterval, dataInterval int

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Start the daemon to collect system statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			return daemon.StartDaemon(logInterval, dataInterval)
		},
	}

	// Flags...
	cmd.Flags().IntVarP(&logInterval, "log-interval", "n", DefaultLogInterval, "Log output interval (in seconds)")
	cmd.Flags().IntVarP(&dataInterval, "data-interval", "m", DefaultDataInterval, "Data collection period (in seconds)")

	return cmd
}
