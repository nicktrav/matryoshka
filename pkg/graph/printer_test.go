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

	wanted := "} foo ✔\n"
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

	wanted := "} foo ✖\n"
	if buf.String() != wanted {
		t.Errorf("wanted string '%s'; got %s", wanted, buf.String())
	}
}
