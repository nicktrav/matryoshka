package actions

// Action represents some action to take
type Action interface {

	// Run attempts to perform the action action, returning an error if the
	// action could not be completed.
	Run() error
}
