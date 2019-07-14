package actions

// Action represents some action to take.
type Action interface {

	// Run attempts to perform the action action, returning an error if the
	// action could not be completed.
	Run() error
}

// Debugger takes an appropriate debug option on an Action.
type Debugger interface {

	// Debugger is also an Action.
	Action

	// Debug takes an appropriate debug action on an Action. For example,
	// setting the output stream to stderr, etc.
	Debug()
}
