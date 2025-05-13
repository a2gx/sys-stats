package command

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewLogsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Log the sys-stats daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Logging sys-stats daemon...")
			return nil
		},
	}

	return cmd
}
