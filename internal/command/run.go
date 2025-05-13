package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewRunCommand() *cobra.Command {
	var detect bool

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run the sys-stats daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Running sys-stats daemon...")
			return nil
		},
	}

	// Adding flags to the command
	cmd.Flags().BoolVarP(&detect, "detect", "d", false, "Run daemon in detect mode")

	return cmd
}
