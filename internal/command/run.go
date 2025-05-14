package command

import (
	"github.com/a2gx/sys-stats/internal/daemon"
	"github.com/spf13/cobra"
)

func NewCommandRun() *cobra.Command {
	var detect bool
	var logInterval, dataInterval int

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Start the daemon to collect system statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			return daemon.StartDaemon(detect, logInterval, dataInterval)
		},
	}

	// Flags...
	cmd.Flags().BoolVarP(&detect, "detect", "d", false, "Run the daemon in background mode")
	cmd.Flags().IntVarP(&logInterval, "log-interval", "n", 5, "Log output interval (in seconds)")
	cmd.Flags().IntVarP(&dataInterval, "data-interval", "m", 15, "Data collection period (in seconds)")

	return cmd
}
