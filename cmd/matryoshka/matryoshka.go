package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nicktrav/matryoshka/pkg/graph"
	"github.com/nicktrav/matryoshka/pkg/lang"
)

// Enforce state from dependency files in a given directory.
//
// Usage:
//
//  -dir string (required)
//        the directory of files to parse
//  -debug
//        print debug output
//  -dep string
//        the dependency to start from (default "all")
//  -dry-run
//        do not run the 'meet' actions
//  -no-color
//        disable color in output
//

func main() {
	var dir, dep string
	var disableColor, dryRun, debug bool
	flag.StringVar(&dir, "dir", "", "the directory of files to parse")
	flag.StringVar(&dep, "dep", "all", "the dependency to start from")
	flag.BoolVar(&disableColor, "no-color", false, "disable color in output")
	flag.BoolVar(&debug, "debug", false, "print debug output")
	flag.BoolVar(&dryRun, "dry-run", false, "do not run the 'meet' actions")
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

	var printOptions []graph.PrintOption
	if !disableColor {
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
	err = walker.Walk(depGraph, dep)
	if err != nil {
		log.Fatal(err)
	}
}
