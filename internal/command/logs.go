package command

import (
	"github.com/a2gx/sys-stats/internal/daemon"
	"github.com/spf13/cobra"
)

func NewLogs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Display logs from the sys-stats daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Создаем менеджер демона с параметрами по умолчанию
			dm := daemon.NewDaemonManager(nil)

			// Отображаем логи демона
			return dm.LogsDaemon()
		},
	}

	return cmd
}
