package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nicktrav/matryoshka/pkg/graph"
	"github.com/nicktrav/matryoshka/pkg/lang"
)

// Print a minimal set of metadata about each dependency in a given directory.
//
// Usage of ./print:
//
//  -dir string (required)
//        the directory of files to parse
//  -no-color
//        disable color in output
//  -start string
//        the dependency to start from (default "all")
//

func main() {
	var dir, start string
	var disableColor bool
	flag.StringVar(&dir, "dir", "", "the directory of files to parse")

	flag.StringVar(&start, "start", "all", "the dependency to start from")
	flag.BoolVar(&disableColor, "no-color", false, "disable color in output")
	flag.Parse()

	if dir == "" {
		fmt.Println("-dir is a required argument")
		os.Exit(1)
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
