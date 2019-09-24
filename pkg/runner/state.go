package runner

// State indicates the current state
// of a task/runner
type State int

const (
	// INITIALIZED indicates the Runner
	// is configured and ready to start
	INITIALIZED State = iota + 1
	// RUNNING indicates the runner is
	// currently timing the given task
	RUNNING
	// BREAKING indicates it is time
	// to break and move on to the next
	// pomodoro if configured
	BREAKING
	// COMPLETE indicates the end of each
	// configured pomodoro for a given task
	// has been reached
	COMPLETE
	// SUSPENDED indicates the user has
	// suspended the current task
	SUSPENDED
)

func (s State) String() string {
	switch s {
	case INITIALIZED:
		return "INITIALIZED"
	case RUNNING:
		return "RUNNING"
	case BREAKING:
		return "BREAKING"
	case COMPLETE:
		return "COMPLETE"
	case SUSPENDED:
		return "SUSPENDED"
	}
	return ""
}
