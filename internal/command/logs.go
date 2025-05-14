package command

import (
	"github.com/a2gx/sys-stats/internal/daemon"
	"github.com/spf13/cobra"
)

func NewCommandLogs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Show daemon logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return daemon.LogsDaemon()
		},
	}

	return cmd
}
