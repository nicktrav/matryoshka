package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/nicktrav/matryoshka/cmd/apply"
	"github.com/nicktrav/matryoshka/cmd/print"
)

const commandName = "matryoshka"

func main() {
	rootCmd := newCommand(os.Args[1:])
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// newCommand returns a new command that acts as the root command, to which
// all other sub-commands are bound.
func newCommand(args []string) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use: commandName,
	}

	rootCmd.SetArgs(args)
	rootCmd.AddCommand(print.NewCommand())
	rootCmd.AddCommand(apply.NewCommand())

	return rootCmd
}
