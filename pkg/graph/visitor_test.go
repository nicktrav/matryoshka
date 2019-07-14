package graph

import (
	"errors"
	"strings"
	"testing"
)

type testVisitor struct {
	visitCalled     bool
	preVisitCalled  bool
	postVisitCalled bool
}

func (v *testVisitor) Visit(dep *Dependency) error {
	v.visitCalled = true
	return nil
}

func (v *testVisitor) PreVisit(dep *Dependency) {
	v.preVisitCalled = true
}

func (v *testVisitor) PostVisit(dep *Dependency) {
	v.postVisitCalled = true
}

func TestCompositeVisitor_Visit(t *testing.T) {
	v1 := &testVisitor{}
	v2 := &testVisitor{}

	compositeVisitor := NewCompositeVisitor(v1, v2)
	err := compositeVisitor.Visit(NewDependency("foo"))

	if err != nil {
		t.Fatalf("wanted no error; got %+v", err)
	}

	if !v1.visitCalled {
		t.Error("wanted v1.Visit called")
	}

	if !v2.visitCalled {
		t.Error("wanted v2.Visit called")
	}
}

func TestCompositeVisitor_PreVisit(t *testing.T) {
	v1 := &testVisitor{}
	v2 := &testVisitor{}

	compositeVisitor := NewCompositeVisitor(v1, v2)
	compositeVisitor.PreVisit(NewDependency("foo"))

	if !v1.preVisitCalled {
		t.Error("wanted v1.PreVisit called")
	}

	if !v2.preVisitCalled {
		t.Error("wanted v2.PreVisit called")
	}
}

func TestCompositeVisitor_PostVisit(t *testing.T) {
	v1 := &testVisitor{}
	v2 := &testVisitor{}

	compositeVisitor := NewCompositeVisitor(v1, v2)
	compositeVisitor.PostVisit(NewDependency("foo"))

	if !v1.postVisitCalled {
		t.Error("wanted v1.PostVisit called")
	}

	if !v2.postVisitCalled {
		t.Error("wanted v2.PostVisit called")
	}
}

type failingVisitor struct {
	err error
}

func (v *failingVisitor) Visit(dep *Dependency) error {
	return v.err
}

func (v *failingVisitor) PreVisit(dep *Dependency) {
}

func (v *failingVisitor) PostVisit(dep *Dependency) {
}

func TestCompositeVisitor_Visit_returnsError(t *testing.T) {
	e := errors.New("oh oh")
	v1 := &failingVisitor{err: e}
	compositeVisitor := NewCompositeVisitor(v1)

	err := compositeVisitor.Visit(NewDependency("foo"))
	if err == nil {
		t.Fatalf("wanted error; got none")
	}

	if !strings.HasPrefix(err.Error(), "composite") {
		t.Errorf("wanted error to have prefix 'composite'; got %+v", err)
	}

	if !strings.Contains(err.Error(), v1.err.Error()) {
		t.Errorf("wanted error to conatin wrapped error '%s'; got %+v", v1.err, err)
	}
}
