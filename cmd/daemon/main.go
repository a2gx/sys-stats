package main

import (
	"fmt"
	"log"
	"os"

	"github.com/a2gx/sys-stats/internal/command"

	"github.com/spf13/cobra"
)

var (
	release   = "dev"
	buildDate = "unknown"
	gitHash   = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "sys-stats",
		Short:   "System Statistics Daemon",
		Version: fmt.Sprintf("%s (%s) built %s", release, gitHash, buildDate),
	}

	// Adding subcommands
	rootCmd.AddCommand(
		command.NewRunCommand(),
		command.NewLogsCommand(),
		command.NewStopCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
