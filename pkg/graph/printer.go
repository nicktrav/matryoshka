package graph

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	// ANSI control sequences
	colorGreen = 92
	colorRed   = 91
	escape     = "\x1b"
)

type PrintOption func(*depPrinter)

// WithColor is a PrintOption to enable color in the output
var WithColor = func(printer *depPrinter) {
	printer.colorize = true
}

// depPrinter is a NodeVisitor that prints out some metadata about each Dep that
// it visits. The output is indented to represent the dependency graph.
type depPrinter struct {

	// indentLevel is the amount to indent each log line
	indentLevel int

	// writer is the destination for any output that may be printed
	writer io.Writer

	// colorize determines whether to print the output with color
	colorize bool
}

// NewDepPrinter returns a new DepVisitor that will print the dependency graph
// to Stdout.
func NewDepPrinter(options ...PrintOption) DepVisitor {
	printer := &depPrinter{writer: os.Stdout}

	for _, option := range options {
		option(printer)
	}

	return printer
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
		icon = p.green(fmt.Sprintf("✔ %s", dep.Name))
	case false:
		icon = p.red(fmt.Sprintf("✖ %s", dep.Name))
	}

	p.printf("} %s", icon)
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

// red returns the given string output as red.
func (p *depPrinter) red(s string) string {
	return p.wrap(s, colorRed)
}

// green returns the given string output as green.
func (p *depPrinter) green(s string) string {
	return p.wrap(s, colorGreen)
}

// wrap returns the the given string with the given color, unless coloring is
// disabled.
func (p *depPrinter) wrap(s string, color int) string {
	if !p.colorize {
		return s
	}
	return fmt.Sprintf("%s[%dm%s%s[0m", escape, color, s, escape)
}
