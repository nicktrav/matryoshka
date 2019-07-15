package graph

import (
	"bytes"
	"os"
	"testing"
)

func TestNewDepPrinter(t *testing.T) {
	visitor := NewDepPrinter()

	printer, ok := visitor.(*depPrinter)
	if !ok {
		t.Fatalf("wanted visitor to be a depPrinter; got %+v", visitor)
	}

	if printer.indentLevel != 0 {
		t.Errorf("wanted indent level to start at zero; got %d", printer.indentLevel)
	}

	if printer.writer != os.Stdout {
		t.Errorf("wanted default writer to be os.Stdout; got %+v", printer.writer)
	}
}

func TestDepPrinter_PreVisit(t *testing.T) {
	buf := new(bytes.Buffer)
	printer := depPrinter{writer: buf}

	dep := NewDependency("foo")
	printer.PreVisit(dep)

	if printer.indentLevel != 1 {
		t.Errorf("wanted indentLevel one; got %d", printer.indentLevel)
	}

	wanted := "foo {\n"
	if buf.String() != wanted {
		t.Errorf("wanted string '%s'; got %s", wanted, buf.String())
	}
}

func TestDepPrinter_PostVisit_IsSatisfied(t *testing.T) {
	buf := new(bytes.Buffer)
	printer := depPrinter{writer: buf, indentLevel: 1}

	dep := NewDependency("foo")
	dep.State = Satisfied
	printer.PostVisit(dep)

	if printer.indentLevel != 0 {
		t.Errorf("wanted indentLevel zero; got %d", printer.indentLevel)
	}

	wanted := "} ✔ foo\n"
	if buf.String() != wanted {
		t.Errorf("wanted string '%s'; got %s", wanted, buf.String())
	}
}

func TestDepPrinter_PostVisit_IsUnsatisfied(t *testing.T) {
	buf := new(bytes.Buffer)
	printer := depPrinter{writer: buf, indentLevel: 1}

	dep := NewDependency("foo")
	dep.State = Unsatisfied
	printer.PostVisit(dep)

	if printer.indentLevel != 0 {
		t.Errorf("wanted indentLevel zero; got %d", printer.indentLevel)
	}

	wanted := "} ✖ foo\n"
	if buf.String() != wanted {
		t.Errorf("wanted string '%s'; got %s", wanted, buf.String())
	}
}

func TestDepPrinter_NoColor(t *testing.T) {
	printer := depPrinter{colorize: false}

	want := "foo"
	outputs := []string{
		printer.wrap(want, colorGreen),
		printer.wrap(want, colorRed),
	}

	for _, output := range outputs {
		if output != want {
			t.Errorf("wanted plain output '%s'; got '%s'", want, output)
		}
	}
}

func TestDepPrinter_Color(t *testing.T) {
	printer := depPrinter{colorize: true}

	type testInput struct {
		s     string
		color int
	}

	inputs := []testInput{
		{"foo", colorGreen},
		{"foo", colorRed},
	}

	want := []string{
		"[92mfoo[0m",
		"[91mfoo[0m",
	}

	for i, input := range inputs {
		got := printer.wrap(input.s, input.color)
		if got != want[i] {
			t.Errorf("wanted '%s'; got '%s'", want[i], got)
		}
	}
}
