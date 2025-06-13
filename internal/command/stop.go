package command

import (
	"github.com/a2gx/sys-stats/internal/daemon"
	"github.com/spf13/cobra"
)

func NewStop() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop running daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Создаем менеджер демона с параметрами по умолчанию
			dm := daemon.NewDaemonManager(nil)

			// Останавливаем демон
			return dm.StopDaemon()
		},
	}

	return cmd
}
