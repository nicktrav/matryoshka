package graph

import (
	"github.com/nicktrav/matryoshka/pkg/actions"
)

// State is the current state of the Dependency
type State int

const (

	// The state of this dep has not yet been evaluated.
	Unknown State = iota

	// The dep has been evaluated and either the deps of this dep are
	// unsatisfied, or at least one of the met actions returns with an error,
	// or both
	Unsatisfied

	// The dep has been evaluated and all the deps of this dep are satisfied
	// and the met actions all return without error.
	Satisfied
)

// A dependency represents a node in the dependency graph.
type Dependency struct {

	// Name is the name of the dependency.
	Name string

	// Dependencies is the list of dependencies that must be satisfied before
	// this dependency is satisfied.
	Dependencies []*Dependency

	// MetAction is the list of commands to run to determine whether the dependency is satisfied.
	// TODO(nickt): Change schema to be only one action
	MetActions []actions.Action

	// MeetAction is the list of commands to run to attempt to satisfy the dependency.
	MeetActions []actions.Action

	// State is the cached state of the Dependency.
	State
}

// NewDependency returns a pointer to a new Dependency.
func NewDependency(name string) *Dependency {
	return &Dependency{
		Name:         name,
		Dependencies: []*Dependency{},
		MetActions:   []actions.Action{},
		MeetActions:  []actions.Action{},
		State:        Unknown,
	}
}
