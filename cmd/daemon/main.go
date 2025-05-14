package main

import (
	"os"

	"github.com/a2gx/sys-stats/internal/command"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "sys-stats",
		Short: "System Statistics Daemon",
	}

	// Добавление подкоманд
	rootCmd.AddCommand(
		command.NewCommandRun(),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
