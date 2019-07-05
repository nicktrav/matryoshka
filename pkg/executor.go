package pkg

// DepExecutor is a NodeVisitor that will attempt to satisfy each dependency,
// collecting any errors it encounters along the way.
type DepExecutor struct {

	// errors are the errors collector by this executor.
	errors []error
}

func (e *DepExecutor) PreVisit(dep *Dependency) {
	// do nothing
}

// Visit attempts to satisfy the Dependency.
func (e *DepExecutor) Visit(dep *Dependency) {
	errors := dep.Satisfy()
	if len(errors) > 0 {
		e.errors = append(e.errors, errors...)
		return
	}
}

func (e *DepExecutor) PostVisit(dep *Dependency) {
	// do nothing
}

// Errors returns a slice of any errors that were collected by the executor.
func (e *DepExecutor) Errors() []error {
	return e.errors
}
