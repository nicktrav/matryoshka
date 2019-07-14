package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	g "github.com/nicktrav/matryoshka/pkg/graph"
	"github.com/nicktrav/matryoshka/pkg/lang"
)

// Print a minimal set of metadata about each dependency in a given directory.
//
// Usage of ./print:
//
//  -dir string
//        the directory of files to parse
//  -start string
//        the directory of files to parse (default "all")

func main() {
	var dir, start string
	flag.StringVar(&dir, "dir", "", "the directory of files to parse")
	flag.StringVar(&start, "start", "all", "the dependency to start from")
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

	graph := g.NewDependencyGraph()
	graph.Construct(parser.Deps())

	fmt.Println("Found the following dependencies:")
	fmt.Println()
	for i, dep := range graph.Deps() {
		fmt.Printf("\t%3d: %s -> %s\n", i, dep.Name, getDeps(dep))
	}

	fmt.Printf("\nWalking the graph from node '%s' ...\n", start)
	fmt.Println()

	printer := g.NewDepPrinter()
	walker := g.NewWalker(printer)

	err = walker.Walk(graph, start)
	if err != nil {
		log.Fatal(err)
	}
}

// getDeps returns a slice of names for each of the given Dependency's own
// dependencies.
func getDeps(dep *g.Dependency) []string {
	var names []string
	for _, d := range dep.Dependencies {
		names = append(names, d.Name)
	}
	return names
}
