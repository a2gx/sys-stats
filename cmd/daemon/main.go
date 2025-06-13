package main

import (
	"fmt"
	"os"

	"github.com/a2gx/sys-stats/internal/command"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cobra.Command{
		Use:     "daemon",
		Short:   "System Statistics Daemon",
		Version: getVersion(),
	}

	rootCmd.AddCommand(
		command.NewRun(),
		command.NewLogs(),
		command.NewStop(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error executing command: %v", err)
		os.Exit(1)
	}
}
