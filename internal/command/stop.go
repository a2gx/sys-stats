package command

import (
	"github.com/a2gx/sys-stats/internal/daemon"
	"github.com/spf13/cobra"
)

func NewCommandStop() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Остановка демона системной статистики",
		RunE: func(cmd *cobra.Command, args []string) error {
			return daemon.StopDaemon()
		},
	}

	return cmd
}
