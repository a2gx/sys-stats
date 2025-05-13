package command

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewStopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the sys-stats daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Stopping sys-stats daemon...")
			return nil
		},
	}

	return cmd
}
