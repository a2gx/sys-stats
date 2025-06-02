package command

import (
	"github.com/a2gx/sys-stats/internal/config"
	"github.com/a2gx/sys-stats/internal/daemon"
	"github.com/spf13/cobra"
)

func NewCommandRun() *cobra.Command {
	var detect bool
	var logInterval, dataInterval int
	var configFile string

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Start the daemon to collect system statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.NewConfig(configFile)
			if err != nil {
				return err
			}

			return daemon.StartDaemon(cfg, detect, logInterval, dataInterval)
		},
	}

	// Flags...
	cmd.Flags().StringVar(&configFile, "config", "/configs/config.yaml", "Path to configuration file")
	cmd.Flags().BoolVarP(&detect, "detect", "d", false, "Run the daemon in background mode")
	cmd.Flags().IntVarP(&logInterval, "log-interval", "n", 5, "Log output interval (in seconds)")
	cmd.Flags().IntVarP(&dataInterval, "data-interval", "m", 15, "Data collection period (in seconds)")

	return cmd
}
