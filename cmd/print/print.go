package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nicktrav/matryoshka/pkg"
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

	graph := pkg.NewDependencyGraph()
	graph.Construct(parser.Modules())

	fmt.Println("Found the following dependencies:")
	fmt.Println()
	for i, dep := range graph.Deps() {
		fmt.Printf("\t%d: %s -> %s\n", i, dep.Name, getDeps(dep))
	}

	fmt.Printf("\nWalking the graph from node '%s' ...\n", start)
	fmt.Println()

	visitor := &depPrinter{}
	err = graph.Walk(start, visitor)
	if err != nil {
		log.Fatal(err)
	}
}

// getDeps returns a slice of names for each of the given Dependency's own
// dependencies.
func getDeps(dep *pkg.Dependency) []string {
	var names []string
	for _, d := range dep.Dependencies {
		names = append(names, d.Name)
	}
	return names
}

// depPrinter is a NodeVisitor that prints out some metadata about each Dep that
// it visits. The output is indented to represent the dependency graph.
type depPrinter struct {

	// indentLevel is the amount to indent each log line
	indentLevel int
}

// Visit will print information about the Dependency.
func (p *depPrinter) Visit(dep *pkg.Dependency) {
	p.printf(dep, "visit")
}

// PreVisit increments the indentation before printing the pre-visit message.
func (p *depPrinter) PreVisit(dep *pkg.Dependency) {
	p.indentLevel++
	p.printf(dep, "pre-visit")
}

// PreVisit printe the post-visit message before decrementing the indentation.
func (p *depPrinter) PostVisit(dep *pkg.Dependency) {
	p.printf(dep, "post-visit")
	p.indentLevel--
}

// Errors returns an empty list of errors, as there's nothing that can go wrong
// with this visitor.
func (p *depPrinter) Errors() []error {
	return []error{}
}

// printf is a simple wrapper around fmt.Println that adds the requisite amount
// of indentation to the line, as well as prefixes the message with the name of
// the Dependency.
func (p *depPrinter) printf(dep *pkg.Dependency, format string, a ...interface{}) {
	var indent string
	for i := 0; i < p.indentLevel; i++ {
		indent += "\t"
	}

	prefix := fmt.Sprintf("%s[%s]: ", indent, dep.Name)
	fmt.Printf(prefix+format+"\n", a...)
}
