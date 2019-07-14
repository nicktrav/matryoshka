package graph

import (
	"fmt"
)

// A nodeVisitor iterates over nodes in the DependencyGraph
type DepVisitor interface {

	// Visit is the main action to take when visiting a node in the graph.
	Visit(dep *Dependency) error

	// PreVisit is a hook called before the Visit function is called.
	PreVisit(dep *Dependency)

	// PostVisit is a hook called after the Visit function is called.
	PostVisit(dep *Dependency)
}

// NewCompositeVistor returns a new DepVisitor that runs the given DepVisitors
// in order on given Dependencies.
func NewCompositeVisitor(visitors ...DepVisitor) DepVisitor {
	return &compositeVisitor{
		visitors: visitors,
	}
}

// compositeVisitor is a DepVisitor that will call the given DepVisitors in
// for each Dependency evaluated.
type compositeVisitor struct {

	// visitors is a slice of DepVisitors that will run, in order when
	// evaluating a given Dependency.
	visitors []DepVisitor
}

// Visit runs Visit, in order on the given Dependency.
func (v *compositeVisitor) Visit(dep *Dependency) error {
	for _, v := range v.visitors {
		if err := v.Visit(dep); err != nil {
			return fmt.Errorf("composite: %+v", err)
		}
	}
	return nil
}

// PreVisit runs PreVisit, in order on the given Dependency.
func (v *compositeVisitor) PreVisit(dep *Dependency) {
	for _, v := range v.visitors {
		v.PreVisit(dep)
	}
}

// PostVisit runs PostVisit, in order on the given Dependency.
func (v *compositeVisitor) PostVisit(dep *Dependency) {
	for _, v := range v.visitors {
		v.PostVisit(dep)
	}
}
