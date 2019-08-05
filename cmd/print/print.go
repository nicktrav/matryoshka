package print

import (
	"errors"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/nicktrav/matryoshka/pkg/graph"
	"github.com/nicktrav/matryoshka/pkg/lang"
)

var (
	dir     string
	rootDep string
)

// NewCommand returns a new command for the printing dep graph.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "print",
		Short: "Print the dep tree",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}

	cmd.Flags().StringVar(&dir, "dir", "", "Directory of deps")
	cmd.Flags().StringVar(&rootDep, "root", "", "Root of the dependency graph")

	return cmd
}

// Print a minimal set of metadata about each dependency in a given directory.
func run() error {
	if dir == "" {
		return errors.New("dir is a required argument")
	}

	parser := lang.NewParser(dir)
	err := parser.Run()
	if err != nil {
		log.Fatal(err)
	}

	depGraph := graph.NewDependencyGraph()
	depGraph.Construct(parser.Deps())

	fmt.Println("Found the following dependencies:")
	fmt.Println()
	for i, dep := range depGraph.Deps() {
		fmt.Printf("\t%3d: %s -> %s\n", i, dep.Name, getDeps(dep))
	}

	return nil
}

// getDeps returns a slice of names for each of the given Dependency's own
// dependencies.
func getDeps(dep *graph.Dependency) []string {
	var names []string
	for _, d := range dep.Dependencies {
		names = append(names, d.Name)
	}
	return names
}
