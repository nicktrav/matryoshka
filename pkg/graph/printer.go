package graph

import (
	"fmt"
	"io"
	"log"
	"os"
)

// depPrinter is a NodeVisitor that prints out some metadata about each Dep that
// it visits. The output is indented to represent the dependency graph.
type depPrinter struct {

	// indentLevel is the amount to indent each log line
	indentLevel int

	// the destination for any output that may be printed
	writer io.Writer
}

// NewDepPrinter returns a new DepVisitor that will print the dependency graph
// to Stdout.
func NewDepPrinter() DepVisitor {
	return &depPrinter{writer: os.Stdout}
}

// Visit does nothing.
func (p *depPrinter) Visit(dep *Dependency) error {
	return nil
}

// PreVisit increments the indentation before printing the pre-visit message.
func (p *depPrinter) PreVisit(dep *Dependency) {
	p.printf("%s {", dep.Name)
	p.indentLevel++
}

// PreVisit prints the post-visit message before decrementing the indentation.
func (p *depPrinter) PostVisit(dep *Dependency) {
	p.indentLevel--

	var icon string
	switch isMet(dep) {
	case true:
		icon = "✔"
	case false:
		icon = "✖"
	}

	p.printf("} %s %s", dep.Name, icon)
}

// printf is a simple wrapper around fmt.Println that adds the requisite amount
// of indentation to the line, as well as prefixes the message with the name of
// the Dependency.
func (p *depPrinter) printf(format string, a ...interface{}) {
	var indent string
	for i := 0; i < p.indentLevel; i++ {
		indent += "  "
	}

	_, err := fmt.Fprintf(p.writer, indent+format+"\n", a...)
	if err != nil {
		log.Fatal(err)
	}
}

// isMet recursively checks the state of the dep to determine the state of the
// dep.
func isMet(d *Dependency) bool {
	if d.State != Unknown {
		return d.State == Satisfied
	}

	for _, dep := range d.Dependencies {
		if !isMet(dep) {
			d.State = Unsatisfied
			return false
		}
	}

	for _, a := range d.MetActions {
		if err := a.Run(); err != nil {
			d.State = Unsatisfied
			return false
		}
	}

	d.State = Satisfied
	return true
}
