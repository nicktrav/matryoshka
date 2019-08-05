package apply

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/nicktrav/matryoshka/pkg/graph"
	"github.com/nicktrav/matryoshka/pkg/lang"
)

var (
	dir     string
	rootDep string
	noColor bool
	debug   bool
	dryRun  bool
)

const defaultRoot = "all"

// NewCommand returns a new command for applying dependencies.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}

	cmd.Flags().StringVar(&dir, "dir", "", "Directory of deps")
	cmd.Flags().StringVar(&rootDep, "dep", defaultRoot, "Root of the dependency graph")
	cmd.Flags().BoolVar(&noColor, "no-color", false, "Disable color printing")
	cmd.Flags().BoolVar(&debug, "debug", false, "Enable debug output")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Do not attempt to satisfy dependencies")

	return cmd
}

// run attempts to enforce state from dependency files in a given directory.
func run() error {
	if dir == "" {
		return errors.New("dir is a required argument")
	}

	parser := lang.NewParser(dir)
	err := parser.Run()
	if err != nil {
		return err
	}

	depGraph := graph.NewDependencyGraph()
	depGraph.Construct(parser.Deps())

	var printOptions []graph.PrintOption
	if !noColor {
		printOptions = append(printOptions, graph.WithColor)
	}
	printer := graph.NewDepPrinter(printOptions...)

	var executorOptions []graph.ExecutorOption
	if debug {
		executorOptions = append(executorOptions, graph.Debug)
	}
	if dryRun {
		executorOptions = append(executorOptions, graph.DryRun)
	}
	executor := graph.NewExecutor(executorOptions...)

	v := graph.NewCompositeVisitor(printer, executor)
	walker := graph.NewWalker(v)
	err = walker.Walk(depGraph, rootDep)
	if err != nil {
		return err
	}

	return nil
}
